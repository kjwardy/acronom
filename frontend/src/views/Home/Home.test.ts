import axios from 'axios';
import {
  isSafeResultLink,
  removeAcronymResult,
  replaceAcronymResult,
  searchAcronyms,
} from './Home';

jest.mock('axios');

const mockedGet = axios.get as jest.MockedFunction<typeof axios.get>;

describe('isSafeResultLink', () => {
  it('only allows HTTP and HTTPS links', () => {
    expect(isSafeResultLink('https://example.com')).toBe(true);
    expect(isSafeResultLink('http://example.com')).toBe(true);
    expect(isSafeResultLink(['javascript', ':alert(1)'].join(''))).toBe(false);
    expect(isSafeResultLink('not a URL')).toBe(false);
  });
});

describe('replaceAcronymResult', () => {
  it('replaces only the edited entry by ID', () => {
    const results = [
      { id: 1, acronym: 'PC', definition: 'Personal Computer' },
      { id: 2, acronym: 'PC', definition: 'Probable Cause' },
    ];
    const updated = {
      id: 2,
      acronym: 'PC',
      definition: 'Police Constable',
      updated_at: '2026-09-27T20:00:00Z',
    };

    expect(replaceAcronymResult(results, updated)).toEqual([
      results[0],
      updated,
    ]);
  });
});

describe('removeAcronymResult', () => {
  it('removes only the deleted entry by ID', () => {
    const results = [
      { id: 1, acronym: 'PC', definition: 'Personal Computer' },
      { id: 2, acronym: 'PC', definition: 'Probable Cause' },
    ];

    expect(removeAcronymResult(results, 1)).toEqual([results[1]]);
  });
});

describe('searchAcronyms', () => {
  afterEach(() => jest.clearAllMocks());

  it('calls the search API with a trimmed query', async () => {
    const results = [
      {
        id: 1,
        acronym: 'API',
        definition: 'Application Programming Interface',
      },
    ];
    mockedGet.mockResolvedValue({ data: results });

    await expect(searchAcronyms('  API  ')).resolves.toEqual(results);
    expect(mockedGet).toHaveBeenCalledWith('/api/search', {
      params: { q: 'API' },
    });
  });

  it('returns an empty result set', async () => {
    mockedGet.mockResolvedValue({ data: [] });

    await expect(searchAcronyms('missing')).resolves.toEqual([]);
  });

  it('rejects a non-API response', async () => {
    mockedGet.mockResolvedValue({ data: '<html></html>' });

    await expect(searchAcronyms('API')).rejects.toThrow(
      'Invalid search response',
    );
  });

  it('propagates failures for the page to display', async () => {
    const error = new Error('request failed');
    mockedGet.mockRejectedValue(error);

    await expect(searchAcronyms('API')).rejects.toBe(error);
  });
});
