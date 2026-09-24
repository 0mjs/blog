package main

import (
	"strings"
	"testing"
)

// Post embeds (e.g. the Rob Pike talk in gopher-nation.md) need an explicit
// frame-src, or default-src 'self' blocks them in production.
func TestHeadersAllowYouTubeEmbeds(t *testing.T) {
	if !strings.Contains(headers, "frame-src https://www.youtube.com https://www.youtube-nocookie.com;") {
		t.Fatalf("CSP does not allow YouTube iframes:\n%s", headers)
	}
}
