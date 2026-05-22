import Image from "next/image";
import Link from "next/link";
import { Mail } from "lucide-react";
import { FaGithub, FaLinkedin, FaXTwitter, FaYoutube } from "react-icons/fa6";
import { versionedAsset } from "@/lib/assets";
import { SiteNav } from "@/components/site-nav";
// import { ThemeToggle } from "@/components/theme-toggle";

const socialItems = [
  { href: "mailto:dev@mattjs.me", label: "Email", icon: <Mail size={16} strokeWidth={1.8} /> },
  { href: "https://github.com/0mjs", label: "GitHub", icon: <FaGithub size={16} /> },
  { href: "https://x.com/0mjs_", label: "X", icon: <FaXTwitter size={15} /> },
  { href: "https://linkedin.com/in/matt-j-stevenson", label: "LinkedIn", icon: <FaLinkedin size={16} /> },
  { href: "https://www.youtube.com/@slaikers/playlists", label: "YouTube", icon: <FaYoutube size={17} /> },
];

const logoSrc = versionedAsset("/image/logo.png");

export function SiteShell({ children }: { children: React.ReactNode }) {
  return (
    <div className="mx-auto grid min-h-dvh w-full max-w-6xl grid-cols-1 px-5 py-10 sm:px-8 lg:grid-cols-[220px_minmax(0,1fr)] lg:gap-16 lg:py-24">
      <aside className="lg:sticky lg:top-24 lg:flex lg:h-[calc(100dvh-12rem)] lg:flex-col">
        <div className="flex items-center justify-between gap-4 lg:block">
          <Link href="/" className="inline-flex no-underline" aria-label="Matt J. Stevenson">
            <Image
              src={logoSrc}
              alt="Matt J. Stevenson"
              width={34}
              height={34}
              className="rounded-md"
              unoptimized
              priority
            />
          </Link>
          {/* <div className="lg:hidden">
            <ThemeToggle />
          </div> */}
        </div>

        <SiteNav />

        {/* <div className="mt-auto hidden lg:block">
          <div className="mb-5">
            <ThemeToggle />
          </div>
        </div> */}
      </aside>

      <div className="flex min-w-0 flex-col">
        <main className="min-w-0 flex-1 pt-16 lg:pt-0">{children}</main>
        <footer className="mt-24 border-t border-[var(--border)] py-8 text-sm text-[var(--muted)]">
          <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
            <p>© {new Date().getFullYear()} Matt J. Stevenson.</p>
            <div className="flex items-center gap-2 sm:justify-end">
              {socialItems.map((item) => (
                <a
                  key={item.href}
                  href={item.href}
                  aria-label={item.label}
                  title={item.label}
                  className="inline-flex size-8 items-center justify-center rounded-md text-[var(--muted)] transition hover:bg-[var(--surface)] hover:text-[var(--text)]"
                  target={item.href.startsWith("http") ? "_blank" : undefined}
                  rel={item.href.startsWith("http") ? "noopener noreferrer" : undefined}
                >
                  {item.icon}
                </a>
              ))}
            </div>
          </div>
        </footer>
      </div>
    </div>
  );
}
