import fs from "node:fs";
import path from "node:path";

const publicDirectory = path.join(process.cwd(), "public");

export function versionedAsset(src: string) {
  if (!src.startsWith("/") || src.includes("?")) {
    return src;
  }

  try {
    const filePath = path.join(publicDirectory, src);
    const stat = fs.statSync(filePath);
    return `${src}?v=${Math.round(stat.mtimeMs)}`;
  } catch {
    return src;
  }
}
