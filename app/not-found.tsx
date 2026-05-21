import Link from "next/link";

export default function NotFound() {
  return (
    <section className="max-w-xl">
      <p className="mb-4 text-sm font-medium uppercase tracking-[0.18em] text-[var(--accent)]">404</p>
      <h1 className="text-4xl font-semibold text-[var(--text)]">Page not found.</h1>
      <p className="mt-4 text-[var(--muted)]">That page is not part of the current archive.</p>
      <Link href="/" className="mt-8 inline-flex no-underline hover:text-[var(--accent)]">
        Back to About
      </Link>
    </section>
  );
}
