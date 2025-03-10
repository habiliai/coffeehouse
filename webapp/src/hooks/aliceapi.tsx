import { createContext, ReactNode, useContext, useState } from 'react';
import { AliceApiClient } from '@/proto/aliceapi';

const aliceApiClientContext = createContext<AliceApiClient | null>(null);

export function AliceApiClientProvider({ children }: { children: ReactNode }) {
  const [client] = useState(
    () =>
      new AliceApiClient(
        process.env.NEXT_PUBLIC_ALICEAPI_URL ?? 'http://localhost:8001',
      ),
  );
  return (
    <aliceApiClientContext.Provider value={client}>
      {children}
    </aliceApiClientContext.Provider>
  );
}

export function useAliceApiClient() {
  const client = useContext(aliceApiClientContext);
  if (!client) {
    throw new Error(
      'useAliceApiClient must be used within a AliceApiClientProvider',
    );
  }

  return client;
}
