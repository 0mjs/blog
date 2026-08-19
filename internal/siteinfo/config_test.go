package siteinfo

import "testing"

func TestLoadFromEnvironment(t *testing.T) {
	t.Setenv("SITE_NAME", "A Name")
	t.Setenv("SITE_TAGLINE", "A tagline")
	t.Setenv("SITE_URL", "https://example.com/")
	t.Setenv("SITE_LANGUAGE", "en-GB")

	got := Load()
	if got.Name != "A Name" || got.Tagline != "A tagline" || got.BaseURL != "https://example.com" || got.Language != "en-GB" {
		t.Fatalf("unexpected config: %#v", got)
	}
}
