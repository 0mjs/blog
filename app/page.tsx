import { Prose } from "@/components/prose";
import { SiteTitle } from "@/components/site-title";
import { getEntryBySlug } from "@/lib/content";

export default async function HomePage() {
  const about = await getEntryBySlug("about");

  if (!about) {
    throw new Error("Missing site/content/about.md");
  }

  return (
    <section className="max-w-[65ch]">
      <SiteTitle />
      <Prose entry={about} />
    </section>
  );
}
