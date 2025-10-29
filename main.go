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
	"time"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	mdhtml "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/redis/go-redis/v9"
)

//go:embed templates/* static/*
var content embed.FS

type Post struct {
	Title   string        `json:"title"`
	Date    time.Time     `json:"date"`
	Draft   bool          `json:"draft,omitempty"`
	Slug    string        `json:"slug"`
	Content string        `json:"content"`
	HTML    template.HTML `json:"-"`
}

type Analytics struct {
	client *redis.Client
	ctx    context.Context
}

var (
	posts     []Post
	postsMap  = make(map[string]Post)
	analytics *Analytics
	ctx       = context.Background()
	templates = make(map[string]*template.Template)
)

func NewAnalytics(redisURL string) *Analytics {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("Redis URL parse error: %v", err)
	}

	client := redis.NewClient(opt)
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Redis connection error: %v", err)
	}

	log.Println("✅ Redis connected")
	return &Analytics{client: client, ctx: ctx}
}

func (a *Analytics) getIPHash(ip string) string {
	hash := sha256.Sum256([]byte(ip))
	return fmt.Sprintf("%x", hash[:12])
}

func (a *Analytics) Track(page string, r *http.Request) {
	ip := getRealIP(r)
	ipHash := a.getIPHash(ip)
	
	// Add IP hash to page's visitor set (Redis handles uniqueness)
	key := fmt.Sprintf("views:%s", page)
	a.client.SAdd(a.ctx, key, ipHash)
	a.client.Expire(a.ctx, key, 30*24*time.Hour) // 30 day expiry
}

func (a *Analytics) GetUniqueVisitors(page string) int {
	key := fmt.Sprintf("views:%s", page)
	count, _ := a.client.SCard(a.ctx, key).Result()
	return int(count)
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

func init() {
	entries, err := os.ReadDir("content")
	if err != nil {
		log.Printf("Warning: Could not read content directory: %v", err)
		return
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
	templates["index"] = template.Must(template.ParseFS(content, "templates/layout.html", "templates/index.html"))
	templates["post"] = template.Must(template.ParseFS(content, "templates/layout.html", "templates/post.html"))
	templates["about"] = template.Must(template.ParseFS(content, "templates/layout.html", "templates/about.html"))
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

			style := styles.Get("monokai")
			if style == nil {
				style = styles.Fallback
			}

			formatter := html.New(html.WithClasses(true), html.TabWidth(4))
			iterator, err := lexer.Tokenise(nil, string(code.Literal))
			if err != nil {
				io.WriteString(w, string(code.Literal))
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

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	staticFS, _ := fs.Sub(content, "static")
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/post/", handlePost)
	http.HandleFunc("/about", handleAbout)

	log.Println("Server starting on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	analytics.Track("home", r)

	templates["index"].Execute(w, map[string]interface{}{
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

	templates["post"].Execute(w, map[string]interface{}{
		"Post":     post,
		"Now":      time.Now(),
		"Visitors": analytics.GetUniqueVisitors("post:" + slug),
	})
}

func handleAbout(w http.ResponseWriter, r *http.Request) {
	analytics.Track("about", r)

	templates["about"].Execute(w, map[string]interface{}{
		"Now":      time.Now(),
		"Visitors": analytics.GetUniqueVisitors("about"),
	})
}

