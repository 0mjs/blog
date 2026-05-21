import Link from "next/link";

type SiteTitleProps = {
  className?: string;
};

export function SiteTitle({ className = "mb-20" }: SiteTitleProps) {
  return (
    <h1 className={`site-heading ${className}`}>
      <Link href="/">Matt J. Stevenson</Link>
    </h1>
  );
}
