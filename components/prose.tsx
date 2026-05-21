import type { ContentEntry } from "@/lib/content";

export function Prose({ entry, className = "" }: { entry: ContentEntry; className?: string }) {
  return (
    <article className={`prose-content ${className}`}>
      <div dangerouslySetInnerHTML={{ __html: entry.html }} />
    </article>
  );
}
