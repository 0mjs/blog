package site

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	htmlstd "html"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/0mjs/zinc"
	"github.com/gomarkdown/markdown"
	mdhtml "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	htmlnode "golang.org/x/net/html"

	publicassets "blog.0mjs.dev/public"
)

//go:embed templates/* content/*
var embeddedContent embed.FS

type Post struct {
	Title    string        `json:"title"`
	Date     time.Time     `json:"date"`
	Draft    bool          `json:"draft,omitempty"`
	Tags     []string      `json:"tags,omitempty"`
	Slug     string        `json:"slug"`
	Content  string        `json:"content"`
	ReadTime int           `json:"read_time"`
	HTML     template.HTML `json:"-"`
}

type Site struct {
	posts      []Post
	postsByID  map[string]Post
	renderer   zinc.Renderer
	chartMu    sync.Mutex
	chart      *GitHubChart
	chartAt    time.Time
	chartRetry time.Time
}

type GitHubChart struct {
	Username   string
	ProfileURL string
	Year       int
	Total      int
	Months     []GitHubChartMonth
	Weeks      []GitHubChartWeek
	Rows       []GitHubChartRow
}

type GitHubChartMonth struct {
	Name    string
	ColSpan int
}

type GitHubChartWeek struct {
	Days []GitHubChartDay
}

type GitHubChartRow struct {
	Label string
	Days  []GitHubChartDay
}

type GitHubChartDay struct {
	Date    string
	Count   int
	Level   int
	Tooltip string
	Empty   bool
}

const (
	frontMatterDelimiter = "---"
	staticMaxAge         = 604800
	githubUsername       = "0mjs"
	githubChartTimeout   = 2 * time.Second
	githubChartTTL       = 12 * time.Hour
	githubChartRetryTTL  = 30 * time.Minute
)

var (
	defaultSite = mustLoadSite()
	staticTypes = map[string]string{
		".css": "text/css; charset=utf-8",
		".js":  "text/javascript; charset=utf-8",
		".ttf": "font/ttf",
	}
	githubChartFetcher = fetchGitHubChart
)

func NewApp() *zinc.App {
	return defaultSite.newApp()
}

func mustLoadSite() *Site {
	entries, err := embeddedContent.ReadDir("content")
	if err != nil {
		log.Fatalf("could not read content directory: %v", err)
	}

	site := &Site{
		postsByID: make(map[string]Post),
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		filePath := filepath.Join("content", entry.Name())
		content, err := embeddedContent.ReadFile(filePath)
		if err != nil {
			log.Printf("error reading content file %s: %v", entry.Name(), err)
			continue
		}

		post, ok := loadPost(entry.Name(), content)
		if !ok {
			continue
		}

		site.postsByID[post.Slug] = post
		if !post.Draft && post.Slug != "about" {
			site.posts = append(site.posts, post)
		}
	}

	sort.Slice(site.posts, func(i, j int) bool {
		return site.posts[i].Date.After(site.posts[j].Date)
	})

	site.renderer = zinc.NewHTMLTemplateRenderer(
		loadTemplates(),
		zinc.WithTemplateSuffixes(".html"),
	)

	return site
}

func loadPost(name string, content []byte) (Post, bool) {
	parts := strings.SplitN(string(content), frontMatterDelimiter, 3)
	if len(parts) != 3 {
		log.Printf("invalid front matter in %s", name)
		return Post{}, false
	}

	var post Post
	if err := json.Unmarshal([]byte(parts[1]), &post); err != nil {
		log.Printf("error parsing front matter in %s: %v", name, err)
		return Post{}, false
	}

	post.Slug = strings.TrimSuffix(name, ".md")
	post.Content = parts[2]
	post.HTML = renderMarkdown(post.Content)

	return post, true
}

func loadTemplates() *template.Template {
	funcMap := template.FuncMap{
		"formatNumber": func(value int) string {
			text := strconv.Itoa(value)
			if len(text) <= 3 {
				return text
			}

			var parts []string
			for len(text) > 3 {
				parts = append([]string{text[len(text)-3:]}, parts...)
				text = text[:len(text)-3]
			}
			parts = append([]string{text}, parts...)
			return strings.Join(parts, ",")
		},
		"tagColor": func(tag string) int {
			hash := 0
			for _, c := range tag {
				hash += int(c)
			}
			if hash < 0 {
				hash = -hash
			}
			return hash % 5
		},
	}

	return template.Must(
		template.New("").Funcs(funcMap).ParseFS(
			embeddedContent,
			"templates/layout.html",
			"templates/index.html",
			"templates/post.html",
			"templates/about.html",
		),
	)
}

func renderMarkdown(content string) template.HTML {
	p := parser.NewWithExtensions(parser.CommonExtensions | parser.AutoHeadingIDs)
	doc := p.Parse([]byte(content))
	opts := mdhtml.RendererOptions{Flags: mdhtml.CommonFlags | mdhtml.HrefTargetBlank}
	return template.HTML(markdown.Render(doc, mdhtml.NewRenderer(opts)))
}

func (s *Site) newApp() *zinc.App {
	app := zinc.NewWithConfig(zinc.Config{
		Renderer: s.renderer,
	})

	s.registerStatic(app)
	must(app.Get("/", s.handleHome))
	must(app.Get("/about", s.handleAbout))
	app.Route("/post", func(posts *zinc.Group) {
		must(posts.Get("/:slug", s.handlePost))
	})

	return app
}

func (s *Site) registerStatic(app *zinc.App) {
	s.registerStaticDir(app, "/css", "css")
	s.registerStaticDir(app, "/image", "image")
	s.registerStaticDir(app, "/js", "js")

	s.registerStaticFile(app, "/favicon.ico", "favicon.ico")
	s.registerStaticFile(app, "/apple-touch-icon.png", "apple-touch-icon.png")
	s.registerStaticFile(app, "/android-chrome-192x192.png", "android-chrome-192x192.png")
	s.registerStaticFile(app, "/android-chrome-512x512.png", "android-chrome-512x512.png")
	s.registerStaticFile(app, "/favicon-16x16.png", "favicon-16x16.png")
	s.registerStaticFile(app, "/favicon-32x32.png", "favicon-32x32.png")
}

func (s *Site) registerStaticDir(app *zinc.App, prefix, dir string) {
	app.UsePrefix(prefix, staticCacheHeaders())
	must(app.Get(prefix+"/*path", func(c *zinc.Context) error {
		assetPath, ok := staticAssetPath(c.Param("*"))
		if !ok {
			return zinc.ErrNotFound
		}
		return serveEmbeddedFile(c, path.Join(dir, assetPath))
	}))
}

func (s *Site) registerStaticFile(app *zinc.App, routePath, filePath string) {
	app.UsePrefix(routePath, staticCacheHeaders())
	must(app.Get(routePath, func(c *zinc.Context) error {
		return serveEmbeddedFile(c, filePath)
	}))
}

func staticAssetPath(value string) (string, bool) {
	cleaned := strings.TrimPrefix(path.Clean("/"+strings.TrimSpace(value)), "/")
	if cleaned == "" || cleaned == "." {
		return "", false
	}
	return cleaned, true
}

func staticCacheHeaders() zinc.HandlerFunc {
	return func(c *zinc.Context) error {
		c.SetHeader("Cache-Control", fmt.Sprintf("public, max-age=%d", staticMaxAge))
		if ct := staticTypes[filepath.Ext(c.Path())]; ct != "" {
			c.SetHeader("Content-Type", ct)
		}
		return c.Next()
	}
}

func serveEmbeddedFile(c *zinc.Context, filePath string) error {
	err := c.FileFS(filePath, publicassets.FS)
	if errors.Is(err, fs.ErrNotExist) || errors.Is(err, fs.ErrInvalid) {
		return zinc.ErrNotFound
	}
	return err
}

func (s *Site) handleHome(c *zinc.Context) error {
	return s.render(c, "index", zinc.Map{
		"Posts":            s.posts,
		"GitHubChart":      s.githubChart(),
		"GitHubProfileURL": githubProfileURL(githubUsername),
	})
}

func (s *Site) handleAbout(c *zinc.Context) error {
	post, ok := s.postsByID["about"]
	if !ok {
		return zinc.ErrNotFound
	}
	return s.render(c, "about", zinc.Map{"Post": post})
}

func (s *Site) handlePost(c *zinc.Context) error {
	post, ok := s.postsByID[c.Param("slug")]
	if !ok || post.Draft {
		return zinc.ErrNotFound
	}
	return s.render(c, "post", zinc.Map{"Post": post})
}

func (s *Site) githubChart() *GitHubChart {
	now := time.Now()

	s.chartMu.Lock()
	cached := s.chart
	fetchedAt := s.chartAt
	retryAt := s.chartRetry
	s.chartMu.Unlock()

	if cached != nil && now.Sub(fetchedAt) < githubChartTTL {
		return cached
	}
	if !retryAt.IsZero() && now.Before(retryAt) {
		return cached
	}

	chart, err := githubChartFetcher(githubUsername, now)
	s.chartMu.Lock()
	defer s.chartMu.Unlock()

	if err != nil {
		log.Printf("error fetching GitHub chart: %v", err)
		s.chartRetry = now.Add(githubChartRetryTTL)
		return s.chart
	}

	s.chart = chart
	s.chartAt = now
	s.chartRetry = time.Time{}
	return s.chart
}

func fetchGitHubChart(username string, now time.Time) (*GitHubChart, error) {
	year := now.Year()
	from := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(year, time.December, 31, 23, 59, 59, 0, time.UTC)

	url := fmt.Sprintf(
		"https://github.com/users/%s/contributions?from=%s&to=%s",
		username,
		from.Format("2006-01-02"),
		to.Format("2006-01-02"),
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "blog.0mjs.dev")

	client := &http.Client{Timeout: githubChartTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github contributions returned status %d", resp.StatusCode)
	}

	days, err := parseGitHubContributionDays(resp.Body)
	if err != nil {
		return nil, err
	}
	if len(days) == 0 {
		return nil, fmt.Errorf("github contributions returned no day cells")
	}

	total := 0
	for _, day := range days {
		total += day.Count
	}

	months, weeks, rows := buildGitHubChartLayout(days, year)

	return &GitHubChart{
		Username:   username,
		ProfileURL: githubProfileURL(username),
		Year:       year,
		Total:      total,
		Months:     months,
		Weeks:      weeks,
		Rows:       rows,
	}, nil
}

func parseGitHubContributionDays(body io.Reader) ([]GitHubChartDay, error) {
	tokenizer := htmlnode.NewTokenizer(body)
	days := make([]GitHubChartDay, 0, 366)
	var (
		inTooltip bool
		tooltip   strings.Builder
	)

	for {
		switch tokenizer.Next() {
		case htmlnode.ErrorToken:
			err := tokenizer.Err()
			if err == io.EOF {
				return days, nil
			}
			return nil, err
		case htmlnode.StartTagToken, htmlnode.SelfClosingTagToken:
			tagName, hasAttr := tokenizer.TagName()
			tag := string(tagName)

			switch tag {
			case "td":
				if !hasAttr {
					continue
				}

				var (
					date  string
					level int
				)
				for {
					key, value, more := tokenizer.TagAttr()
					switch string(key) {
					case "data-date":
						date = string(value)
					case "data-level":
						parsed, err := strconv.Atoi(string(value))
						if err == nil {
							level = parsed
						}
					}
					if !more {
						break
					}
				}

				if date != "" {
					days = append(days, GitHubChartDay{
						Date:  date,
						Level: level,
					})
				}
			case "tool-tip":
				if len(days) == 0 {
					continue
				}
				inTooltip = true
				tooltip.Reset()
			}
		case htmlnode.TextToken:
			if inTooltip {
				tooltip.Write(tokenizer.Text())
			}
		case htmlnode.EndTagToken:
			tagName, _ := tokenizer.TagName()
			if string(tagName) != "tool-tip" || len(days) == 0 {
				continue
			}

			text := strings.TrimSpace(htmlstd.UnescapeString(tooltip.String()))
			days[len(days)-1].Tooltip = text
			days[len(days)-1].Count = parseGitHubTooltipCount(text)
			inTooltip = false
		}
	}
}

func parseGitHubTooltipCount(text string) int {
	if text == "" || strings.HasPrefix(text, "No contributions") {
		return 0
	}

	fields := strings.Fields(text)
	if len(fields) == 0 {
		return 0
	}

	count, err := strconv.Atoi(strings.ReplaceAll(fields[0], ",", ""))
	if err != nil {
		return 0
	}
	return count
}

func githubProfileURL(username string) string {
	return "https://github.com/" + username
}

func buildGitHubChartLayout(days []GitHubChartDay, year int) ([]GitHubChartMonth, []GitHubChartWeek, []GitHubChartRow) {
	yearStart := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(year, time.December, 31, 0, 0, 0, 0, time.UTC)
	gridStart := yearStart.AddDate(0, 0, -int(yearStart.Weekday()))
	gridEnd := yearEnd.AddDate(0, 0, 6-int(yearEnd.Weekday()))

	dayByDate := make(map[string]GitHubChartDay, len(days))
	for _, day := range days {
		dayByDate[day.Date] = day
	}

	weeks := make([]GitHubChartWeek, 0, 54)
	for weekStart := gridStart; !weekStart.After(gridEnd); weekStart = weekStart.AddDate(0, 0, 7) {
		week := GitHubChartWeek{Days: make([]GitHubChartDay, 0, 7)}
		for offset := 0; offset < 7; offset++ {
			dayTime := weekStart.AddDate(0, 0, offset)
			if dayTime.Year() != year {
				week.Days = append(week.Days, GitHubChartDay{Empty: true})
				continue
			}

			date := dayTime.Format("2006-01-02")
			day, ok := dayByDate[date]
			if !ok {
				day = GitHubChartDay{
					Date:    date,
					Tooltip: noContributionsTooltip(dayTime),
				}
			}
			week.Days = append(week.Days, day)
		}
		weeks = append(weeks, week)
	}

	monthColumns := make([]int, 0, 12)
	monthNames := make([]string, 0, 12)
	for month := time.January; month <= time.December; month++ {
		monthStart := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
		column := int(monthStart.Sub(gridStart) / (7 * 24 * time.Hour))
		monthColumns = append(monthColumns, column)
		monthNames = append(monthNames, monthStart.Format("Jan"))
	}

	months := make([]GitHubChartMonth, 0, len(monthColumns))
	for i, column := range monthColumns {
		next := len(weeks)
		if i+1 < len(monthColumns) {
			next = monthColumns[i+1]
		}
		if next <= column {
			continue
		}
		months = append(months, GitHubChartMonth{
			Name:    monthNames[i],
			ColSpan: next - column,
		})
	}

	labels := []string{"", "Mon", "", "Wed", "", "Fri", ""}
	rows := make([]GitHubChartRow, 7)
	for dayIndex := range rows {
		row := GitHubChartRow{
			Label: labels[dayIndex],
			Days:  make([]GitHubChartDay, 0, len(weeks)),
		}
		for _, week := range weeks {
			row.Days = append(row.Days, week.Days[dayIndex])
		}
		rows[dayIndex] = row
	}

	return months, weeks, rows
}

func noContributionsTooltip(day time.Time) string {
	return fmt.Sprintf("No contributions on %s.", day.Format("January 2"))
}

func (s *Site) render(c *zinc.Context, name string, data zinc.Map) error {
	if data == nil {
		data = zinc.Map{}
	}
	data["Now"] = time.Now()
	return c.Render(name, data)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
