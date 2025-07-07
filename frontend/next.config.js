/** @type {import('next').NextConfig} */
const nextConfig = {
  // Enable standalone mode for Docker optimization
  output: 'standalone',
  
  // Optimize images
  images: {
    unoptimized: false,
    remotePatterns: [
      {
        protocol: 'https',
        hostname: '**',
      },
    ],
  },

  // Enable SWC minification
  swcMinify: true,

  // Compress responses
  compress: true,

  // Optimize bundle
  experimental: {
    optimizeCss: true,
    optimizePackageImports: ['lucide-react', '@radix-ui/react-icons'],
  },

  // Environment variables
  env: {
    NEXT_PUBLIC_API_URL: process.env.NEXT_PUBLIC_API_URL,
  },
}

module.exports = nextConfig
