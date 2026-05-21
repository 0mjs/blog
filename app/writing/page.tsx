import Link from "next/link";
import type { Metadata } from "next";
import { SiteTitle } from "@/components/site-title";
import { getPublishedPosts } from "@/lib/content";
import { formatLongDate } from "@/lib/format";

export const metadata: Metadata = {
  title: "Writing",
  description: "Posts by Matt J. Stevenson.",
};

export default async function WritingPage() {
  const posts = await getPublishedPosts();

  return (
    <section>
      <SiteTitle />

      <div className="-mt-[14px] flex flex-col divide-y divide-[var(--border)]">
        {posts.map((post) => (
          <Link
            key={post.slug}
            href={`/writing/${post.slug}`}
            className="text-[var(--text)] no-underline transition hover:opacity-70"
          >
            <div className="flex flex-col items-start pb-4 pt-3 md:flex-row md:pb-2">
              <div className="min-w-0 flex-auto md:w-9/12">
                <h2 className="nimbus truncate text-[var(--text)]">{post.title}</h2>
              </div>
              <time
                className="relative top-px shrink-0 text-xs text-[var(--muted)] md:w-3/12 md:pl-4"
                dateTime={post.date?.toISOString()}
              >
                {post.date ? formatLongDate(post.date) : null}
              </time>
            </div>
          </Link>
        ))}
      </div>
    </section>
  );
}
