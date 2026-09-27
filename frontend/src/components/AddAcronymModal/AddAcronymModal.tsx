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

interface AcronymEntry {
  id: number;
  acronym: string;
  definition: string;
  link?: string;
}

interface AddAcronymModalProps {
  onClose: () => void;
  onCreated: (entry: AcronymEntry) => void;
}

interface FormErrors {
  acronym?: string;
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
}));

export const validateAcronym = (
  acronym: string,
  definition: string,
  link: string,
): FormErrors => {
  const errors: FormErrors = {};
  if (!acronym.trim()) errors.acronym = 'Acronym is required';
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

export const createAcronym = async (
  acronym: string,
  definition: string,
  link: string,
): Promise<AcronymEntry> => {
  const payload = {
    acronym: acronym.trim(),
    definition: definition.trim(),
    ...(link.trim() ? { link: link.trim() } : {}),
  };
  const { data } = await axios.post<AcronymEntry>('/api/acronyms', payload);
  return data;
};

const AddAcronymModal: React.FC<AddAcronymModalProps> = ({
  onClose,
  onCreated,
}) => {
  const classes = useStyles();
  const [acronym, setAcronym] = useState('');
  const [definition, setDefinition] = useState('');
  const [link, setLink] = useState('');
  const [errors, setErrors] = useState<FormErrors>({});
  const [requestError, setRequestError] = useState('');
  const [submitting, setSubmitting] = useState(false);

  const submit = async (event: FormEvent) => {
    event.preventDefault();
    const validationErrors = validateAcronym(acronym, definition, link);
    setErrors(validationErrors);
    setRequestError('');
    if (Object.keys(validationErrors).length > 0) return;

    setSubmitting(true);
    try {
      const entry = await createAcronym(acronym, definition, link);
      onCreated(entry);
    } catch (error) {
      const response = (error as any).response?.data;
      setRequestError(
        response?.message ||
          (typeof response === 'string' ? response : '') ||
          'Unable to add the acronym. Please try again.',
      );
      setSubmitting(false);
    }
  };

  return (
    <Dialog
      aria-labelledby="add-acronym-title"
      fullWidth
      maxWidth="sm"
      onClose={submitting ? undefined : onClose}
      open
    >
      <form noValidate onSubmit={submit}>
        <DialogTitle id="add-acronym-title">Add new acronym</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Add an acronym and its definition to the shared glossary.
          </DialogContentText>
          <TextField
            autoComplete="off"
            autoFocus
            className={classes.field}
            disabled={submitting}
            error={Boolean(errors.acronym)}
            fullWidth
            helperText={errors.acronym}
            id="acronym"
            label="Acronym"
            onChange={(event) => setAcronym(event.target.value)}
            required
            value={acronym}
            variant="outlined"
          />
          <TextField
            className={classes.field}
            disabled={submitting}
            error={Boolean(errors.definition)}
            fullWidth
            helperText={errors.definition}
            id="definition"
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
            disabled={submitting}
            error={Boolean(errors.link)}
            fullWidth
            helperText={errors.link || 'Optional'}
            id="link"
            inputProps={{ inputMode: 'url' }}
            label="URL"
            onChange={(event) => setLink(event.target.value)}
            placeholder="https://example.com"
            value={link}
            variant="outlined"
          />
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
            {submitting ? 'Adding…' : 'Add acronym'}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  );
};

export default AddAcronymModal;
