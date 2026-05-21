import { notFound } from "next/navigation";
import type { Metadata } from "next";
import { Prose } from "@/components/prose";
import { SiteTitle } from "@/components/site-title";
import { getEntryBySlug, getPublishedPosts } from "@/lib/content";
import { formatLongDate } from "@/lib/format";

type PageProps = {
  params: Promise<{ slug: string }>;
};

export async function generateStaticParams() {
  const posts = await getPublishedPosts();
  return posts.map((post) => ({ slug: post.slug }));
}

export async function generateMetadata({ params }: PageProps): Promise<Metadata> {
  const { slug } = await params;
  const post = await getEntryBySlug(slug);

  if (!post || post.draft || post.slug === "about") {
    return {};
  }

  return {
    title: post.title,
    description: `${post.title} by Matt J. Stevenson.`,
  };
}

export default async function WritingPostPage({ params }: PageProps) {
  const { slug } = await params;
  const post = await getEntryBySlug(slug);

  if (!post || post.draft || post.slug === "about") {
    notFound();
  }

  return (
    <article>
      <SiteTitle className="mb-14" />

      <header className="mb-10 border-b border-[var(--border)] pb-7">
        <h1 className="max-w-3xl text-3xl font-semibold leading-tight text-[var(--text)] sm:text-4xl">
          {post.title}
        </h1>
        <div className="mt-4 flex flex-wrap items-center gap-2 text-sm text-[var(--muted)]">
          {post.date ? <time dateTime={post.date.toISOString()}>{formatLongDate(post.date)}</time> : null}
          <span className="rounded-full border border-[color-mix(in_srgb,var(--accent)_42%,transparent)] px-2 py-0.5 text-xs text-[var(--accent)]">
            {post.readTime} min read
          </span>
        </div>
        {post.tags.length ? (
          <div className="mt-4 flex flex-wrap gap-1.5">
            {post.tags.map((tag) => (
              <span
                key={tag}
                className="rounded-full border border-[color-mix(in_srgb,var(--accent)_28%,transparent)] px-2 py-0.5 text-[11px] text-[var(--accent)] opacity-80"
              >
                {tag}
              </span>
            ))}
          </div>
        ) : null}
      </header>

      <Prose entry={post} />
    </article>
  );
}
