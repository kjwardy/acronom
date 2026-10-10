import React from 'react';
import axios from 'axios';
import '@testing-library/jest-dom';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import CategorySelector from './CategorySelector';
import { CategoriesProvider } from '../../contexts/Categories';

jest.mock('axios');

const mockedGet = axios.get as jest.MockedFunction<typeof axios.get>;
const mockedPost = axios.post as jest.MockedFunction<typeof axios.post>;

const categories = [
  { id: 1, name: 'Business', colour: '#388E3C' },
  { id: 2, name: 'Platform', colour: '#1976D2' },
  { id: 3, name: 'Security', colour: '#D32F2F' },
  { id: 4, name: 'Support', colour: '#FFA000' },
  { id: 5, name: 'Technology', colour: '#7B1FA2' },
  { id: 6, name: 'Testing', colour: '#455A64' },
];

const SelectorHarness: React.FC<{ initial?: number[] }> = ({
  initial = [],
}) => {
  const [selected, setSelected] = React.useState(initial);
  return <CategorySelector onChange={setSelected} selected={selected} />;
};

const renderSelector = (initial: number[] = []) =>
  render(
    <CategoriesProvider>
      <SelectorHarness initial={initial} />
    </CategoriesProvider>,
  );

describe('CategorySelector interactions', () => {
  afterEach(() => jest.clearAllMocks());

  it('loads selections, removes chips, and prevents a sixth category', async () => {
    mockedGet.mockResolvedValue({ data: categories });
    renderSelector([1, 2, 3, 4, 5]);

    expect(await screen.findByText('5 of 5 selected')).toBeInTheDocument();
    expect(screen.getByText('Business')).toBeInTheDocument();
    expect(
      screen.getByRole('button', { name: 'Create category' }),
    ).toBeDisabled();

    fireEvent.mouseDown(screen.getByLabelText('Categories'));
    expect(screen.getByRole('option', { name: /Testing/ })).toHaveAttribute(
      'aria-disabled',
      'true',
    );

    fireEvent.keyDown(screen.getByRole('listbox'), { key: 'Escape' });
    fireEvent.click(
      screen.getByRole('button', { name: 'Business' }).querySelector('svg')!,
    );
    expect(await screen.findByText('4 of 5 selected')).toBeInTheDocument();
  });

  it('creates and automatically selects a category without clearing selections', async () => {
    mockedGet.mockResolvedValue({ data: categories.slice(0, 1) });
    mockedPost.mockResolvedValue({
      data: { id: 2, name: 'Platform Tools', colour: '#1976D2' },
    });
    renderSelector([1]);

    await screen.findByText('1 of 5 selected');
    fireEvent.click(screen.getByRole('button', { name: 'Create category' }));
    fireEvent.change(screen.getByLabelText(/Category name/), {
      target: { value: 'Platform Tools' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'Create category' }));

    expect(await screen.findByText('2 of 5 selected')).toBeInTheDocument();
    expect(screen.getByText('Business')).toBeInTheDocument();
    expect(screen.getByText('Platform Tools')).toBeInTheDocument();
  });

  it('preserves category input and reports creation failures', async () => {
    mockedGet.mockResolvedValue({ data: [] });
    mockedPost.mockRejectedValue(new Error('request failed'));
    renderSelector();

    await screen.findByText('0 of 5 selected');
    fireEvent.click(screen.getByRole('button', { name: 'Create category' }));
    const input = screen.getByLabelText(/Category name/);
    fireEvent.change(input, { target: { value: 'Platform' } });
    fireEvent.click(screen.getByRole('button', { name: 'Create category' }));

    expect(await screen.findByRole('alert')).toHaveTextContent(
      'Unable to create the category. Please try again.',
    );
    expect(input).toHaveValue('Platform');
  });

  it('shows a non-blocking category loading failure', async () => {
    mockedGet.mockRejectedValue(new Error('request failed'));
    renderSelector();

    expect(
      await screen.findByText(
        'Categories are unavailable. You can still save this entry.',
      ),
    ).toBeInTheDocument();
    expect(screen.getByLabelText('Categories')).toHaveAttribute(
      'aria-disabled',
      'true',
    );
    await waitFor(() =>
      expect(
        screen.getByRole('button', { name: 'Create category' }),
      ).toBeDisabled(),
    );
  });
});
