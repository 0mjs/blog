# Studio

Dark, minimal blog direction on `codex/ui-studio`.

## Research and interpretation

Reference: [Studio design reference](https://www.vitsoe.com/us/about/good-design).

Unobtrusive, useful objects and attention to proportion. One soft graphite enclosure, centred identity, inset link groups and small warm accents; no decorative controls.

This is an original interpretation of those principles, not a clone. Homepage, archive, tag pages and both articles share the design. Original copy, content, routes, RSS, metadata and project links are retained. The supplied coral headshot is used unchanged. Dark presentation is explicit and ignores old light-mode local storage.

## Run

```sh
make run
```

Or generate the production static export with `npm run build:static` after `make generate`.

## Verification

- templ formatting and generation, CSS compilation, existing Go tests, and static export passed.
- Browser-checked home, archive, both articles and tag navigation at desktop and mobile sizes.
- No horizontal overflow at 390px or 320px on the homepage; image loading and keyboard skip navigation checked.
- Desktop, mobile, archive and article screenshots captured for the comparison gallery.
- No new dependencies, content changes, production deployment, or remote branch push.
