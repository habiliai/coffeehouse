'use client';

import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useState } from 'react';

export function GlobalLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const [queryClient] = useState(() => new QueryClient());

  return (
    <div className="flex h-full w-full flex-1">
      <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    </div>
  );
}
