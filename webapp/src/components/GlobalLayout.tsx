'use client';

import { HabiliApiClientProvider } from '@/hooks/habapi';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useState } from 'react';

export function GlobalLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const [queryClient] = useState(() => new QueryClient({
    defaultOptions: {
      queries: {
        refetchOnWindowFocus: false,
      }
    }
  }));

  return (
    <div className="flex h-full w-full flex-1">
      <QueryClientProvider client={queryClient}>
        <HabiliApiClientProvider>{children}</HabiliApiClientProvider>
      </QueryClientProvider>
    </div>
  );
}
