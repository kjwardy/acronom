import axios from 'axios';
import {
  fetchCategories,
  getCategoryTextColour,
  postCategory,
} from './Categories';

jest.mock('axios');

const mockedGet = axios.get as jest.MockedFunction<typeof axios.get>;
const mockedPost = axios.post as jest.MockedFunction<typeof axios.post>;

describe('category API helpers', () => {
  afterEach(() => jest.clearAllMocks());

  it('loads categories', async () => {
    const categories = [{ id: 1, name: 'Platform', colour: '#1976D2' }];
    mockedGet.mockResolvedValue({ data: categories });

    await expect(fetchCategories()).resolves.toEqual(categories);
    expect(mockedGet).toHaveBeenCalledWith('/api/categories');
  });

  it('creates a trimmed category', async () => {
    const category = { id: 1, name: 'Platform', colour: '#1976D2' };
    mockedPost.mockResolvedValue({ data: category });

    await expect(postCategory(' Platform ', '#1976D2')).resolves.toEqual(
      category,
    );
    expect(mockedPost).toHaveBeenCalledWith('/api/categories', {
      name: 'Platform',
      colour: '#1976D2',
    });
  });
});

describe('getCategoryTextColour', () => {
  it('selects contrasting text colours', () => {
    expect(getCategoryTextColour('#FFA000')).toBe('#111111');
    expect(getCategoryTextColour('#F57C00')).toBe('#111111');
    expect(getCategoryTextColour('#388E3C')).toBe('#111111');
    expect(getCategoryTextColour('#455A64')).toBe('#FFFFFF');
  });
});
