# Design

mattjs.me shares its design system with [carbonsoft.sh](https://carbonsoft.sh): Geist and Geist Mono, true black (or
white) with hairline greys, a single railed column with sections stacked inside it, faint grid backgrounds, and the same
motion vocabulary. Coral (`#f04436`) is the personal accent, the way each product has its own accent on carbonsoft.sh.

What stays Matt's: the hex-editor details (`0x00` row offsets, `// projects` labels, the blinking cursor), the hex
portrait that decodes from its own pixels (`public/js/avatar.js`), the byte field behind the hero (`public/js/bytes.js`)
and the ↑↑↓↓←→←→BA god mode.

## Principles

- One column, comfortable reading width; articles are prose first.
- Light and dark themes, remembered per visitor, painted before first render so a refresh never flashes.
- Motion: things rise a little out of a soft blur on a long ease-out, lists stagger, section labels sweep once,
  pills lift with a spring, the footer wordmark rises with scroll. Only transform, opacity and filter animate, and
  `prefers-reduced-motion` turns all of it off.
- Same-origin navigations morph the portrait and post titles between pages (cross-document view transitions).
- Code blocks stay dark in both themes so Chroma's palette reads the same everywhere.
- Unknown pages get a real 404 (Zinc's `NotFound` and error handler locally, `dist/404.html` on Cloudflare).

## Run

```sh
make run
```

Or generate the production static export with `npm run build:static` after `make generate`.
