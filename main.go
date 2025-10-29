package main

import (
	"context"
	"crypto/sha256"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	mdhtml "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"

	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
)

//go:embed templates/* static/* public/*
var embeddedContent embed.FS

type Post struct {
	Title   string        `json:"title"`
	Date    time.Time     `json:"date"`
	Draft   bool          `json:"draft,omitempty"`
	Slug    string        `json:"slug"`
	Content string        `json:"content"`
	HTML    template.HTML `json:"-"`
}

type Analytics struct {
	client     *redis.Client
	ctx        context.Context
	countCache map[string]int
	cacheMutex sync.RWMutex
}

var (
	posts     []Post
	postsMap  = make(map[string]Post)
	analytics *Analytics
	ctx       = context.Background()
	templates = make(map[string]*template.Template)
)

func init() {
	entries, err := os.ReadDir("content")
	if err != nil {
		log.Fatalf("Could not read content directory: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		content, err := os.ReadFile(filepath.Join("content", entry.Name()))
		if err != nil {
			log.Printf("Error reading %s: %v", entry.Name(), err)
			continue
		}

		parts := strings.SplitN(string(content), "---", 3)
		if len(parts) != 3 {
			log.Printf("Invalid front matter in %s", entry.Name())
			continue
		}

		var post Post
		if err := json.Unmarshal([]byte(parts[1]), &post); err != nil {
			log.Printf("Error parsing front matter in %s: %v", entry.Name(), err)
			continue
		}

		post.Slug = strings.TrimSuffix(entry.Name(), ".md")
		post.Content = parts[2]
		post.HTML = renderMarkdown(post.Content)

		// Add to map (for direct access checking)
		postsMap[post.Slug] = post

		// Only add to posts list if not a draft
		if !post.Draft {
			posts = append(posts, post)
		}
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.After(posts[j].Date)
	})

	// Initialize analytics with Redis
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379"
	}
	analytics = NewAnalytics(redisURL)

	// Pre-parse templates once
	templates["index"] = template.Must(template.ParseFS(embeddedContent, "templates/layout.html", "templates/index.html"))
	templates["post"] = template.Must(template.ParseFS(embeddedContent, "templates/layout.html", "templates/post.html"))
	templates["about"] = template.Must(template.ParseFS(embeddedContent, "templates/layout.html", "templates/about.html"))
}

func NewAnalytics(redisURL string) *Analytics {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("Redis URL parse error: %v", err)
	}

	opt.MaintNotificationsConfig = &maintnotifications.Config{
		Mode: maintnotifications.ModeDisabled,
	}

	client := redis.NewClient(opt)
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Redis connection error: %v", err)
	}
	log.Println("Redis connected")

	a := &Analytics{
		client:     client,
		ctx:        ctx,
		countCache: make(map[string]int),
	}

	// Initial cache population
	a.refreshCounts()

	// Refresh counts every 10 seconds in background
	go func() {
		for {
			time.Sleep(10 * time.Second)
			a.refreshCounts()
		}
	}()

	return a
}

func (a *Analytics) Track(page string, r *http.Request) {
	// Fire and forget - don't block page render
	go trackIP(page, analytics, r)
}

func (a *Analytics) GetUniqueVisitors(page string) int {
	a.cacheMutex.RLock()
	defer a.cacheMutex.RUnlock()
	return a.countCache[page]
}

func (a *Analytics) getIPHash(ip string) string {
	hash := sha256.Sum256([]byte(ip))
	return fmt.Sprintf("%x", hash[:12])
}

func (a *Analytics) refreshCounts() {
	// Fetch all view counts from Redis
	pages := []string{"home", "about"}

	// Get post slugs from postsMap
	for slug := range postsMap {
		pages = append(pages, "post:"+slug)
	}

	newCache := make(map[string]int)
	for _, page := range pages {
		key := fmt.Sprintf("views:%s", page)
		count, _ := a.client.SCard(a.ctx, key).Result()
		newCache[page] = int(count)
	}

	a.cacheMutex.Lock()
	a.countCache = newCache
	a.cacheMutex.Unlock()
}

// Get real IP from request (handles proxies)
func getRealIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Extract IP from RemoteAddr (format: "IP:port" or "[IPv6]:port")
	ip := r.RemoteAddr

	// Handle IPv6 format [::1]:port
	if strings.HasPrefix(ip, "[") {
		if idx := strings.LastIndex(ip, "]"); idx != -1 {
			ip = ip[1:idx]
		}
	} else {
		// Handle IPv4 format IP:port
		if idx := strings.LastIndex(ip, ":"); idx != -1 {
			ip = ip[:idx]
		}
	}

	return ip
}

func trackIP(page string, a *Analytics, r *http.Request) {
	ip := getRealIP(r)
	ipHash := a.getIPHash(ip)
	key := fmt.Sprintf("views:%s", page)
	a.client.SAdd(a.ctx, key, ipHash)
	a.client.Expire(a.ctx, key, 30*24*time.Hour)
}

func renderMarkdown(content string) template.HTML {
	p := parser.NewWithExtensions(parser.CommonExtensions | parser.AutoHeadingIDs | parser.FencedCode)
	doc := p.Parse([]byte(content))

	opts := mdhtml.RendererOptions{
		Flags: mdhtml.CommonFlags | mdhtml.HrefTargetBlank,
		RenderNodeHook: func(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
			code, ok := node.(*ast.CodeBlock)
			if !ok || !entering {
				return ast.GoToNext, false
			}

			lang := string(code.Info)
			if lang == "" {
				lang = "text"
			}

			lexer := lexers.Get(lang)
			if lexer == nil {
				lexer = lexers.Fallback
			}

			style := styles.Get("gruvbox")
			if style == nil {
				style = styles.Fallback
			}

			formatter := html.New(
				html.PreventSurroundingPre(true),
			)
			iterator, err := lexer.Tokenise(nil, string(code.Literal))
			if err != nil {
				w.Write(code.Literal)
				return ast.GoToNext, true
			}

			io.WriteString(w, `<pre class="chroma"><code class="language-`+lang+`">`)
			formatter.Format(w, style, iterator)
			io.WriteString(w, "</code></pre>")
			return ast.GoToNext, true
		},
	}

	return template.HTML(markdown.Render(doc, mdhtml.NewRenderer(opts)))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	analytics.Track("home", r)

	templates["index"].Execute(w, map[string]any{
		"Posts":    posts,
		"Now":      time.Now(),
		"Visitors": analytics.GetUniqueVisitors("home"),
	})
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/post/")
	post, exists := postsMap[slug]
	if !exists || post.Draft {
		http.NotFound(w, r)
		return
	}

	analytics.Track("post:"+slug, r)

	templates["post"].Execute(w, map[string]any{
		"Post":     post,
		"Now":      time.Now(),
		"Visitors": analytics.GetUniqueVisitors("post:" + slug),
	})
}

func handleAbout(w http.ResponseWriter, r *http.Request) {
	analytics.Track("about", r)

	templates["about"].Execute(w, map[string]any{
		"Now":      time.Now(),
		"Visitors": analytics.GetUniqueVisitors("about"),
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	staticFS, _ := fs.Sub(embeddedContent, "static")
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	publicFS, _ := fs.Sub(embeddedContent, "public")
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.FS(publicFS))))

	// Favicon handler at root path
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		favicon, err := embeddedContent.ReadFile("public/favicon.ico")
		if err != nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "image/x-icon")
		w.Header().Set("Cache-Control", "public, max-age=31536000")
		w.Write(favicon)
	})

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/post/", handlePost)
	http.HandleFunc("/about", handleAbout)

	log.Println("Go server listening on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
