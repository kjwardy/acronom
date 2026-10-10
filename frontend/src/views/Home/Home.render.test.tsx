import React from 'react';
import axios from 'axios';
import '@testing-library/jest-dom';
import { fireEvent, render, screen } from '@testing-library/react';
import {
  AcronymEventsProvider,
  useAcronymEvents,
} from '../../contexts/AcronymEvents';
import { MetricsProvider } from '../../contexts/Metrics';
import Home from './Home';

jest.mock('axios');

const mockedGet = axios.get as jest.MockedFunction<typeof axios.get>;

const CreationControl = () => {
  const { notifyCreated } = useAcronymEvents();
  return (
    <button
      onClick={() =>
        notifyCreated({
          id: 2,
          acronym: 'API',
          definition: 'Active Pharmaceutical Ingredient',
          categories: [{ id: 2, name: 'Healthcare', colour: '#388E3C' }],
        })
      }
      type="button"
    >
      Simulate creation
    </button>
  );
};

const renderHome = (withCreationControl = false) =>
  render(
    <MetricsProvider>
      <AcronymEventsProvider>
        <Home />
        {withCreationControl && <CreationControl />}
      </AcronymEventsProvider>
    </MetricsProvider>,
  );

const mockRequests = (searchResults: object[]) => {
  mockedGet.mockImplementation((url: string) => {
    if (url === '/api/metrics/total-acronyms')
      return Promise.resolve({ data: { total_acronyms: 1 } });
    if (url === '/api/search') return Promise.resolve({ data: searchResults });
    return Promise.reject(new Error(`Unexpected request: ${url}`));
  });
};

describe('Home category results', () => {
  afterEach(() => jest.clearAllMocks());

  it('renders category labels returned by search', async () => {
    mockRequests([
      {
        id: 1,
        acronym: 'API',
        definition: 'Application Programming Interface',
        categories: [{ id: 1, name: 'Platform', colour: '#1976D2' }],
      },
    ]);
    renderHome();

    await screen.findByText('1 meaning decoded');
    fireEvent.change(screen.getByLabelText('Search acronyms'), {
      target: { value: 'API' },
    });
    fireEvent.click(screen.getByLabelText('Submit acronym search'));

    const label = await screen.findByText('Platform');
    expect(label).toBeInTheDocument();
    expect(label.parentElement?.style.backgroundColor).toBe(
      'rgb(25, 118, 210)',
    );
    expect(label.parentElement?.style.color).toBe('rgb(255, 255, 255)');
  });

  it('adds a matching created acronym and its labels without another search', async () => {
    mockRequests([
      {
        id: 1,
        acronym: 'API',
        definition: 'Application Programming Interface',
        categories: [],
      },
    ]);
    renderHome(true);

    await screen.findByText('1 meaning decoded');
    fireEvent.change(screen.getByLabelText('Search acronyms'), {
      target: { value: 'API' },
    });
    fireEvent.click(screen.getByLabelText('Submit acronym search'));
    await screen.findByText(/Application Programming Interface/);
    fireEvent.click(screen.getByText('Simulate creation'));

    expect(
      await screen.findByText(/Active Pharmaceutical Ingredient/),
    ).toBeInTheDocument();
    expect(screen.getByText('Healthcare')).toBeInTheDocument();
    expect(
      mockedGet.mock.calls.filter(([url]) => url === '/api/search'),
    ).toHaveLength(1);
  });
});
