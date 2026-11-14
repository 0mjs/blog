package main

import (
	"crypto/sha256"
	"encoding/xml"
	"fmt"
	"image/jpeg"
	"image/png"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	publicfs "mattjs.me/public"
)

func TestRoutes(t *testing.T) {
	app, err := newApp()
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		path        string
		status      int
		contains    string
		contentType string
	}{
		{"/", http.StatusOK, "Matt Stevenson", "text/html"},
		{"/blog", http.StatusOK, "Gopher Nation", "text/html"},
		{"/blog/gopher-nation", http.StatusOK, "Gopher Nation", "text/html"},
		{"/blog/tag/golang", http.StatusOK, "#golang", "text/html"},
		{"/blog/missing", http.StatusNotFound, "Nothing here", "text/html"},
		{"/no-such-page", http.StatusNotFound, "Nothing here", "text/html"},
		{"/rss.xml", http.StatusOK, "<rss", "application/rss+xml"},
		{"/sitemap.xml", http.StatusOK, "<urlset", "application/xml"},
		{"/robots.txt", http.StatusOK, "Sitemap:", "text/plain"},
		{"/assets/app.css", http.StatusOK, "--color-brand", "text/css"},
		{"/assets/favicon.ico", http.StatusOK, "", "image/x-icon"},
		{"/assets/image/social-card.jpg", http.StatusOK, "", "image/jpeg"},
		{"/assets/image/social-chip.png", http.StatusOK, "", "image/png"},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			app.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, tt.path, nil))
			if recorder.Code != tt.status {
				t.Fatalf("status=%d want=%d body=%s", recorder.Code, tt.status, recorder.Body.String())
			}
			if !strings.Contains(strings.ToLower(recorder.Body.String()), strings.ToLower(tt.contains)) {
				t.Fatalf("body missing %q", tt.contains)
			}
			if !strings.Contains(recorder.Header().Get("Content-Type"), tt.contentType) {
				t.Fatalf("content-type=%q", recorder.Header().Get("Content-Type"))
			}
		})
	}
}

func TestSocialPreview(t *testing.T) {
	app, err := newApp()
	if err != nil {
		t.Fatal(err)
	}
	imagePath := versionedAsset(t, "image/social-card.jpg")
	chipPath := versionedAsset(t, "image/social-chip.png")
	// Pages get the wide portrait card; posts get the square paint chip as a small card
	for _, page := range []struct{ path, image, card, height string }{
		{"/", imagePath, "summary_large_image", "630"},
		{"/blog/gopher-nation", chipPath, "summary", "1200"},
	} {
		response := httptest.NewRecorder()
		app.ServeHTTP(response, httptest.NewRequest(http.MethodGet, page.path, nil))
		for _, tag := range []string{
			`property="og:image" content="https://mattjs.me` + page.image + `"`,
			`name="twitter:image" content="https://mattjs.me` + page.image + `"`,
			`name="twitter:card" content="` + page.card + `"`,
			`property="og:image:width" content="1200"`,
			`property="og:image:height" content="` + page.height + `"`,
		} {
			if !strings.Contains(response.Body.String(), tag) {
				t.Errorf("%s missing %s", page.path, tag)
			}
		}
	}
	chip := httptest.NewRecorder()
	app.ServeHTTP(chip, httptest.NewRequest(http.MethodGet, chipPath, nil))
	if size, err := png.DecodeConfig(chip.Body); err != nil || size.Width != 1200 || size.Height != 1200 {
		t.Fatalf("social chip should be a 1200x1200 PNG: %+v %v", size, err)
	}
	response := httptest.NewRecorder()
	app.ServeHTTP(response, httptest.NewRequest(http.MethodGet, imagePath, nil))
	if response.Code != http.StatusOK || !strings.HasPrefix(response.Header().Get("Content-Type"), "image/jpeg") {
		t.Fatalf("social image response: %d %s", response.Code, response.Header().Get("Content-Type"))
	}
	if response.Body.Len() >= 5_000_000 {
		t.Fatal("social image must be below 5 MB")
	}
	config, err := jpeg.DecodeConfig(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	if config.Width != 1200 || config.Height != 630 {
		t.Fatalf("social image is %dx%d, metadata declares 1200x630", config.Width, config.Height)
	}
}

func TestRSSIsWellFormedAndDiscoverable(t *testing.T) {
	app, err := newApp()
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	app.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/rss.xml", nil))
	feed := recorder.Body.String()

	var document struct {
		XMLName xml.Name
	}
	if err := xml.Unmarshal([]byte(feed), &document); err != nil {
		t.Fatalf("RSS is not well-formed XML: %v\n%s", err, feed)
	}
	for _, expected := range []string{
		`xmlns:atom="http://www.w3.org/2005/Atom"`,
		`<language>en</language>`,
		`<atom:link href="https://mattjs.me/rss.xml" rel="self" type="application/rss+xml"></atom:link>`,
		`<description>Becoming a Gopher.</description>`,
	} {
		if !strings.Contains(feed, expected) {
			t.Fatalf("RSS missing %q", expected)
		}
	}
}

func TestFontsArePreloadedAndCached(t *testing.T) {
	app, err := newApp()
	if err != nil {
		t.Fatal(err)
	}

	home := httptest.NewRecorder()
	app.ServeHTTP(home, httptest.NewRequest(http.MethodGet, "/", nil))
	for _, font := range []string{"Geist.woff2", "GeistMono.woff2"} {
		if !strings.Contains(home.Body.String(), `rel="preload" href="/assets/fonts/Geist/`+font+`"`) {
			t.Errorf("homepage does not preload %s", font)
		}
	}

	font := httptest.NewRecorder()
	app.ServeHTTP(font, httptest.NewRequest(http.MethodGet, "/assets/fonts/Geist/GeistMono.woff2", nil))
	if got := font.Header().Get("Cache-Control"); got != "public, max-age=31536000, immutable" {
		t.Fatalf("font Cache-Control=%q", got)
	}
}

func TestHomepageImageIsPreloadedAndCached(t *testing.T) {
	app, err := newApp()
	if err != nil {
		t.Fatal(err)
	}

	home := httptest.NewRecorder()
	app.ServeHTTP(home, httptest.NewRequest(http.MethodGet, "/", nil))
	if !strings.Contains(home.Body.String(), `rel="preload" href="`+versionedAsset(t, "image/matt.png")+`" as="image"`) {
		t.Fatal("homepage does not preload matt.png")
	}

	image := httptest.NewRecorder()
	app.ServeHTTP(image, httptest.NewRequest(http.MethodGet, "/assets/image/matt.png", nil))
	if got := image.Header().Get("Cache-Control"); got != "public, max-age=86400, stale-while-revalidate=604800" {
		t.Fatalf("image Cache-Control=%q", got)
	}
}

func versionedAsset(t *testing.T, name string) string {
	t.Helper()
	data, err := publicfs.FS.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}
	hash := sha256.Sum256(data)
	return fmt.Sprintf("/assets/%s?v=%x", name, hash[:8])
}

func TestVersionedAssetsAndFavicons(t *testing.T) {
	app, err := newApp()
	if err != nil {
		t.Fatal(err)
	}
	home := httptest.NewRecorder()
	app.ServeHTTP(home, httptest.NewRequest(http.MethodGet, "/", nil))
	for _, asset := range []struct {
		name string
		size int
	}{
		{"app.css", 0},
		{"image/matt.png", 1254},
		{"image/pocket-oranges.jpg", 0},
		{"favicon.ico", 0},
		{"favicon-16x16.png", 16},
		{"favicon-32x32.png", 32},
		{"apple-touch-icon.png", 180},
	} {
		t.Run(asset.name, func(t *testing.T) {
			url := versionedAsset(t, asset.name)
			if !strings.Contains(home.Body.String(), `"`+url+`"`) {
				t.Fatalf("homepage missing versioned asset %s", url)
			}
			response := httptest.NewRecorder()
			app.ServeHTTP(response, httptest.NewRequest(http.MethodGet, url, nil))
			if response.Code != http.StatusOK {
				t.Fatalf("asset status=%d", response.Code)
			}
			if asset.size > 0 {
				config, err := png.DecodeConfig(response.Body)
				if err != nil {
					t.Fatal(err)
				}
				if config.Width != asset.size || config.Height != asset.size {
					t.Fatalf("image is %dx%d, want %dx%d", config.Width, config.Height, asset.size, asset.size)
				}
			}
		})
	}
}
