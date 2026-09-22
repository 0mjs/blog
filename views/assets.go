package views

import (
	"crypto/sha256"
	"fmt"

	publicfs "mattjs.me/public"
)

var (
	stylesheetURL     = assetURL("app.css")
	avatarURL         = assetURL("image/matt.png")
	faviconURL        = assetURL("favicon.ico")
	favicon32URL      = assetURL("favicon-32x32.png")
	favicon16URL      = assetURL("favicon-16x16.png")
	appleTouchIconURL = assetURL("apple-touch-icon.png")
)

// Version mutable assets by their content so a deployment bypasses old browser caches.
func assetURL(name string) string {
	data, err := publicfs.FS.ReadFile(name)
	if err != nil {
		panic(fmt.Sprintf("missing site asset %q: %v", name, err))
	}
	hash := sha256.Sum256(data)
	return fmt.Sprintf("/assets/%s?v=%x", name, hash[:8])
}
