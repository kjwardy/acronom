import React, { FormEvent, useState } from 'react';
import Button from '@material-ui/core/Button';
import CircularProgress from '@material-ui/core/CircularProgress';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import TextField from '@material-ui/core/TextField';
import Typography from '@material-ui/core/Typography';
import { makeStyles } from '@material-ui/core/styles';
import {
  CATEGORY_COLOURS,
  Category,
  getCategoryTextColour,
  useCategories,
} from '../../contexts/Categories';

interface CreateCategoryModalProps {
  onClose: () => void;
  onCreated: (category: Category) => void;
}

const useStyles = makeStyles((theme) => ({
  field: {
    marginTop: theme.spacing(2),
  },
  palette: {
    display: 'grid',
    gridTemplateColumns: 'repeat(2, minmax(0, 1fr))',
    gap: theme.spacing(1),
    marginTop: theme.spacing(2),
  },
  colourButton: {
    justifyContent: 'flex-start',
  },
  swatch: {
    width: 20,
    height: 20,
    marginRight: theme.spacing(1),
    border: '1px solid rgba(0, 0, 0, 0.2)',
    borderRadius: '50%',
  },
  error: {
    marginTop: theme.spacing(2),
    color: theme.palette.error.main,
  },
  progress: {
    marginRight: theme.spacing(1),
  },
}));

export const validateCategoryName = (name: string): string => {
  const normalized = name.trim().replace(/\s+/g, ' ');
  if (!normalized) return 'Category name is required';
  if ([...normalized].length > 30)
    return 'Category name must be 30 characters or fewer';
  return '';
};

const CreateCategoryModal: React.FC<CreateCategoryModalProps> = ({
  onClose,
  onCreated,
}) => {
  const classes = useStyles();
  const { createCategory } = useCategories();
  const [name, setName] = useState('');
  const [colour, setColour] = useState(CATEGORY_COLOURS[0].value);
  const [nameError, setNameError] = useState('');
  const [requestError, setRequestError] = useState('');
  const [submitting, setSubmitting] = useState(false);

  const submit = async (event: FormEvent) => {
    event.preventDefault();
    const validationError = validateCategoryName(name);
    setNameError(validationError);
    setRequestError('');
    if (validationError) return;

    setSubmitting(true);
    try {
      onCreated(await createCategory(name, colour));
    } catch (error) {
      const response = (error as any).response?.data;
      setRequestError(
        response?.message ||
          (typeof response === 'string' ? response : '') ||
          'Unable to create the category. Please try again.',
      );
      setSubmitting(false);
    }
  };

  return (
    <Dialog
      aria-labelledby="create-category-title"
      fullWidth
      maxWidth="xs"
      onClose={submitting ? undefined : onClose}
      open
    >
      <form noValidate onSubmit={submit}>
        <DialogTitle id="create-category-title">Create category</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Choose a reusable category name and colour.
          </DialogContentText>
          <TextField
            autoFocus
            className={classes.field}
            disabled={submitting}
            error={Boolean(nameError)}
            fullWidth
            helperText={nameError || 'Maximum 30 characters'}
            id="category-name"
            inputProps={{ maxLength: 30 }}
            label="Category name"
            onChange={(event) => setName(event.target.value)}
            required
            value={name}
            variant="outlined"
          />
          <div
            aria-label="Category colour"
            className={classes.palette}
            role="group"
          >
            {CATEGORY_COLOURS.map((option) => (
              <Button
                aria-label={`${option.name} category colour`}
                aria-pressed={colour === option.value}
                className={classes.colourButton}
                disabled={submitting}
                key={option.value}
                onClick={() => setColour(option.value)}
                style={
                  colour === option.value
                    ? {
                        backgroundColor: option.value,
                        color: getCategoryTextColour(option.value),
                      }
                    : undefined
                }
                type="button"
                variant={colour === option.value ? 'contained' : 'outlined'}
              >
                <span
                  aria-hidden="true"
                  className={classes.swatch}
                  style={{ backgroundColor: option.value }}
                />
                {option.name}
              </Button>
            ))}
          </div>
          {requestError && (
            <Typography className={classes.error} role="alert" variant="body2">
              {requestError}
            </Typography>
          )}
        </DialogContent>
        <DialogActions>
          <Button disabled={submitting} onClick={onClose} type="button">
            Cancel
          </Button>
          <Button
            color="primary"
            disabled={submitting}
            type="submit"
            variant="contained"
          >
            {submitting && (
              <CircularProgress
                className={classes.progress}
                color="inherit"
                size={16}
              />
            )}
            {submitting ? 'Creating…' : 'Create category'}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  );
};

export default CreateCategoryModal;
