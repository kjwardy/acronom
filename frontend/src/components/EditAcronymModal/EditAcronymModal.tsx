import React, { FormEvent, useState } from 'react';
import axios from 'axios';
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

export interface EditableAcronym {
  id: number;
  acronym: string;
  definition: string;
  link?: string;
  created_at?: string;
  updated_at?: string;
}

interface EditAcronymModalProps {
  entry: EditableAcronym;
  onClose: () => void;
  onDeleted: (id: number) => void;
  onUpdated: (entry: EditableAcronym) => void;
}

interface FormErrors {
  definition?: string;
  link?: string;
}

const useStyles = makeStyles((theme) => ({
  field: {
    marginTop: theme.spacing(2),
  },
  error: {
    marginTop: theme.spacing(2),
    color: theme.palette.error.main,
  },
  progress: {
    marginRight: theme.spacing(1),
  },
  editDeleteAction: {
    marginRight: 'auto',
    color: theme.palette.error.main,
  },
  destructiveAction: {
    color: theme.palette.error.main,
  },
}));

export const validateAcronymEdit = (
  definition: string,
  link: string,
): FormErrors => {
  const errors: FormErrors = {};
  if (!definition.trim()) errors.definition = 'Definition is required';
  if (link.trim()) {
    try {
      const url = new URL(link.trim());
      if (url.protocol !== 'http:' && url.protocol !== 'https:') {
        errors.link = 'Enter a valid HTTP or HTTPS URL';
      }
    } catch {
      errors.link = 'Enter a valid HTTP or HTTPS URL';
    }
  }
  return errors;
};

export const updateAcronym = async (
  entry: EditableAcronym,
  definition: string,
  link: string,
): Promise<EditableAcronym> => {
  const payload = {
    definition: definition.trim(),
    link: link.trim() || null,
  };
  const { data } = await axios.put<EditableAcronym>(
    `/api/acronyms/${entry.id}`,
    payload,
  );
  return data;
};

export const deleteAcronym = async (id: number): Promise<void> => {
  await axios.delete(`/api/acronyms/${id}`);
};

const EditAcronymModal: React.FC<EditAcronymModalProps> = ({
  entry,
  onClose,
  onDeleted,
  onUpdated,
}) => {
  const classes = useStyles();
  const [definition, setDefinition] = useState(entry.definition);
  const [link, setLink] = useState(entry.link || '');
  const [errors, setErrors] = useState<FormErrors>({});
  const [requestError, setRequestError] = useState('');
  const [deleteError, setDeleteError] = useState('');
  const [confirmDelete, setConfirmDelete] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [deleting, setDeleting] = useState(false);

  const submit = async (event: FormEvent) => {
    event.preventDefault();
    const validationErrors = validateAcronymEdit(definition, link);
    setErrors(validationErrors);
    setRequestError('');
    if (Object.keys(validationErrors).length > 0) return;

    setSubmitting(true);
    try {
      onUpdated(await updateAcronym(entry, definition, link));
    } catch (error) {
      const response = (error as any).response?.data;
      setRequestError(
        response?.message ||
          (typeof response === 'string' ? response : '') ||
          'Unable to update the acronym. Please try again.',
      );
      setSubmitting(false);
    }
  };

  const remove = async () => {
    setDeleting(true);
    setDeleteError('');
    try {
      await deleteAcronym(entry.id);
      onDeleted(entry.id);
    } catch (error) {
      const response = (error as any).response?.data;
      setDeleteError(
        response?.message ||
          (typeof response === 'string' ? response : '') ||
          'Unable to delete the acronym. Please try again.',
      );
      setDeleting(false);
    }
  };

  const closeConfirmation = () => {
    setConfirmDelete(false);
    setDeleteError('');
  };

  return (
    <>
      <Dialog
        aria-labelledby="edit-acronym-title"
        fullWidth
        maxWidth="sm"
        onClose={submitting || deleting ? undefined : onClose}
        open
      >
        <form noValidate onSubmit={submit}>
          <DialogTitle id="edit-acronym-title">Edit acronym</DialogTitle>
          <DialogContent>
            <DialogContentText>
              Update this acronym definition and its optional reference URL.
            </DialogContentText>
            <TextField
              className={classes.field}
              disabled
              fullWidth
              id="edit-acronym"
              label="Acronym"
              value={entry.acronym}
              variant="outlined"
            />
            <TextField
              autoFocus
              className={classes.field}
              disabled={submitting || deleting}
              error={Boolean(errors.definition)}
              fullWidth
              helperText={errors.definition}
              id="edit-definition"
              label="Definition"
              multiline
              onChange={(event) => setDefinition(event.target.value)}
              required
              rows={3}
              value={definition}
              variant="outlined"
            />
            <TextField
              className={classes.field}
              disabled={submitting || deleting}
              error={Boolean(errors.link)}
              fullWidth
              helperText={errors.link || 'Optional'}
              id="edit-link"
              inputProps={{ inputMode: 'url' }}
              label="URL"
              onChange={(event) => setLink(event.target.value)}
              placeholder="https://example.com"
              value={link}
              variant="outlined"
            />
            {requestError && (
              <Typography
                className={classes.error}
                role="alert"
                variant="body2"
              >
                {requestError}
              </Typography>
            )}
          </DialogContent>
          <DialogActions>
            <Button
              className={classes.editDeleteAction}
              disabled={submitting || deleting}
              onClick={() => setConfirmDelete(true)}
              type="button"
            >
              Delete
            </Button>
            <Button
              disabled={submitting || deleting}
              onClick={onClose}
              type="button"
            >
              Cancel
            </Button>
            <Button
              color="primary"
              disabled={submitting || deleting}
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
              {submitting ? 'Saving…' : 'Save changes'}
            </Button>
          </DialogActions>
        </form>
      </Dialog>
      <Dialog
        aria-describedby="delete-acronym-description"
        aria-labelledby="delete-acronym-title"
        onClose={deleting ? undefined : closeConfirmation}
        open={confirmDelete}
      >
        <DialogTitle id="delete-acronym-title">Delete acronym?</DialogTitle>
        <DialogContent>
          <DialogContentText id="delete-acronym-description">
            Delete “{entry.acronym} - {entry.definition}”? This action cannot be
            undone.
          </DialogContentText>
          {deleteError && (
            <Typography className={classes.error} role="alert" variant="body2">
              {deleteError}
            </Typography>
          )}
        </DialogContent>
        <DialogActions>
          <Button disabled={deleting} onClick={closeConfirmation} type="button">
            Cancel
          </Button>
          <Button
            className={classes.destructiveAction}
            disabled={deleting}
            onClick={remove}
            type="button"
          >
            {deleting && (
              <CircularProgress
                className={classes.progress}
                color="inherit"
                size={16}
              />
            )}
            Confirm Delete
          </Button>
        </DialogActions>
      </Dialog>
    </>
  );
};

export default EditAcronymModal;
