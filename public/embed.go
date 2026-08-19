package public

import "embed"

// FS contains the static assets compiled into the application.
//
//go:embed app.css favicon.ico favicon-16x16.png favicon-32x32.png apple-touch-icon.png fonts image
var FS embed.FS
