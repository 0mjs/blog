package site

import (
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/0mjs/zinc"
	"github.com/a-h/templ"
	contentfs "mattjs.me/content"
	"mattjs.me/internal/blog"
	"mattjs.me/internal/model"
	"mattjs.me/internal/siteinfo"
	publicfs "mattjs.me/public"
	"mattjs.me/views"
)

var projects = []model.Project{
	{Name: "zinc", Description: "go http framework", URL: "https://zinc.carbonsoft.sh"},
	// {Name: "goff", Description: "go feature flag lib", URL: "https://github.com/0mjs/goff"},
	{Name: "cutwise", Description: "calorie tracker app", URL: "https://cutwise.fit"},
	// {Name: "nuddge", Description: "calm wellness app", URL: "https://nuddge.app"},
	// {Name: "samsa", Description: "ai medtech startup", URL: "https://hih.ie/initiatives/hihi-ai/hihi-ai-call-winners-2025/samsa-ltd/"},
}

// NewApp constructs the blog's shared HTTP handler for local and serverless use.
func NewApp() (*zinc.App, error) {
	config := siteinfo.Current
	blogService, err := blog.New(contentfs.FS)
	if err != nil {
		return nil, err
	}
	images, err := fs.Sub(publicfs.FS, "image")
	if err != nil {
		return nil, err
	}

	app := zinc.New()
	app.UseHTTP(cacheAssets)
	if err := app.StaticFS("/assets", publicfs.FS); err != nil {
		return nil, err
	}
	// Preserve the image URLs already used by the Markdown content.
	if err := app.StaticFS("/image", images); err != nil {
		return nil, err
	}

	app.Get("/", func(c *zinc.Context) error { return render(c, views.Home(blogService.Posts(), projects)) })
	app.Get("/blog", func(c *zinc.Context) error { return render(c, views.BlogList(blogService.Posts())) })
	app.Get("/blog/{slug}", func(c *zinc.Context) error {
		post, ok := blogService.Post(c.Param("slug"))
		if !ok {
			return zinc.ErrNotFound
		}
		return render(c, views.BlogPost(post))
	})
	app.Get("/blog/tag/{tag}", func(c *zinc.Context) error {
		tag, err := url.PathUnescape(c.Param("tag"))
		if err != nil || tag == "" {
			return zinc.ErrNotFound
		}
		return render(c, views.BlogListByTag(blogService.PostsByTag(tag), tag))
	})
	app.Get("/rss.xml", func(c *zinc.Context) error {
		return c.Data("application/rss+xml; charset=utf-8", []byte(rss(blogService.Posts(), config)))
	})
	app.Get("/sitemap.xml", func(c *zinc.Context) error {
		return c.Data("application/xml; charset=utf-8", []byte(sitemap(blogService.Posts(), config)))
	})
	app.Get("/robots.txt", func(c *zinc.Context) error {
		return c.String(fmt.Sprintf("User-agent: *\nAllow: /\nSitemap: %s/sitemap.xml\n", config.BaseURL))
	})
	return app, nil
}

func cacheAssets(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasPrefix(r.URL.Path, "/assets/fonts/"):
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		case r.URL.Path == "/assets/image/matt.png":
			w.Header().Set("Cache-Control", "public, max-age=86400, stale-while-revalidate=604800")
		}
		next.ServeHTTP(w, r)
	})
}

func render(c *zinc.Context, component templ.Component) error {
	c.Type("html")
	return component.Render(c.Context(), c.Writer())
}

func rss(posts []*model.Post, config siteinfo.Config) string {
	var body strings.Builder
	body.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	body.WriteString("<rss version=\"2.0\" xmlns:atom=\"http://www.w3.org/2005/Atom\">\n")
	body.WriteString("  <channel>\n")
	fmt.Fprintf(&body, "    <title>%s</title>\n", xmlEscape(config.Name))
	fmt.Fprintf(&body, "    <link>%s</link>\n", xmlEscape(config.BaseURL))
	fmt.Fprintf(&body, "    <description>%s</description>\n", xmlEscape(config.Tagline))
	fmt.Fprintf(&body, "    <language>%s</language>\n", xmlEscape(config.Language))
	fmt.Fprintf(&body, "    <atom:link href=\"%s/rss.xml\" rel=\"self\" type=\"application/rss+xml\"></atom:link>\n", xmlEscape(config.BaseURL))
	for _, post := range posts {
		postURL := config.BaseURL + "/blog/" + post.Slug
		body.WriteString("    <item>\n")
		fmt.Fprintf(&body, "      <title>%s</title>\n", xmlEscape(post.Title))
		fmt.Fprintf(&body, "      <link>%s</link>\n", xmlEscape(postURL))
		fmt.Fprintf(&body, "      <guid>%s</guid>\n", xmlEscape(postURL))
		fmt.Fprintf(&body, "      <pubDate>%s</pubDate>\n", post.Date.Format(time.RFC1123Z))
		fmt.Fprintf(&body, "      <description>%s</description>\n", xmlEscape(post.Subtitle))
		body.WriteString("    </item>\n")
	}
	body.WriteString("  </channel>\n")
	body.WriteString("</rss>\n")
	return body.String()
}

func sitemap(posts []*model.Post, config siteinfo.Config) string {
	var body strings.Builder
	body.WriteString(`<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)
	fmt.Fprintf(&body, "<url><loc>%s/</loc></url><url><loc>%s/blog</loc></url>", xmlEscape(config.BaseURL), xmlEscape(config.BaseURL))
	for _, post := range posts {
		fmt.Fprintf(&body, "<url><loc>%s/blog/%s</loc><lastmod>%s</lastmod></url>", xmlEscape(config.BaseURL), xmlEscape(post.Slug), post.Date.Format("2006-01-02"))
	}
	body.WriteString("</urlset>")
	return body.String()
}

func xmlEscape(value string) string {
	replacer := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;", "'", "&apos;")
	return replacer.Replace(value)
}

var _ http.Handler = (*zinc.App)(nil)
