# mattjs.me

Personal blog built with Zinc, templ, TemplUI, Goldmark, and Tailwind CSS.

## Run

```sh
make install
make run
```

Open <http://localhost:3000>.

For live development, run `make dev` and open <http://localhost:7331>. Changes
to Go, Templ, CSS, and Markdown files rebuild the app and reload the browser.

## Writing

Posts live in `content/blog`. Set the list-page summary with `subtitle` in the
front matter:

```yaml
title: Post title
subtitle: A short description of the post.
date: 2026-08-19
tags: [go]
```

`make dev` regenerates Templ and Tailwind output automatically. Outside the
watcher, use `make generate` for `.templ` changes and `make css` for changes to
`assets/css`.

## Configuration

Favicons are exported directly from `public/image/matt.png`. After replacing the
headshot, run `npm run favicon` (requires ImageMagick). The icon, portrait, and
stylesheet URLs are content-versioned automatically to refresh browser caches.

Site metadata is defined once in `internal/siteinfo/config.go`. Deployments can
override it with `SITE_NAME`, `SITE_TAGLINE`, `SITE_URL`, and `SITE_LANGUAGE`.

## Writing

```sh
make write
```

Opens a local editor at http://127.0.0.1:4000. It edits the Markdown in `content/blog` with a live preview
rendered by the site's own templates, saves as you type, and takes images by drop, paste or the toolbar (they
land in `public/image/blog/<post>/`, shrunk and stripped of location data). Type `/` on a new line to insert
figures, code, tables and the rest. New posts start as drafts; **Publish** commits the post and its images,
pushes, rebuilds and deploys, then checks the live page.

## Cloudflare Workers

The production site is a static export served by Workers Static Assets.

```sh
npm run deploy:dry-run
npm run deploy
```

`wrangler.jsonc` points Workers at `dist`. The deploy scripts rebuild the CSS
and static HTML before uploading, and the exporter writes the Worker `_headers`
rules into the output directory.
