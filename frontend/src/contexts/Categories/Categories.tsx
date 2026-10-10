import React, { createContext, useContext, useEffect, useState } from 'react';
import axios from 'axios';

export interface Category {
  id: number;
  name: string;
  colour: string;
  created_at?: string;
  updated_at?: string;
}

export interface CategoryColour {
  name: string;
  value: string;
}

export const CATEGORY_COLOURS: CategoryColour[] = [
  { name: 'Blue', value: '#1976D2' },
  { name: 'Teal', value: '#00796B' },
  { name: 'Green', value: '#388E3C' },
  { name: 'Amber', value: '#FFA000' },
  { name: 'Orange', value: '#F57C00' },
  { name: 'Red', value: '#D32F2F' },
  { name: 'Purple', value: '#7B1FA2' },
  { name: 'Slate', value: '#455A64' },
];

interface CategoriesContextValue {
  categories: Category[];
  loading: boolean;
  failed: boolean;
  createCategory: (name: string, colour: string) => Promise<Category>;
}

const CategoriesContext = createContext<CategoriesContextValue | undefined>(
  undefined,
);

export const fetchCategories = async (): Promise<Category[]> => {
  const { data } = await axios.get<Category[]>('/api/categories');
  if (!Array.isArray(data)) throw new Error('Invalid categories response');
  return data;
};

export const postCategory = async (
  name: string,
  colour: string,
): Promise<Category> => {
  const { data } = await axios.post<Category>('/api/categories', {
    name: name.trim(),
    colour,
  });
  return data;
};

const relativeLuminance = (colour: string): number => {
  const hex = colour.replace('#', '');
  const channels = [0, 2, 4].map((offset) => {
    const channel = parseInt(hex.slice(offset, offset + 2), 16) / 255;
    return channel <= 0.04045
      ? channel / 12.92
      : ((channel + 0.055) / 1.055) ** 2.4;
  });
  return channels[0] * 0.2126 + channels[1] * 0.7152 + channels[2] * 0.0722;
};

const contrastRatio = (first: string, second: string): number => {
  const lighter = Math.max(relativeLuminance(first), relativeLuminance(second));
  const darker = Math.min(relativeLuminance(first), relativeLuminance(second));
  return (lighter + 0.05) / (darker + 0.05);
};

export const getCategoryTextColour = (colour: string): string =>
  contrastRatio(colour, '#111111') >= contrastRatio(colour, '#FFFFFF')
    ? '#111111'
    : '#FFFFFF';

export const CategoriesProvider: React.FC = ({ children }) => {
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);
  const [failed, setFailed] = useState(false);

  useEffect(() => {
    let mounted = true;
    fetchCategories()
      .then((loaded) => {
        if (mounted) setCategories(loaded);
      })
      .catch(() => {
        if (mounted) setFailed(true);
      })
      .finally(() => {
        if (mounted) setLoading(false);
      });
    return () => {
      mounted = false;
    };
  }, []);

  const createCategory = async (name: string, colour: string) => {
    const category = await postCategory(name, colour);
    setCategories((current) =>
      [...current, category].sort((first, second) =>
        first.name.localeCompare(second.name),
      ),
    );
    return category;
  };

  return (
    <CategoriesContext.Provider
      value={{ categories, loading, failed, createCategory }}
    >
      {children}
    </CategoriesContext.Provider>
  );
};

export const useCategories = (): CategoriesContextValue => {
  const context = useContext(CategoriesContext);
  if (!context)
    throw new Error('useCategories must be used within CategoriesProvider');
  return context;
};
