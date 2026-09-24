import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Header from '../../components/Header';

const useStyles = makeStyles((theme) => ({
  root: {
    minHeight: '100vh',
    color: theme.palette.text.primary,
    backgroundColor: theme.palette.background.default,
    backgroundImage:
      theme.palette.type === 'light'
        ? 'radial-gradient(circle at 85% 0%, rgba(64, 84, 178, 0.08), transparent 28%)'
        : 'radial-gradient(circle at 85% 0%, rgba(142, 162, 255, 0.08), transparent 28%)',
  },
}));

interface LayoutProps {
  mode: 'light' | 'dark';
  onToggleMode: () => void;
}

const Layout: React.FC<LayoutProps> = ({ children, mode, onToggleMode }) => {
  const classes = useStyles();

  return (
    <div className={classes.root}>
      <Header mode={mode} onToggleMode={onToggleMode} />
      <main>{children}</main>
    </div>
  );
};

export default Layout;
