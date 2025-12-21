/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  // Separate client and admin apps using subpaths
  async rewrites() {
    return [
      {
        source: '/client/:path*',
        destination: '/client/:path*',
      },
      {
        source: '/admin/:path*',
        destination: '/admin/:path*',
      },
    ];
  },
  env: {
    NEXT_PUBLIC_API_URL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1',
  },
};

module.exports = nextConfig;

