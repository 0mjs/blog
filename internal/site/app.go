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
	publicfs "mattjs.me/public"
	"mattjs.me/views"
)

var projects = []model.Project{
	{Name: "zinc", Description: "go http framework", URL: "https://zinc.carbonsoft.sh"},
	{Name: "cutwise", Description: "ai calorie tracker", URL: "https://cutwise.fit"},
	// {Name: "nuddge", Description: "calm wellness for combatting overwhelm", URL: "https://nuddge.app"},
	// {Name: "samsa", Description: "AI MedTech system", URL: "https://hih.ie/initiatives/hihi-ai/hihi-ai-call-winners-2025/samsa-ltd/"},
}

// NewApp constructs the blog's shared HTTP handler for local and serverless use.
func NewApp() (*zinc.App, error) {
	blogService, err := blog.New(contentfs.FS)
	if err != nil {
		return nil, err
	}
	images, err := fs.Sub(publicfs.FS, "image")
	if err != nil {
		return nil, err
	}

	app := zinc.New()
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
		return c.Data("application/rss+xml; charset=utf-8", []byte(rss(blogService.Posts())))
	})
	app.Get("/sitemap.xml", func(c *zinc.Context) error {
		return c.Data("application/xml; charset=utf-8", []byte(sitemap(blogService.Posts())))
	})
	app.Get("/robots.txt", func(c *zinc.Context) error {
		return c.String("User-agent: *\nAllow: /\nSitemap: https://mattjs.me/sitemap.xml\n")
	})
	return app, nil
}

func render(c *zinc.Context, component templ.Component) error {
	c.Type("html")
	return component.Render(c.Context(), c.Writer())
}

func rss(posts []*model.Post) string {
	var body strings.Builder
	body.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	body.WriteString("<rss version=\"2.0\" xmlns:atom=\"http://www.w3.org/2005/Atom\">\n")
	body.WriteString("  <channel>\n")
	body.WriteString("    <title>Matt Stevenson</title>\n")
	body.WriteString("    <link>https://mattjs.me</link>\n")
	body.WriteString("    <description>I build schtuff.</description>\n")
	body.WriteString("    <language>en</language>\n")
	body.WriteString("    <atom:link href=\"https://mattjs.me/rss.xml\" rel=\"self\" type=\"application/rss+xml\"></atom:link>\n")
	for _, post := range posts {
		body.WriteString("    <item>\n")
		fmt.Fprintf(&body, "      <title>%s</title>\n", xmlEscape(post.Title))
		fmt.Fprintf(&body, "      <link>https://mattjs.me/blog/%s</link>\n", post.Slug)
		fmt.Fprintf(&body, "      <guid>https://mattjs.me/blog/%s</guid>\n", post.Slug)
		fmt.Fprintf(&body, "      <pubDate>%s</pubDate>\n", post.Date.Format(time.RFC1123Z))
		fmt.Fprintf(&body, "      <description>%s</description>\n", xmlEscape(post.Subtitle))
		body.WriteString("    </item>\n")
	}
	body.WriteString("  </channel>\n")
	body.WriteString("</rss>\n")
	return body.String()
}

func sitemap(posts []*model.Post) string {
	var body strings.Builder
	body.WriteString(`<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"><url><loc>https://mattjs.me/</loc></url><url><loc>https://mattjs.me/blog</loc></url>`)
	for _, post := range posts {
		fmt.Fprintf(&body, "<url><loc>https://mattjs.me/blog/%s</loc><lastmod>%s</lastmod></url>", post.Slug, post.Date.Format("2006-01-02"))
	}
	body.WriteString("</urlset>")
	return body.String()
}

func xmlEscape(value string) string {
	replacer := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;", "'", "&apos;")
	return replacer.Replace(value)
}

var _ http.Handler = (*zinc.App)(nil)
