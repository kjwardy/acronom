import { normalizeCategorySelection } from './CategorySelector';

describe('normalizeCategorySelection', () => {
  it('deduplicates IDs and preserves selection order', () => {
    expect(normalizeCategorySelection([3, 1, 3, 2])).toEqual([3, 1, 2]);
  });

  it('limits selections to five categories', () => {
    expect(normalizeCategorySelection([1, 2, 3, 4, 5, 6])).toEqual([
      1, 2, 3, 4, 5,
    ]);
  });
});
