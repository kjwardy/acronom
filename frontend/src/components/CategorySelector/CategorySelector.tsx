import React, { useState } from 'react';
import Button from '@material-ui/core/Button';
import Checkbox from '@material-ui/core/Checkbox';
import Chip from '@material-ui/core/Chip';
import FormControl from '@material-ui/core/FormControl';
import FormHelperText from '@material-ui/core/FormHelperText';
import InputLabel from '@material-ui/core/InputLabel';
import ListItemText from '@material-ui/core/ListItemText';
import MenuItem from '@material-ui/core/MenuItem';
import Select from '@material-ui/core/Select';
import { makeStyles } from '@material-ui/core/styles';
import CreateCategoryModal from '../CreateCategoryModal';
import {
  Category,
  getCategoryTextColour,
  useCategories,
} from '../../contexts/Categories';

interface CategorySelectorProps {
  disabled?: boolean;
  selected: number[];
  onChange: (categoryIDs: number[]) => void;
}

export const normalizeCategorySelection = (categoryIDs: number[]): number[] =>
  Array.from(new Set(categoryIDs)).slice(0, 5);

const useStyles = makeStyles((theme) => ({
  root: {
    marginTop: theme.spacing(2),
  },
  chips: {
    display: 'flex',
    flexWrap: 'wrap',
    gap: theme.spacing(0.5),
  },
  chip: {
    fontWeight: 600,
  },
  createButton: {
    marginTop: theme.spacing(1),
  },
}));

const CategorySelector: React.FC<CategorySelectorProps> = ({
  disabled = false,
  selected,
  onChange,
}) => {
  const classes = useStyles();
  const { categories, failed, loading } = useCategories();
  const [createOpen, setCreateOpen] = useState(false);

  const selectedCategories = selected
    .map((id) => categories.find((category) => category.id === id))
    .filter((category): category is Category => Boolean(category));

  const remove = (id: number) => {
    onChange(selected.filter((categoryID) => categoryID !== id));
  };

  const created = (category: Category) => {
    setCreateOpen(false);
    onChange(normalizeCategorySelection([...selected, category.id]));
  };

  return (
    <>
      {createOpen && (
        <CreateCategoryModal
          onClose={() => setCreateOpen(false)}
          onCreated={created}
        />
      )}
      <FormControl className={classes.root} fullWidth variant="outlined">
        <InputLabel id="categories-label">Categories</InputLabel>
        <Select
          disabled={disabled || loading || failed}
          inputProps={{ 'aria-label': 'Categories' }}
          label="Categories"
          labelId="categories-label"
          multiple
          onChange={(event) =>
            onChange(normalizeCategorySelection(event.target.value as number[]))
          }
          renderValue={() => (
            <div className={classes.chips}>
              {selectedCategories.map((category) => (
                <Chip
                  className={classes.chip}
                  key={category.id}
                  label={category.name}
                  onDelete={() => remove(category.id)}
                  onMouseDown={(event) => event.stopPropagation()}
                  size="small"
                  style={{
                    backgroundColor: category.colour,
                    color: getCategoryTextColour(category.colour),
                  }}
                />
              ))}
            </div>
          )}
          value={selected}
        >
          {categories.map((category) => {
            const checked = selected.includes(category.id);
            return (
              <MenuItem
                disabled={!checked && selected.length >= 5}
                key={category.id}
                value={category.id}
              >
                <Checkbox checked={checked} />
                <ListItemText primary={category.name} />
                <Chip
                  aria-hidden="true"
                  label={category.name}
                  size="small"
                  style={{
                    backgroundColor: category.colour,
                    color: getCategoryTextColour(category.colour),
                  }}
                />
              </MenuItem>
            );
          })}
        </Select>
        <FormHelperText error={failed}>
          {loading && 'Loading categories…'}
          {failed &&
            'Categories are unavailable. You can still save this entry.'}
          {!loading && !failed && `${selected.length} of 5 selected`}
        </FormHelperText>
      </FormControl>
      <Button
        className={classes.createButton}
        color="primary"
        disabled={disabled || failed || selected.length >= 5}
        onClick={() => setCreateOpen(true)}
        type="button"
      >
        Create category
      </Button>
    </>
  );
};

export default CategorySelector;
