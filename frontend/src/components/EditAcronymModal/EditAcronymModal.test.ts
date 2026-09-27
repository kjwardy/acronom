import axios from 'axios';
import {
  deleteAcronym,
  updateAcronym,
  validateAcronymEdit,
} from './EditAcronymModal';

jest.mock('axios');

const mockedDelete = axios.delete as jest.MockedFunction<typeof axios.delete>;
const mockedPut = axios.put as jest.MockedFunction<typeof axios.put>;

describe('validateAcronymEdit', () => {
  it('requires a definition', () => {
    expect(validateAcronymEdit(' ', '')).toEqual({
      definition: 'Definition is required',
    });
  });

  it('accepts a definition and optional HTTP URL', () => {
    expect(validateAcronymEdit('Personal Computer', '')).toEqual({});
    expect(
      validateAcronymEdit('Personal Computer', 'https://example.com'),
    ).toEqual({});
  });

  it('rejects an invalid URL', () => {
    expect(validateAcronymEdit('Personal Computer', 'example.com')).toEqual({
      link: 'Enter a valid HTTP or HTTPS URL',
    });
  });
});

describe('updateAcronym', () => {
  afterEach(() => jest.clearAllMocks());

  it('puts trimmed editable fields to the entry endpoint', async () => {
    const entry = {
      id: 7,
      acronym: 'PC',
      definition: 'Personal Computer',
      link: 'https://old.example.com',
    };
    const updated = {
      ...entry,
      definition: 'Probable Cause',
      link: 'https://example.com',
      updated_at: '2026-09-27T20:00:00Z',
    };
    mockedPut.mockResolvedValue({ data: updated });

    await expect(
      updateAcronym(entry, ' Probable Cause ', ' https://example.com '),
    ).resolves.toEqual(updated);
    expect(mockedPut).toHaveBeenCalledWith('/api/acronyms/7', {
      definition: 'Probable Cause',
      link: 'https://example.com',
    });
  });

  it('sends null to clear an optional URL', async () => {
    const entry = { id: 7, acronym: 'PC', definition: 'Personal Computer' };
    mockedPut.mockResolvedValue({ data: entry });

    await updateAcronym(entry, 'Personal Computer', ' ');
    expect(mockedPut).toHaveBeenCalledWith('/api/acronyms/7', {
      definition: 'Personal Computer',
      link: null,
    });
  });

  it('propagates API failures for the modal to display', async () => {
    const error = new Error('request failed');
    mockedPut.mockRejectedValue(error);

    await expect(
      updateAcronym(
        { id: 7, acronym: 'PC', definition: 'Personal Computer' },
        'Probable Cause',
        '',
      ),
    ).rejects.toBe(error);
  });
});

describe('deleteAcronym', () => {
  afterEach(() => jest.clearAllMocks());

  it('deletes the entry by ID', async () => {
    mockedDelete.mockResolvedValue({});

    await deleteAcronym(7);
    expect(mockedDelete).toHaveBeenCalledWith('/api/acronyms/7');
  });

  it('propagates API failures for the confirmation dialog to display', async () => {
    const error = new Error('request failed');
    mockedDelete.mockRejectedValue(error);

    await expect(deleteAcronym(7)).rejects.toBe(error);
  });
});
