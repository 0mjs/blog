package blog

import (
	"bytes"
	"fmt"
	"html"
	"io/fs"
	"regexp"
	"sort"
	"strings"
	"time"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	goldmarkhtml "github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
	"go.abhg.dev/goldmark/frontmatter"
	"mattjs.me/internal/model"
)

type Service struct {
	markdown goldmark.Markdown
	content  fs.FS
	posts    []*model.Post
	bySlug   map[string]*model.Post
}

var markup = regexp.MustCompile(`(?s)<[^>]*>|!\[[^]]*\]\([^)]*\)|\[([^]]+)\]\([^)]*\)|[*_#>` + "`" + `~]`)

func New(content fs.FS) (*Service, error) {
	service := &Service{
		markdown: goldmark.New(
			goldmark.WithExtensions(extension.GFM, extension.Footnote, extension.Typographer, &frontmatter.Extender{},
				highlighting.NewHighlighting(
					highlighting.WithFormatOptions(chromahtml.WithClasses(true)),
					highlighting.WithWrapperRenderer(codeBlock),
				),
			),
			goldmark.WithParserOptions(parser.WithAutoHeadingID()),
			goldmark.WithRendererOptions(goldmarkhtml.WithUnsafe()),
		),
		content: content,
		bySlug:  make(map[string]*model.Post),
	}
	// No content: a service that only renders what it's handed (the local editor's preview)
	if content == nil {
		return service, nil
	}

	files, err := fs.Glob(content, "blog/*.md")
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		post, err := service.load(file)
		if err != nil {
			return nil, err
		}
		service.bySlug[post.Slug] = post
		if !post.Draft && !post.Date.IsZero() {
			service.posts = append(service.posts, post)
		}
	}
	sort.Slice(service.posts, func(i, j int) bool { return service.posts[i].Date.After(service.posts[j].Date) })
	return service, nil
}

// codeBlock wraps each fenced block so the page can label it with its language.
// Highlighting uses CSS classes; the token colours live in assets/css/input.css.
func codeBlock(w util.BufWriter, c highlighting.CodeBlockContext, entering bool) {
	if entering {
		_, _ = w.WriteString(`<div class="code"`)
		if lang, ok := c.Language(); ok && len(lang) > 0 {
			_, _ = w.WriteString(` data-lang="` + html.EscapeString(string(lang)) + `"`)
		}
		_, _ = w.WriteString(">")
		if !c.Highlighted() {
			_, _ = w.WriteString("<pre><code>")
		}
		return
	}
	if !c.Highlighted() {
		_, _ = w.WriteString("</code></pre>")
	}
	_, _ = w.WriteString("</div>\n")
}

func (s *Service) Posts() []*model.Post {
	return append([]*model.Post(nil), s.posts...)
}

func (s *Service) Post(slug string) (*model.Post, bool) {
	post, ok := s.bySlug[slug]
	return post, ok && !post.Draft
}

func (s *Service) PostsByTag(tag string) []*model.Post {
	var posts []*model.Post
	for _, post := range s.posts {
		for _, candidate := range post.Tags {
			if strings.EqualFold(candidate, tag) {
				posts = append(posts, post)
				break
			}
		}
	}
	return posts
}

func (s *Service) load(file string) (*model.Post, error) {
	source, err := fs.ReadFile(s.content, file)
	if err != nil {
		return nil, err
	}
	return s.Render(strings.TrimSuffix(strings.TrimPrefix(file, "blog/"), ".md"), source)
}

// Render turns one post's source (front matter and Markdown) into a post, exactly as the site does.
func (s *Service) Render(slug string, source []byte) (*model.Post, error) {
	ctx := parser.NewContext()
	var rendered bytes.Buffer
	if err := s.markdown.Convert(source, &rendered, parser.WithContext(ctx)); err != nil {
		return nil, fmt.Errorf("render %s: %w", slug, err)
	}

	var meta Meta
	data := frontmatter.Get(ctx)
	if data == nil {
		return nil, fmt.Errorf("%s has no front matter", slug)
	}
	if err := data.Decode(&meta); err != nil {
		return nil, fmt.Errorf("decode %s: %w", slug, err)
	}

	date, err := parseDate(meta.Date)
	if err != nil {
		return nil, fmt.Errorf("date in %s: %w", slug, err)
	}
	body := bodyAfterFrontMatter(string(source))
	subtitle := meta.Subtitle
	if subtitle == "" {
		subtitle = summarize(body)
	}
	readTime := meta.ReadTime
	if readTime == 0 {
		readTime = max(1, (len(strings.Fields(body))+199)/200)
	}

	return &model.Post{
		Title:    meta.Title,
		Subtitle: subtitle,
		Slug:     slug,
		Date:     date,
		Tags:     meta.Tags,
		HTML:     rendered.String(),
		ReadTime: readTime,
		Draft:    meta.Draft,
		Mood:     meta.Mood,
	}, nil
}

func parseDate(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}
	if parsed, err := time.Parse(time.RFC3339, value); err == nil {
		return parsed, nil
	}
	return time.Parse("2006-01-02", value)
}

func bodyAfterFrontMatter(source string) string {
	if _, body, ok := cutFrontMatter(source); ok {
		return strings.TrimSpace(body)
	}
	parts := strings.SplitN(source, "---", 3)
	if len(parts) != 3 {
		return source
	}
	return strings.TrimSpace(parts[2])
}

func summarize(body string) string {
	paragraph := strings.SplitN(body, "\n\n", 2)[0]
	paragraph = markup.ReplaceAllString(paragraph, "$1")
	paragraph = strings.Join(strings.Fields(paragraph), " ")
	const limit = 160
	if len(paragraph) <= limit {
		return paragraph
	}
	return strings.TrimSpace(paragraph[:limit]) + "…"
}
