import type { Metadata } from "next";
import "@fontsource/geist/400.css";
import "@fontsource/geist/500.css";
import "@fontsource/geist/600.css";
import "@fontsource/geist/700.css";
import "@fontsource/ibm-plex-mono/400.css";
import "@fontsource/ibm-plex-mono/500.css";
import "@fontsource/ibm-plex-mono/600.css";
import "@fontsource/ibm-plex-mono/400-italic.css";
import "./globals.css";
import { SiteShell } from "@/components/site-shell";

export const metadata: Metadata = {
  metadataBase: new URL("https://blog.0mjs.dev"),
  title: {
    default: "Matt J. Stevenson",
    template: "%s - Matt J. Stevenson",
  },
  description: "Writing about software engineering, systems, and building useful products.",
};

export default function RootLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body>
        <script
          dangerouslySetInnerHTML={{
            __html: `(() => { try { const theme = localStorage.getItem("theme"); if (theme === "light" || theme === "dark") document.documentElement.dataset.theme = theme; } catch {} })();`,
          }}
        />
        <SiteShell>{children}</SiteShell>
      </body>
    </html>
  );
}
