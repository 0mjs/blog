import fs from "node:fs/promises";
import path from "node:path";
import { cache } from "react";
import MarkdownIt from "markdown-it";
import markdownItAnchor from "markdown-it-anchor";

const contentDirectory = path.join(process.cwd(), "site", "content");

type FrontMatter = {
  title: string;
  date?: string;
  draft?: boolean;
  tags?: string[];
  read_time?: number;
};

export type ContentEntry = {
  slug: string;
  title: string;
  date: Date | null;
  draft: boolean;
  tags: string[];
  readTime: number;
  body: string;
  html: string;
};

const markdown = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
}).use(markdownItAnchor, {
  level: [2, 3],
});

const defaultLinkOpen =
  markdown.renderer.rules.link_open ??
  ((tokens, index, options, _env, self) => self.renderToken(tokens, index, options));

markdown.renderer.rules.link_open = (tokens, index, options, env, self) => {
  const href = tokens[index].attrGet("href") ?? "";

  if (/^https?:\/\//.test(href)) {
    tokens[index].attrSet("target", "_blank");
    tokens[index].attrSet("rel", "noopener noreferrer");
  }

  return defaultLinkOpen(tokens, index, options, env, self);
};

export const getAllEntries = cache(async () => {
  const files = await fs.readdir(contentDirectory);
  const entries = await Promise.all(
    files
      .filter((file) => file.endsWith(".md"))
      .map(async (file) => {
        const raw = await fs.readFile(path.join(contentDirectory, file), "utf8");
        return parseMarkdownFile(file, raw);
      }),
  );

  return entries.sort((a, b) => {
    if (!a.date && !b.date) return a.title.localeCompare(b.title);
    if (!a.date) return 1;
    if (!b.date) return -1;
    return b.date.getTime() - a.date.getTime();
  });
});

export const getPublishedPosts = cache(async () => {
  const entries = await getAllEntries();
  return entries.filter((entry) => !entry.draft && entry.slug !== "about" && entry.date);
});

export const getEntryBySlug = cache(async (slug: string) => {
  const entries = await getAllEntries();
  return entries.find((entry) => entry.slug === slug) ?? null;
});

function parseMarkdownFile(fileName: string, raw: string): ContentEntry {
  const trimmed = raw.trimStart();
  const delimiter = "---";

  if (!trimmed.startsWith(delimiter)) {
    throw new Error(`${fileName} is missing front matter`);
  }

  const endIndex = trimmed.indexOf(`\n${delimiter}`, delimiter.length);
  if (endIndex === -1) {
    throw new Error(`${fileName} has unterminated front matter`);
  }

  const frontMatterText = trimmed.slice(delimiter.length, endIndex).trim();
  const body = trimmed.slice(endIndex + delimiter.length + 1).trimStart();
  const frontMatter = JSON.parse(frontMatterText) as FrontMatter;
  const slug = fileName.replace(/\.md$/, "");

  if (!frontMatter.title) {
    throw new Error(`${fileName} is missing a title`);
  }

  return {
    slug,
    title: frontMatter.title,
    date: frontMatter.date ? new Date(frontMatter.date) : null,
    draft: frontMatter.draft ?? false,
    tags: frontMatter.tags ?? [],
    readTime: frontMatter.read_time ?? estimateReadTime(body),
    body,
    html: markdown.render(body),
  };
}

function estimateReadTime(markdownBody: string) {
  const words = markdownBody
    .replace(/<[^>]*>/g, " ")
    .split(/\s+/)
    .filter(Boolean).length;

  return Math.max(1, Math.ceil(words / 220));
}
