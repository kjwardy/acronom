import { validateCategoryName } from './CreateCategoryModal';

describe('validateCategoryName', () => {
  it('requires a name', () => {
    expect(validateCategoryName('  ')).toBe('Category name is required');
  });

  it('enforces the 30-character limit after whitespace normalization', () => {
    expect(validateCategoryName('a'.repeat(31))).toBe(
      'Category name must be 30 characters or fewer',
    );
    expect(validateCategoryName('  Platform   Tools  ')).toBe('');
  });
});
