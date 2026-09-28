import React, { createContext, useContext, useEffect, useState } from 'react';
import axios from 'axios';

interface MetricsResponse {
  total_acronyms: number;
}

interface MetricsContextValue {
  totalAcronyms: number | null;
  loading: boolean;
  failed: boolean;
  incrementTotal: () => void;
  decrementTotal: () => void;
}

const MetricsContext = createContext<MetricsContextValue | undefined>(
  undefined,
);

export const fetchTotalAcronyms = async (): Promise<number> => {
  const { data } = await axios.get<MetricsResponse>(
    '/api/metrics/total-acronyms',
    { timeout: 3000 },
  );
  if (!Number.isInteger(data?.total_acronyms) || data.total_acronyms < 0) {
    throw new Error('Invalid metrics response');
  }
  return data.total_acronyms;
};

export const adjustMetricTotal = (total: number, delta: number): number =>
  Math.max(0, total + delta);

export const formatAcronymCount = (total: number): string =>
  `${total} ${total === 1 ? 'meaning' : 'meanings'} decoded`;

export const MetricsProvider: React.FC = ({ children }) => {
  const [totalAcronyms, setTotalAcronyms] = useState<number | null>(null);
  const [loading, setLoading] = useState(true);
  const [failed, setFailed] = useState(false);

  useEffect(() => {
    let mounted = true;
    fetchTotalAcronyms()
      .then((total) => {
        if (mounted) setTotalAcronyms(total);
      })
      .catch(() => {
        if (mounted) setFailed(true);
      })
      .finally(() => {
        if (mounted) setLoading(false);
      });
    return () => {
      mounted = false;
    };
  }, []);

  const adjustTotal = (delta: number) => {
    setTotalAcronyms((total) =>
      total === null ? null : adjustMetricTotal(total, delta),
    );
  };

  return (
    <MetricsContext.Provider
      value={{
        totalAcronyms,
        loading,
        failed,
        incrementTotal: () => adjustTotal(1),
        decrementTotal: () => adjustTotal(-1),
      }}
    >
      {children}
    </MetricsContext.Provider>
  );
};

export const useMetrics = (): MetricsContextValue => {
  const metrics = useContext(MetricsContext);
  if (!metrics)
    throw new Error('useMetrics must be used within MetricsProvider');
  return metrics;
};
