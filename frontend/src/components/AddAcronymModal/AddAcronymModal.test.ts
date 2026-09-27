import axios from 'axios';
import { createAcronym, validateAcronym } from './AddAcronymModal';

jest.mock('axios');

const mockedPost = axios.post as jest.MockedFunction<typeof axios.post>;

describe('validateAcronym', () => {
  it('requires an acronym and definition', () => {
    expect(validateAcronym(' ', '', '')).toEqual({
      acronym: 'Acronym is required',
      definition: 'Definition is required',
    });
  });

  it('allows a valid entry without a URL', () => {
    expect(
      validateAcronym('API', 'Application Programming Interface', ''),
    ).toEqual({});
  });

  it('rejects invalid and non-HTTP URLs', () => {
    expect(validateAcronym('API', 'Definition', 'example.com').link).toBe(
      'Enter a valid HTTP or HTTPS URL',
    );
    expect(validateAcronym('API', 'Definition', 'ftp://example.com').link).toBe(
      'Enter a valid HTTP or HTTPS URL',
    );
  });
});

describe('createAcronym', () => {
  afterEach(() => jest.clearAllMocks());

  it('posts a trimmed payload to the acronym API', async () => {
    const entry = {
      id: 1,
      acronym: 'API',
      definition: 'Application Programming Interface',
      link: 'https://example.com',
    };
    mockedPost.mockResolvedValue({ data: entry });

    await expect(
      createAcronym(
        ' API ',
        ' Application Programming Interface ',
        ' https://example.com ',
      ),
    ).resolves.toEqual(entry);
    expect(mockedPost).toHaveBeenCalledWith('/api/acronyms', {
      acronym: 'API',
      definition: 'Application Programming Interface',
      link: 'https://example.com',
    });
  });

  it('omits an empty optional URL', async () => {
    mockedPost.mockResolvedValue({
      data: { id: 1, acronym: 'API', definition: 'Definition' },
    });

    await createAcronym('API', 'Definition', ' ');
    expect(mockedPost).toHaveBeenCalledWith('/api/acronyms', {
      acronym: 'API',
      definition: 'Definition',
    });
  });

  it('propagates API failures for the modal to display', async () => {
    const error = new Error('request failed');
    mockedPost.mockRejectedValue(error);

    await expect(createAcronym('API', 'Definition', '')).rejects.toBe(error);
  });
});
