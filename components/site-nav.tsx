"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

const navItems = [
  { href: "/", label: "About" },
  { href: "/writing", label: "Writing" },
  { href: "/projects", label: "Projects" },
  { href: "/media", label: "Media" },
];

export function SiteNav() {
  const pathname = usePathname();

  return (
    <nav className="mt-8 flex items-center gap-2 overflow-x-auto text-sm lg:flex-col lg:items-start lg:gap-1">
      {navItems.map((item) => {
        const isActive = item.href === "/" ? pathname === "/" : pathname.startsWith(item.href);

        return (
          <Link
            key={item.href}
            href={item.href}
            aria-current={isActive ? "page" : undefined}
            className="group flex items-center gap-2 rounded-md px-2 py-1.5 text-[var(--muted)] no-underline transition hover:opacity-60"
          >
            <span
              aria-hidden="true"
              className={`size-1 rounded-full bg-[var(--accent)] transition-transform ${isActive ? "scale-125" : "scale-0"}`}
            />
            <span className={isActive ? "text-[var(--accent)]" : ""}>{item.label}</span>
          </Link>
        );
      })}
    </nav>
  );
}
