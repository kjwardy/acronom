import React, { FormEvent, useState } from 'react';
import axios from 'axios';
import CircularProgress from '@material-ui/core/CircularProgress';
import IconButton from '@material-ui/core/IconButton';
import InputBase from '@material-ui/core/InputBase';
import Paper from '@material-ui/core/Paper';
import Snackbar from '@material-ui/core/Snackbar';
import Tooltip from '@material-ui/core/Tooltip';
import Typography from '@material-ui/core/Typography';
import EditIcon from '@material-ui/icons/Edit';
import SearchIcon from '@material-ui/icons/Search';
import { makeStyles } from '@material-ui/core/styles';
import EditAcronymModal from '../../components/EditAcronymModal';

export interface AcronymResult {
  id: number;
  acronym: string;
  definition: string;
  link?: string;
  created_at?: string;
  updated_at?: string;
}

export const searchAcronyms = async (
  query: string,
): Promise<AcronymResult[]> => {
  const { data } = await axios.get<AcronymResult[]>('/api/search', {
    params: { q: query.trim() },
  });
  if (!Array.isArray(data)) throw new Error('Invalid search response');
  return data;
};

export const isSafeResultLink = (link?: string): boolean => {
  if (!link) return false;
  try {
    const url = new URL(link);
    return url.protocol === 'http:' || url.protocol === 'https:';
  } catch {
    return false;
  }
};

export const replaceAcronymResult = (
  results: AcronymResult[],
  updated: AcronymResult,
): AcronymResult[] =>
  results.map((result) => (result.id === updated.id ? updated : result));

export const removeAcronymResult = (
  results: AcronymResult[],
  deletedID: number,
): AcronymResult[] => results.filter((result) => result.id !== deletedID);

const useStyles = makeStyles((theme) => ({
  root: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    minHeight: 'calc(100vh - 64px)',
    padding: theme.spacing(3),
  },
  rootWithResults: {
    alignItems: 'flex-start',
    paddingTop: theme.spacing(6),
  },
  content: {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    width: '100%',
  },
  title: {
    marginBottom: theme.spacing(4),
    fontWeight: 600,
  },
  logo: {
    width: 96,
    height: 96,
    marginBottom: theme.spacing(1.5),
    borderRadius: theme.shape.borderRadius,
    objectFit: 'cover',
    filter: theme.palette.type === 'light' ? 'none' : 'invert(1)',
    mixBlendMode: theme.palette.type === 'light' ? 'multiply' : 'screen',
  },
  info: {
    marginTop: theme.spacing(3),
    padding: theme.spacing(1, 2),
    border: `1px solid ${theme.palette.divider}`,
    color: theme.palette.text.secondary,
  },
  search: {
    display: 'flex',
    alignItems: 'center',
    width: '100%',
    maxWidth: 640,
    padding: theme.spacing(0.5, 1, 0.5, 2),
    border: `1px solid ${theme.palette.divider}`,
    borderRadius: 28,
    boxShadow:
      theme.palette.type === 'light'
        ? '0 8px 24px rgba(20, 29, 60, 0.12)'
        : '0 8px 24px rgba(0, 0, 0, 0.28)',
    transition: theme.transitions.create(['border-color', 'box-shadow']),
    '&:focus-within': {
      borderColor: theme.palette.primary.main,
      boxShadow:
        theme.palette.type === 'light'
          ? '0 10px 30px rgba(64, 84, 178, 0.18)'
          : '0 10px 30px rgba(0, 0, 0, 0.36)',
    },
  },
  input: {
    flex: 1,
    fontSize: '1.1rem',
  },
  searchButton: {
    color: theme.palette.text.secondary,
  },
  message: {
    width: '100%',
    maxWidth: 640,
    marginTop: theme.spacing(3),
    textAlign: 'center',
    color: theme.palette.text.secondary,
  },
  error: {
    color: theme.palette.error.main,
  },
  results: {
    display: 'grid',
    gap: theme.spacing(2),
    width: '100%',
    maxWidth: 640,
    marginTop: theme.spacing(3),
  },
  result: {
    display: 'flex',
    alignItems: 'center',
    gap: theme.spacing(1),
    padding: theme.spacing(2, 2.5),
    border: `1px solid ${theme.palette.divider}`,
  },
  resultText: {
    flex: 1,
    minWidth: 0,
  },
  resultAcronym: {
    fontWeight: 600,
  },
  resultDefinition: {
    color: theme.palette.text.secondary,
  },
  resultLink: {
    color: theme.palette.primary.main,
    textDecoration: 'none',
    '&:hover, &:focus': {
      textDecoration: 'underline',
    },
  },
  updatedToast: {
    width: 'auto',
    minWidth: 'unset',
    color: theme.palette.common.white,
    backgroundColor: '#1976d2',
    '& .MuiSnackbarContent-message': {
      textAlign: 'center',
    },
  },
  deletedToast: {
    width: 'auto',
    minWidth: 'unset',
    color: theme.palette.common.white,
    backgroundColor: '#2e7d32',
    '& .MuiSnackbarContent-message': {
      textAlign: 'center',
    },
  },
}));

const Home: React.FC = () => {
  const [query, setQuery] = useState('');
  const [results, setResults] = useState<AcronymResult[] | null>(null);
  const [editing, setEditing] = useState<AcronymResult | null>(null);
  const [successMessage, setSuccessMessage] = useState('');
  const [successType, setSuccessType] = useState<'updated' | 'deleted'>(
    'updated',
  );
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const classes = useStyles();

  const handleSubmit = async (event: FormEvent) => {
    event.preventDefault();
    const trimmedQuery = query.trim();
    if (!trimmedQuery) {
      setError('Enter an acronym to search.');
      setResults(null);
      return;
    }

    setLoading(true);
    setError('');
    setResults(null);
    try {
      setResults(await searchAcronyms(trimmedQuery));
    } catch {
      setError('Unable to search acronyms right now. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const handleUpdated = (updated: AcronymResult) => {
    setResults((current) =>
      current ? replaceAcronymResult(current, updated) : null,
    );
    setEditing(null);
    setSuccessType('updated');
    setSuccessMessage('Acronym updated successfully!');
  };

  const handleDeleted = (deletedID: number) => {
    setResults((current) =>
      current ? removeAcronymResult(current, deletedID) : null,
    );
    setSuccessType('deleted');
    setSuccessMessage('Acronym deleted successfully!');
    setEditing(null);
  };

  const closeSuccess = (_event?: React.SyntheticEvent, reason?: string) => {
    if (reason !== 'clickaway') setSuccessMessage('');
  };

  const hasSearchState = loading || Boolean(error) || results !== null;

  return (
    <div
      className={`${classes.root} ${
        hasSearchState ? classes.rootWithResults : ''
      }`}
    >
      {editing && (
        <EditAcronymModal
          entry={editing}
          onClose={() => setEditing(null)}
          onDeleted={handleDeleted}
          onUpdated={handleUpdated}
        />
      )}
      <div className={classes.content}>
        <img
          alt="Acronom logo"
          className={classes.logo}
          src={`${process.env.PUBLIC_URL}/logo.png`}
        />
        <Typography className={classes.title} variant="h4">
          Acronom
        </Typography>
        <Paper
          className={classes.search}
          component="form"
          elevation={0}
          onSubmit={handleSubmit}
        >
          <InputBase
            autoFocus
            className={classes.input}
            disabled={loading}
            fullWidth
            inputProps={{ 'aria-label': 'Search acronyms' }}
            onChange={(event) => setQuery(event.target.value)}
            placeholder="Search acronyms…"
            value={query}
          />
          <IconButton
            aria-label="Submit acronym search"
            className={classes.searchButton}
            disabled={loading}
            type="submit"
          >
            <SearchIcon />
          </IconButton>
        </Paper>
        <Paper className={classes.info} elevation={0}>
          <Typography variant="body2">X number of acronyms stored!</Typography>
        </Paper>
        {loading && (
          <div className={classes.message} role="status">
            <CircularProgress size={28} />
          </div>
        )}
        {error && (
          <Typography
            className={`${classes.message} ${classes.error}`}
            role="alert"
            variant="body2"
          >
            {error}
          </Typography>
        )}
        {results && results.length === 0 && (
          <Typography className={classes.message} variant="body1">
            No acronyms matched “{query.trim()}”.
          </Typography>
        )}
        {results && results.length > 0 && (
          <div aria-label="Acronym search results" className={classes.results}>
            {results.map((result) => (
              <Paper className={classes.result} elevation={0} key={result.id}>
                <Typography className={classes.resultText} variant="body1">
                  {result.link && isSafeResultLink(result.link) ? (
                    <a
                      className={`${classes.resultAcronym} ${classes.resultLink}`}
                      href={result.link}
                      rel="noreferrer"
                      target="_blank"
                    >
                      {result.acronym}
                    </a>
                  ) : (
                    <span className={classes.resultAcronym}>
                      {result.acronym}
                    </span>
                  )}
                  <span className={classes.resultDefinition}>
                    {' - '}
                    {result.definition}
                  </span>
                </Typography>
                <Tooltip title={`Edit ${result.acronym}`}>
                  <IconButton
                    aria-label={`Edit ${result.acronym}: ${result.definition}`}
                    onClick={() => setEditing(result)}
                    size="small"
                  >
                    <EditIcon fontSize="small" />
                  </IconButton>
                </Tooltip>
              </Paper>
            ))}
          </div>
        )}
      </div>
      <Snackbar
        anchorOrigin={{ vertical: 'top', horizontal: 'center' }}
        autoHideDuration={5000}
        ContentProps={{
          className:
            successType === 'updated'
              ? classes.updatedToast
              : classes.deletedToast,
        }}
        message={successMessage}
        onClose={closeSuccess}
        open={Boolean(successMessage)}
      />
    </div>
  );
};

export default Home;
