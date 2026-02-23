

const nextConfig = {
  output: 'standalone',
  images: {
    remotePatterns: [
      { protocol: 'https', hostname: 'image.tmdb.org' },
      { protocol: 'https', hostname: '*.coverartarchive.org' },
      { protocol: 'https', hostname: 'images.igdb.com' },
    ],
  },
}

export default nextConfig
