package content

import "embed"

// FS contains the Markdown posts compiled into the application.
//
//go:embed blog/*.md
var FS embed.FS
