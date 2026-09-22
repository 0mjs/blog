#!/bin/sh
set -eu

# Run from the repository root. Preserve the headshot's composition and colours.
command -v magick >/dev/null 2>&1 || {
  echo 'ImageMagick is required to regenerate the favicons.' >&2
  exit 1
}

for size in 16 32; do
  magick public/image/matt.png -filter Lanczos -resize "${size}x${size}" -strip "public/favicon-${size}x${size}.png"
done
magick public/image/matt.png -filter Lanczos -resize 180x180 -strip public/apple-touch-icon.png
magick public/image/matt.png -filter Lanczos -define icon:auto-resize=48,32,16 -strip public/favicon.ico
