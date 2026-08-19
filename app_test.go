package main

import (
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
		{"/blog", http.StatusOK, "constraints and simplicity keep me focused on building.", "text/html"},
		{"/blog/why-go", http.StatusOK, "The Thing About TypeScript", "text/html"},
		{"/blog/tag/go", http.StatusOK, "#go", "text/html"},
		{"/blog/missing", http.StatusNotFound, "Not Found", "text/plain"},
		{"/rss.xml", http.StatusOK, "<rss", "application/rss+xml"},
		{"/sitemap.xml", http.StatusOK, "<urlset", "application/xml"},
		{"/robots.txt", http.StatusOK, "Sitemap:", "text/plain"},
		{"/assets/app.css", http.StatusOK, "--color-brand", "text/css"},
		{"/assets/favicon.ico", http.StatusOK, "", "image/x-icon"},
		{"/image/meme/go-ts-node.jpg", http.StatusOK, "", "image/jpeg"},
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
		`Why Go&apos;s constraints and simplicity keep me focused on building.`,
	} {
		if !strings.Contains(feed, expected) {
			t.Fatalf("RSS missing %q", expected)
		}
	}
}
