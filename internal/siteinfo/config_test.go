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

func TestLoadBuildFromEnvironment(t *testing.T) {
	t.Setenv("BUILD_COMMIT", "ac11eed5b0a1c2d3")

	got := LoadBuild()
	if got.Commit != "ac11eed" || got.URL != "https://github.com/0mjs/blog/commit/ac11eed5b0a1c2d3" || got.Dirty {
		t.Fatalf("unexpected build: %#v", got)
	}
}
