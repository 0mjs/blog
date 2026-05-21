import { SiteTitle } from "@/components/site-title";
import { versionedAsset } from "@/lib/assets";
import { gamingLinks, musicArtists } from "@/lib/media";
import { ArrowUpRight, Gamepad2, Music2 } from "lucide-react";
import type { Metadata } from "next";
import Image from "next/image";

export const metadata: Metadata = {
  title: "Media",
  description: "Gaming videos and music from Matt J. Stevenson.",
};

export default function MediaPage() {
  return (
    <section>
      <SiteTitle />

      <div className="grid gap-8">
        <section className="border-t border-[var(--border)] pt-6">
          <div className="mb-4 flex items-center gap-2 text-[var(--text)]">
            <Gamepad2 size={18} className="text-[var(--accent)]" />
            <h2 className="text-xl font-medium">Gaming</h2>
          </div>
          <div className="grid gap-3 sm:grid-cols-2">
            {gamingLinks.map((item) => (
              <a
                key={item.href}
                href={item.href}
                target="_blank"
                rel="noopener noreferrer"
                className="group rounded-lg border border-[var(--border)] p-4 no-underline transition hover:border-[var(--accent)] hover:bg-[var(--surface)]"
              >
                <div className="flex items-start justify-between gap-3">
                  <div className="flex items-center gap-2">
                    <Image
                      src={versionedAsset(item.image)}
                      alt=""
                      width={28}
                      height={28}
                      className="rounded-md"
                      unoptimized
                    />
                    <p className="font-medium text-[var(--text)]">{item.label}</p>
                  </div>
                  <ArrowUpRight
                    size={16}
                    className="text-[var(--muted)] transition group-hover:-translate-y-0.5 group-hover:translate-x-0.5 group-hover:text-[var(--accent)]"
                  />
                </div>
                <p className="mt-2 text-xs leading-5 text-[var(--muted)]">{item.detail}</p>
              </a>
            ))}
          </div>
        </section>

        <section className="border-t border-[var(--border)] pt-6">
          <div className="mb-4 flex items-center gap-2 text-[var(--text)]">
            <Music2 size={18} className="text-[var(--accent)]" />
            <h2 className="text-xl font-medium">Music</h2>
          </div>
          <div className="grid gap-3 sm:grid-cols-2">
            {musicArtists.map((artist) => (
              <a
                key={artist.href}
                href={artist.href}
                target="_blank"
                rel="noopener noreferrer"
                className="group relative aspect-[5/3] overflow-hidden rounded-lg border border-[var(--border)] bg-[var(--surface)] no-underline"
              >
                <Image
                  src={versionedAsset(artist.image)}
                  alt={artist.label}
                  fill
                  sizes="(min-width: 640px) 50vw, 100vw"
                  className="object-cover transition duration-300 group-hover:scale-[1.03]"
                  unoptimized
                />
                <div className="absolute inset-0 bg-gradient-to-t from-black/70 via-black/10 to-transparent" />
                <div className="absolute inset-x-0 bottom-0 flex items-center justify-between gap-3 p-4">
                  <p className="font-medium text-white">{artist.label}</p>
                  <ArrowUpRight size={16} className="text-white/80" />
                </div>
              </a>
            ))}
          </div>
        </section>
      </div>
    </section>
  );
}
