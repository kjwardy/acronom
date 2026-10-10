import React, { createContext, useContext, useState } from 'react';
import { Category } from '../Categories';

export interface AcronymEntry {
  id: number;
  acronym: string;
  definition: string;
  link?: string;
  created_at?: string;
  updated_at?: string;
  categories?: Category[];
}

interface CreatedAcronymEvent {
  entry: AcronymEntry;
  sequence: number;
}

interface AcronymEventsContextValue {
  created: CreatedAcronymEvent | null;
  notifyCreated: (entry: AcronymEntry) => void;
}

const AcronymEventsContext = createContext<
  AcronymEventsContextValue | undefined
>(undefined);

export const AcronymEventsProvider: React.FC = ({ children }) => {
  const [created, setCreated] = useState<CreatedAcronymEvent | null>(null);

  const notifyCreated = (entry: AcronymEntry) => {
    setCreated((current) => ({
      entry,
      sequence: (current?.sequence || 0) + 1,
    }));
  };

  return (
    <AcronymEventsContext.Provider value={{ created, notifyCreated }}>
      {children}
    </AcronymEventsContext.Provider>
  );
};

export const useAcronymEvents = (): AcronymEventsContextValue => {
  const context = useContext(AcronymEventsContext);
  if (!context)
    throw new Error(
      'useAcronymEvents must be used within AcronymEventsProvider',
    );
  return context;
};
