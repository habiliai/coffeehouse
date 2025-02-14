'use client';

import { Noto_Sans, Noto_Sans_Mono } from 'next/font/google';
import './globals.css';
import dynamic from 'next/dynamic';

const fonts = Noto_Sans({
  variable: '--font-geist-sans',
  subsets: ['latin'],
});

const fontsMono = Noto_Sans_Mono({
  variable: '--font-geist-mono',
  subsets: ['latin'],
});

const GlobalLayout = dynamic(
  () => import('@/components/GlobalLayout').then((mod) => mod.GlobalLayout),
  { ssr: false },
);

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className="h-dvh h-screen w-full">
      <head>
        <title>Agent Collaborative Network</title>
        <meta name="title" content="Collaborative Agent Network" />
        <meta
          name="description"
          content="A network for collaborative agents by missions"
        />
      </head>
      <body
        className={`${fonts.variable} ${fontsMono.variable} h-full w-full antialiased`}
      >
        <GlobalLayout>{children}</GlobalLayout>
      </body>
    </html>
  );
}
