"use client";

import { Monitor, Moon, Sun } from "lucide-react";
import { useState } from "react";

type Theme = "system" | "light" | "dark";

const themes: Theme[] = ["system", "light", "dark"];

export function ThemeToggle() {
  const [theme, setTheme] = useState<Theme>(() => {
    if (typeof window === "undefined") {
      return "system";
    }

    const saved = window.localStorage.getItem("theme");
    if (saved === "light" || saved === "dark" || saved === "system") {
      return saved;
    }

    return "system";
  });

  function cycleTheme() {
    const next = themes[(themes.indexOf(theme) + 1) % themes.length];
    setTheme(next);
    applyTheme(next);
    window.localStorage.setItem("theme", next);
  }

  const Icon = theme === "dark" ? Moon : theme === "light" ? Sun : Monitor;

  return (
    <button
      type="button"
      className="inline-flex size-9 items-center justify-center rounded-md border border-[var(--border)] text-[var(--muted)] transition hover:border-[var(--text)] hover:text-[var(--text)]"
      onClick={cycleTheme}
      aria-label={`Theme: ${theme}`}
      title={`Theme: ${theme}`}
      suppressHydrationWarning
    >
      <Icon size={16} strokeWidth={1.8} />
    </button>
  );
}

function applyTheme(theme: Theme) {
  const root = document.documentElement;
  if (theme === "system") {
    root.removeAttribute("data-theme");
    return;
  }
  root.dataset.theme = theme;
}
