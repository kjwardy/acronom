import React from 'react';
import axios from 'axios';
import '@testing-library/jest-dom';
import { fireEvent, render, screen } from '@testing-library/react';
import AddAcronymModal from '../AddAcronymModal';
import EditAcronymModal from '../EditAcronymModal';
import { CategoriesProvider } from '../../contexts/Categories';

jest.mock('axios');

const mockedGet = axios.get as jest.MockedFunction<typeof axios.get>;

const platform = { id: 1, name: 'Platform', colour: '#1976D2' };

const renderWithCategories = (component: React.ReactElement) => {
  mockedGet.mockResolvedValue({ data: [platform] });
  return render(<CategoriesProvider>{component}</CategoriesProvider>);
};

describe('acronym category forms', () => {
  afterEach(() => jest.clearAllMocks());

  it('preserves add form values after cancelling category creation', async () => {
    renderWithCategories(
      <AddAcronymModal onClose={jest.fn()} onCreated={jest.fn()} />,
    );
    await screen.findByText('0 of 5 selected');
    fireEvent.change(screen.getByLabelText('Acronym *'), {
      target: { value: 'API' },
    });
    fireEvent.change(screen.getByLabelText('Definition *'), {
      target: { value: 'Application Programming Interface' },
    });

    fireEvent.click(screen.getByRole('button', { name: 'Create category' }));
    fireEvent.click(screen.getByRole('button', { name: 'Cancel' }));

    expect(screen.getByLabelText('Acronym *')).toHaveValue('API');
    expect(screen.getByLabelText('Definition *')).toHaveValue(
      'Application Programming Interface',
    );
  });

  it('pre-populates and removes edit form category assignments', async () => {
    renderWithCategories(
      <EditAcronymModal
        entry={{
          id: 1,
          acronym: 'API',
          definition: 'Application Programming Interface',
          categories: [platform],
        }}
        onClose={jest.fn()}
        onDeleted={jest.fn()}
        onUpdated={jest.fn()}
      />,
    );

    expect(await screen.findByText('1 of 5 selected')).toBeInTheDocument();
    const chip = screen.getByRole('button', { name: 'Platform' });
    fireEvent.click(chip.querySelector('svg')!);
    expect(await screen.findByText('0 of 5 selected')).toBeInTheDocument();
  });
});
