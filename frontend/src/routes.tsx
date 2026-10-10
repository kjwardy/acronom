import React from 'react';
import {
  BrowserRouter as Router,
  Redirect,
  Route,
  Switch,
} from 'react-router-dom';
import { AcronymEventsProvider } from './contexts/AcronymEvents';
import { CategoriesProvider } from './contexts/Categories';
import { MetricsProvider } from './contexts/Metrics';
import Home from './views/Home';
import Layout from './views/Layout';

interface RoutesProps {
  mode: 'light' | 'dark';
  onToggleMode: () => void;
}

const Routes: React.FC<RoutesProps> = ({ mode, onToggleMode }) => (
  <Router basename="/acronom">
    <MetricsProvider>
      <CategoriesProvider>
        <AcronymEventsProvider>
          <Layout mode={mode} onToggleMode={onToggleMode}>
            <Switch>
              <Route exact path="/" component={Home} />
              <Redirect to="/" />
            </Switch>
          </Layout>
        </AcronymEventsProvider>
      </CategoriesProvider>
    </MetricsProvider>
  </Router>
);

export default Routes;
