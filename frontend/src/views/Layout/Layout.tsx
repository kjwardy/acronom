import React, { useState } from 'react';
import Fab from '@material-ui/core/Fab';
import Snackbar from '@material-ui/core/Snackbar';
import Tooltip from '@material-ui/core/Tooltip';
import AddIcon from '@material-ui/icons/Add';
import { makeStyles } from '@material-ui/core/styles';
import AddAcronymModal from '../../components/AddAcronymModal';
import Header from '../../components/Header';
import { useMetrics } from '../../contexts/Metrics';

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
  addButton: {
    position: 'fixed',
    right: theme.spacing(3),
    bottom: theme.spacing(3),
    boxShadow: '0 8px 24px rgba(20, 29, 60, 0.24)',
  },
  successToast: {
    width: 'auto',
    minWidth: 'unset',
    color: theme.palette.common.white,
    backgroundColor: '#2e7d32',
    '& .MuiSnackbarContent-message': {
      textAlign: 'center',
    },
  },
}));

interface LayoutProps {
  mode: 'light' | 'dark';
  onToggleMode: () => void;
}

const Layout: React.FC<LayoutProps> = ({ children, mode, onToggleMode }) => {
  const classes = useStyles();
  const { incrementTotal } = useMetrics();
  const [addOpen, setAddOpen] = useState(false);
  const [successMessage, setSuccessMessage] = useState('');

  const handleCreated = () => {
    setAddOpen(false);
    incrementTotal();
    setSuccessMessage('Acronym added successfully!');
  };

  const closeSuccess = (_event?: React.SyntheticEvent, reason?: string) => {
    if (reason !== 'clickaway') setSuccessMessage('');
  };

  return (
    <div className={classes.root}>
      <Header mode={mode} onToggleMode={onToggleMode} />
      {addOpen && (
        <AddAcronymModal
          onClose={() => setAddOpen(false)}
          onCreated={handleCreated}
        />
      )}
      <main>{children}</main>
      <Tooltip title="Add new acronym">
        <Fab
          aria-label="Add new acronym"
          className={classes.addButton}
          color="secondary"
          onClick={() => setAddOpen(true)}
        >
          <AddIcon />
        </Fab>
      </Tooltip>
      <Snackbar
        anchorOrigin={{ vertical: 'top', horizontal: 'center' }}
        autoHideDuration={5000}
        ContentProps={{ className: classes.successToast }}
        message={successMessage}
        onClose={closeSuccess}
        open={Boolean(successMessage)}
      />
    </div>
  );
};

export default Layout;
