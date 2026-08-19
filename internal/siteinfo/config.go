package siteinfo

import (
	"os"
	"strings"
)

// Config contains the site-wide values shared by pages, feeds, and crawlers.
type Config struct {
	Name     string
	Tagline  string
	BaseURL  string
	Language string
}

// Current is loaded once at startup, so deployment environment variables can
// change site metadata without changing source code.
var Current = Load()

// Load returns the site configuration with sensible local defaults.
func Load() Config {
	return Config{
		Name:     envOr("SITE_NAME", "Matt Stevenson"),
		Tagline:  envOr("SITE_TAGLINE", "I build things. Some of them work."),
		BaseURL:  strings.TrimRight(envOr("SITE_URL", "https://mattjs.me"), "/"),
		Language: envOr("SITE_LANGUAGE", "en"),
	}
}

func envOr(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}
