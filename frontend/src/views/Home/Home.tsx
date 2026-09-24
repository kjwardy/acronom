import React, { FormEvent, useState } from 'react';
import InputBase from '@material-ui/core/InputBase';
import Paper from '@material-ui/core/Paper';
import Typography from '@material-ui/core/Typography';
import SearchIcon from '@material-ui/icons/Search';
import { makeStyles } from '@material-ui/core/styles';

const useStyles = makeStyles((theme) => ({
  root: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    minHeight: 'calc(100vh - 64px)',
    padding: theme.spacing(3),
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
    padding: theme.spacing(0.75, 2),
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
  icon: {
    display: 'flex',
    marginRight: theme.spacing(1.5),
    color: theme.palette.text.secondary,
  },
  input: {
    flex: 1,
    fontSize: '1.1rem',
  },
}));

const Home: React.FC = () => {
  const [query, setQuery] = useState('');
  const classes = useStyles();

  const handleSubmit = (event: FormEvent) => {
    event.preventDefault();
  };

  return (
    <div className={classes.root}>
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
          <span className={classes.icon}>
            <SearchIcon />
          </span>
          <InputBase
            autoFocus
            className={classes.input}
            fullWidth
            inputProps={{ 'aria-label': 'Search acronyms' }}
            onChange={(event) => setQuery(event.target.value)}
            placeholder="Search acronyms…"
            value={query}
          />
        </Paper>
        <Paper className={classes.info} elevation={0}>
          <Typography variant="body2">X number of acronyms stored!</Typography>
        </Paper>
      </div>
    </div>
  );
};

export default Home;
