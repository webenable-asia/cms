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

  // Compress responses
  compress: true,

  // Optimize bundle
  experimental: {
    // optimizeCss: true, // Disabled temporarily due to critters dependency issue
    optimizePackageImports: ['lucide-react', '@radix-ui/react-icons'],
  },

  // Environment variables
  env: {
    NEXT_PUBLIC_API_URL: process.env.NEXT_PUBLIC_API_URL,
  },
}

module.exports = nextConfig
