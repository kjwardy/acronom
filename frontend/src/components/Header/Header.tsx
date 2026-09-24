import React from 'react';
import AppBar from '@material-ui/core/AppBar';
import IconButton from '@material-ui/core/IconButton';
import Toolbar from '@material-ui/core/Toolbar';
import Tooltip from '@material-ui/core/Tooltip';
import Typography from '@material-ui/core/Typography';
import Brightness4Icon from '@material-ui/icons/Brightness4';
import Brightness7Icon from '@material-ui/icons/Brightness7';
import { makeStyles } from '@material-ui/core/styles';

const useStyles = makeStyles((theme) => ({
  link: {
    display: 'flex',
    alignItems: 'center',
    padding: theme.spacing(1),
    color: theme.palette.common.white,
    textDecoration: 'none',
    borderRadius: theme.shape.borderRadius,
    transition: theme.transitions.create('background-color'),
    '&:hover': {
      backgroundColor: 'rgba(255, 255, 255, 0.08)',
    },
  },
  brand: {
    fontWeight: 600,
  },
  grow: {
    flexGrow: 1,
  },
}));

interface HeaderProps {
  mode: 'light' | 'dark';
  onToggleMode: () => void;
}

const Header: React.FC<HeaderProps> = ({ mode, onToggleMode }) => {
  const classes = useStyles();

  return (
    <AppBar position="static">
      <Toolbar>
        <a className={classes.link} href="/acronom">
          <Typography className={classes.brand} color="inherit" variant="h6">
            Acronom
          </Typography>
        </a>
        <div className={classes.grow} />
        <Tooltip
          title={`Switch to ${mode === 'light' ? 'dark' : 'light'} mode`}
        >
          <IconButton
            aria-label="Toggle dark mode"
            color="inherit"
            onClick={onToggleMode}
          >
            {mode === 'light' ? <Brightness4Icon /> : <Brightness7Icon />}
          </IconButton>
        </Tooltip>
      </Toolbar>
    </AppBar>
  );
};

export default Header;
