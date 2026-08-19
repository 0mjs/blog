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

Run `make generate` after editing `.templ` files and `make css` after editing
`assets/css`.
