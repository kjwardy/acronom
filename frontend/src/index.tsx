import React, { useMemo, useState } from 'react';
import ReactDOM from 'react-dom';
import CssBaseline from '@material-ui/core/CssBaseline';
import { createTheme, ThemeProvider } from '@material-ui/core/styles';
import 'typeface-roboto';
import './index.css';
import Routes from './routes';

const getInitialMode = () => {
  const savedMode = localStorage.getItem('theme');
  if (savedMode === 'light' || savedMode === 'dark') return savedMode;
  return window.matchMedia('(prefers-color-scheme: dark)').matches
    ? 'dark'
    : 'light';
};

const App = () => {
  const [mode, setMode] = useState<'light' | 'dark'>(getInitialMode);
  const theme = useMemo(
    () =>
      createTheme({
        palette: {
          type: mode,
          primary: { main: mode === 'light' ? '#4054b2' : '#8ea2ff' },
          secondary: { main: '#00a58e' },
          background: {
            default: mode === 'light' ? '#f4f6fb' : '#11141b',
            paper: mode === 'light' ? '#ffffff' : '#1b202a',
          },
        },
        shape: { borderRadius: 10 },
        typography: {
          fontFamily: 'Roboto, sans-serif',
          h4: { fontWeight: 600 },
          h6: { fontWeight: 600 },
          button: { fontWeight: 600, textTransform: 'none' },
        },
        overrides: {
          MuiAppBar: {
            colorPrimary: {
              color: '#ffffff',
              backgroundColor: mode === 'light' ? '#33469b' : '#1b202a',
            },
          },
          MuiPaper: {
            rounded: { borderRadius: 12 },
          },
          MuiButton: {
            root: { borderRadius: 8, paddingLeft: 16, paddingRight: 16 },
          },
        },
      }),
    [mode],
  );

  const toggleMode = () => {
    const nextMode = mode === 'light' ? 'dark' : 'light';
    localStorage.setItem('theme', nextMode);
    setMode(nextMode);
  };

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Routes mode={mode} onToggleMode={toggleMode} />
    </ThemeProvider>
  );
};

ReactDOM.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
  document.getElementById('root'),
);
