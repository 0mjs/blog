import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  async redirects() {
    return [
      {
        source: "/about",
        destination: "/",
        permanent: true,
      },
      {
        source: "/post/:slug",
        destination: "/writing/:slug",
        permanent: true,
      },
      {
        source: "/misc",
        destination: "/projects",
        permanent: true,
      },
    ];
  },
};

export default nextConfig;
