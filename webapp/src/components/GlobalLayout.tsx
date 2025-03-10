'use client';

import { AliceApiClientProvider } from '@/hooks/aliceapi';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useState } from 'react';

export function GlobalLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const [queryClient] = useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: {
            refetchOnWindowFocus: false,
          },
        },
      }),
  );

  return (
    <div className="flex h-full w-full flex-1">
      <QueryClientProvider client={queryClient}>
        <AliceApiClientProvider>{children}</AliceApiClientProvider>
      </QueryClientProvider>
    </div>
  );
}
