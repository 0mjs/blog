# Matt's blog

A small, server-rendered personal site inspired by axeladrian.com, using:

- Zinc v0.2.1
- templ v0.3.1020
- TemplUI v1.12.0
- Goldmark
- Tailwind CSS v4

```sh
make install
make run
```

Then open <http://localhost:8090>.

After editing a `.templ` file, regenerate its Go source:

```sh
make generate
```

Run `make css` only after changing files in `assets/css` (it requires the
one-time `make install` step).

## Writing

Posts live in `content/blog`. Add a deliberate subtitle to the front matter for
the blog index, RSS feed, and page metadata:

```json
{
  "title": "Post title",
  "subtitle": "A short description of what the post is about.",
  "date": "2026-08-19",
  "tags": ["go"]
}
```

If `subtitle` is omitted, the first paragraph is used as a fallback.
