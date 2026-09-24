import React from 'react';
import {
  BrowserRouter as Router,
  Redirect,
  Route,
  Switch,
} from 'react-router-dom';
import Home from './views/Home';
import Layout from './views/Layout';

interface RoutesProps {
  mode: 'light' | 'dark';
  onToggleMode: () => void;
}

const Routes: React.FC<RoutesProps> = ({ mode, onToggleMode }) => (
  <Router basename="/acronom">
    <Layout mode={mode} onToggleMode={onToggleMode}>
      <Switch>
        <Route exact path="/" component={Home} />
        <Redirect to="/" />
      </Switch>
    </Layout>
  </Router>
);

export default Routes;
