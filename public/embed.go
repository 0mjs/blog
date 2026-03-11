package publicassets

import "embed"

// FS contains the site's static assets so production doesn't depend on a runtime public/ directory.
//
//go:embed css image js *.png *.ico
var FS embed.FS
