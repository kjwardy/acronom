import axios from 'axios';
import {
  adjustMetricTotal,
  fetchTotalAcronyms,
  formatAcronymCount,
} from './Metrics';

jest.mock('axios');

const mockedGet = axios.get as jest.MockedFunction<typeof axios.get>;

describe('fetchTotalAcronyms', () => {
  afterEach(() => jest.clearAllMocks());

  it('loads the total with a three-second timeout', async () => {
    mockedGet.mockResolvedValue({ data: { total_acronyms: 42 } });

    await expect(fetchTotalAcronyms()).resolves.toBe(42);
    expect(mockedGet).toHaveBeenCalledWith('/api/metrics/total-acronyms', {
      timeout: 3000,
    });
  });

  it('accepts an empty database count', async () => {
    mockedGet.mockResolvedValue({ data: { total_acronyms: 0 } });

    await expect(fetchTotalAcronyms()).resolves.toBe(0);
  });

  it('rejects request and malformed response failures', async () => {
    mockedGet.mockRejectedValueOnce(new Error('request failed'));
    await expect(fetchTotalAcronyms()).rejects.toThrow('request failed');

    mockedGet.mockResolvedValueOnce({ data: { total_acronyms: -1 } });
    await expect(fetchTotalAcronyms()).rejects.toThrow(
      'Invalid metrics response',
    );
  });
});

describe('formatAcronymCount', () => {
  it('formats zero, singular, and plural counts', () => {
    expect(formatAcronymCount(0)).toBe('0 meanings decoded');
    expect(formatAcronymCount(1)).toBe('1 meaning decoded');
    expect(formatAcronymCount(2)).toBe('2 meanings decoded');
  });
});

describe('adjustMetricTotal', () => {
  it('adjusts creates and deletes without falling below zero', () => {
    expect(adjustMetricTotal(2, 1)).toBe(3);
    expect(adjustMetricTotal(2, -1)).toBe(1);
    expect(adjustMetricTotal(0, -1)).toBe(0);
  });
});
