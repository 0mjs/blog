# mattjs.me

Personal blog built with Zinc, templ, TemplUI, Goldmark, and Tailwind CSS.

## Run

```sh
make install
make run
```

Open <http://localhost:3000>.

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
