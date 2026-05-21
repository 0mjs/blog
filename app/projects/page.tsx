import { SiteTitle } from "@/components/site-title";
import { versionedAsset } from "@/lib/assets";
import { ArrowUpRight, Wrench } from "lucide-react";
import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Projects",
  description: "Products and projects from Matt J. Stevenson.",
};

type ProjectLink = {
  label: string;
  href: string;
  detail: string;
  logo: string;
  logoStyle?: "cutwise";
  logoScale?: "small";
};

const projectLinks: ProjectLink[] = [
  {
    label: "CutWise",
    href: "https://cutwise.fit",
    detail: "An AI agent for weight loss, dieting, and training",
    logo: "/image/projects/cutwise.png",
    logoStyle: "cutwise",
  },
  {
    label: "Zinc",
    href: "https://zinc.carbonsoft.sh",
    detail: "A Go API framework inspired by Express.js",
    logo: "/image/projects/zinc.png",
    logoScale: "small",
  },
  {
    label: "QuitShark",
    href: "https://apps.apple.com/us/app/quitshark-quit-nicotine/id6759577475",
    detail: "An iOS app to help young people quit nicotine",
    logo: "/image/projects/quitshark.png",
  },
  {
    label: "Nuddge",
    href: "https://nuddge.app",
    detail: "A calm wellness iOS app for combatting overwhelm",
    logo: "/image/projects/nuddge.png",
  },
  {
    label: "Samsa",
    href: "https://hih.ie/initiatives/hihi-ai/hihi-ai-call-winners-2025/samsa-ltd/",
    detail: "An AI MedTech system and HiHi-AI Call 2025 winner",
    logo: "/image/projects/samsa.png",
  },
  {
    label: "IconAI",
    href: "https://www.producthunt.com/products/iconai",
    detail: "A solo AI image generation SaaS using DALL-E 3",
    logo: "/image/projects/iconai.png",
    logoScale: "small",
  },
];

export default function ProjectsPage() {
  return (
    <section>
      <SiteTitle />

      <div className="grid gap-8">
        <section className="border-t border-[var(--border)] pt-6">
          <div className="mb-4 flex items-center gap-2 text-[var(--text)]">
            <Wrench size={18} className="text-[var(--accent)]" />
            <h2 className="text-xl font-medium">Projects</h2>
          </div>
          <div className="grid gap-3 sm:grid-cols-2">
            {projectLinks.map((item) => (
              <a
                key={item.label}
                href={item.href}
                target="_blank"
                rel="noopener noreferrer"
                className="group rounded-lg border border-[var(--border)] p-4 no-underline transition hover:border-[var(--accent)] hover:bg-[var(--surface)]"
              >
                <div className="flex items-start justify-between gap-3">
                  <div className="flex min-w-0 items-center gap-3">
                    <span
                      className={
                        item.logoStyle === "cutwise"
                          ? "grid size-9 shrink-0 place-items-center rounded-lg border border-white/30 bg-[#3d6a4a] shadow-[inset_0_1px_0_rgba(255,255,255,0.18),0_1px_2px_rgba(25,23,15,0.08)]"
                          : "grid size-9 shrink-0 place-items-center overflow-hidden rounded-lg border border-[var(--border)] bg-[var(--bg)]"
                      }
                      aria-hidden="true"
                    >
                      <img
                        src={versionedAsset(item.logo)}
                        alt=""
                        width={item.logoStyle === "cutwise" ? 24 : 36}
                        height={item.logoStyle === "cutwise" ? 24 : 36}
                        className={
                          item.logoStyle === "cutwise"
                            ? "size-[68%] object-contain"
                            : item.logoScale === "small"
                              ? "size-[78%] object-contain"
                              : "size-full object-contain"
                        }
                      />
                    </span>
                    <p className="truncate font-medium text-[var(--text)]">{item.label}</p>
                  </div>
                  <ArrowUpRight
                    size={16}
                    className="shrink-0 text-[var(--muted)] transition group-hover:-translate-y-0.5 group-hover:translate-x-0.5 group-hover:text-[var(--accent)]"
                  />
                </div>
                <p className="mt-2 text-xs leading-5 text-[var(--muted)]">{item.detail}</p>
              </a>
            ))}
          </div>
        </section>
      </div>
    </section>
  );
}
