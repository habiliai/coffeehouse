import { createContext, ReactNode, useContext, useState } from 'react';
import { HabiliApiClient } from '@/proto/habapi';

const habiliApiClientContext = createContext<HabiliApiClient | null>(null);

export function HabiliApiClientProvider({ children }: { children: ReactNode }) {
  const [client] = useState(
    () =>
      new HabiliApiClient(
        process.env.NEXT_PUBLIC_HABILIAPI_URL ?? 'http://localhost:8001',
      ),
  );
  return (
    <habiliApiClientContext.Provider value={client}>
      {children}
    </habiliApiClientContext.Provider>
  );
}

export function useHabiliApiClient() {
  const client = useContext(habiliApiClientContext);
  if (!client) {
    throw new Error(
      'useHabiliApiClient must be used within a HabiliApiClientProvider',
    );
  }

  return client;
}
