// Command export renders the blog as a static site for Cloudflare Workers.
package main

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/a-h/templ"
	contentfs "mattjs.me/content"
	"mattjs.me/internal/blog"
	"mattjs.me/internal/model"
	"mattjs.me/internal/siteinfo"
	publicfs "mattjs.me/public"
	"mattjs.me/site"
	"mattjs.me/views"
)

const outputDir = "dist"

func main() {
	if err := export(outputDir); err != nil {
		fmt.Fprintln(os.Stderr, "static export failed:", err)
		os.Exit(1)
	}
}

func export(output string) error {
	if err := os.RemoveAll(output); err != nil {
		return err
	}
	if err := os.MkdirAll(output, 0o755); err != nil {
		return err
	}

	posts, err := blog.New(contentfs.FS)
	if err != nil {
		return err
	}
	allPosts := posts.Posts()
	config := siteinfo.Current

	pages := []struct {
		path      string
		component templ.Component
	}{
		{"/", views.Home(allPosts, site.Projects)},
		{"/blog", views.BlogList(allPosts)},
	}
	for _, post := range allPosts {
		pages = append(pages, struct {
			path      string
			component templ.Component
		}{"/blog/" + post.Slug, views.BlogPost(post)})
	}
	for _, tag := range tags(allPosts) {
		pages = append(pages, struct {
			path      string
			component templ.Component
		}{"/blog/tag/" + tag, views.BlogListByTag(posts.PostsByTag(tag), tag)})
	}
	for _, page := range pages {
		if err := writeHTML(output, page.path, page.component); err != nil {
			return err
		}
	}

	if err := writeFile(output, "rss.xml", []byte(site.RSS(allPosts, config))); err != nil {
		return err
	}
	if err := writeFile(output, "sitemap.xml", []byte(site.Sitemap(allPosts, config))); err != nil {
		return err
	}
	robots := fmt.Sprintf("User-agent: *\nAllow: /\nSitemap: %s/sitemap.xml\n", config.BaseURL)
	if err := writeFile(output, "robots.txt", []byte(robots)); err != nil {
		return err
	}
	if err := writeFile(output, "_headers", []byte(headers)); err != nil {
		return err
	}
	return copyAssets(output)
}

const headers = `/*
  X-Content-Type-Options: nosniff
  Referrer-Policy: strict-origin-when-cross-origin
  Permissions-Policy: camera=(), microphone=(), geolocation=()
  Content-Security-Policy: default-src 'self'; object-src 'none'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'; connect-src 'self'; frame-ancestors 'self'; base-uri 'self'; form-action 'self'

/assets/*
  Cache-Control: public, max-age=86400, stale-while-revalidate=604800

/assets/fonts/*
  Cache-Control: public, max-age=31536000, immutable
`

func writeHTML(output, route string, component templ.Component) error {
	var rendered bytes.Buffer
	if err := component.Render(context.Background(), &rendered); err != nil {
		return err
	}
	path := "index.html"
	if route != "/" {
		path = filepath.Join(strings.TrimPrefix(route, "/"), "index.html")
	}
	return writeFile(output, path, rendered.Bytes())
}

func writeFile(output, path string, content []byte) error {
	filename := filepath.Join(output, filepath.FromSlash(path))
	if err := os.MkdirAll(filepath.Dir(filename), 0o755); err != nil {
		return err
	}
	return os.WriteFile(filename, content, 0o644)
}

func copyAssets(output string) error {
	return fs.WalkDir(publicfs.FS, ".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == "." || entry.IsDir() {
			return nil
		}
		content, err := fs.ReadFile(publicfs.FS, path)
		if err != nil {
			return err
		}
		return writeFile(output, filepath.ToSlash(filepath.Join("assets", path)), content)
	})
}

func tags(posts []*model.Post) []string {
	seen := make(map[string]struct{})
	for _, post := range posts {
		for _, tag := range post.Tags {
			seen[tag] = struct{}{}
		}
	}
	values := make([]string, 0, len(seen))
	for tag := range seen {
		values = append(values, tag)
	}
	sort.Strings(values)
	return values
}
