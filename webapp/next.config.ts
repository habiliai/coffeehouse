import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  images: {
    unoptimized: true,
    remotePatterns: [
      {
        protocol: 'https',
        hostname: 'img.logo.dev',
        pathname: '/**',
      },
      {
        protocol: 'https',
        hostname: 'u6mo491ntx4iwuoz.public.blob.vercel-storage.com',
        pathname: '/**',
      }
    ]
  },
  reactStrictMode: true,
};

export default nextConfig;
