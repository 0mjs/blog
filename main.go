package main

import (
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

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	mdhtml "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
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
	mu       sync.RWMutex
	Pages    map[string]map[string]time.Time `json:"pages"` // page -> visitor_hash -> last_visit
	filePath string
}

var (
	posts     []Post
	postsMap  = make(map[string]Post)
	analytics *Analytics
)

func NewAnalytics(filePath string) *Analytics {
	a := &Analytics{
		Pages:    make(map[string]map[string]time.Time),
		filePath: filePath,
	}
	a.load()

	// Save periodically
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			a.cleanup()
			a.save()
		}
	}()

	return a
}

// Create a unique visitor hash from IP only
func (a *Analytics) getVisitorHash(r *http.Request) string {
	ip := getRealIP(r)
	hash := sha256.Sum256([]byte(ip))
	return fmt.Sprintf("%x", hash[:16]) // Use first 16 bytes
}

// Track a unique visitor to a specific page
func (a *Analytics) Track(page string, r *http.Request) {
	visitorHash := a.getVisitorHash(r)

	a.mu.Lock()
	defer a.mu.Unlock()

	if a.Pages[page] == nil {
		a.Pages[page] = make(map[string]time.Time)
	}

	// Only count if not seen before, or not seen in last 24 hours
	if lastSeen, exists := a.Pages[page][visitorHash]; !exists || time.Since(lastSeen) > 24*time.Hour {
		a.Pages[page][visitorHash] = time.Now()
	}
}

// Get unique visitor count for a specific page
func (a *Analytics) GetUniqueVisitors(page string) int {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if visitors, exists := a.Pages[page]; exists {
		return len(visitors)
	}
	return 0
}

// Cleanup old entries (older than 60 days)
func (a *Analytics) cleanup() {
	a.mu.Lock()
	defer a.mu.Unlock()

	cutoff := time.Now().Add(-60 * 24 * time.Hour)
	for page, visitors := range a.Pages {
		for hash, lastSeen := range visitors {
			if lastSeen.Before(cutoff) {
				delete(visitors, hash)
			}
		}
		// Remove empty pages
		if len(visitors) == 0 {
			delete(a.Pages, page)
		}
	}
}

func (a *Analytics) save() {
	a.mu.RLock()
	defer a.mu.RUnlock()

	data, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		log.Printf("Error marshaling analytics: %v", err)
		return
	}

	if err := os.WriteFile(a.filePath, data, 0644); err != nil {
		log.Printf("Error saving analytics: %v", err)
	}
}

func (a *Analytics) load() {
	data, err := os.ReadFile(a.filePath)
	if err != nil {
		return
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	if err := json.Unmarshal(data, a); err != nil {
		log.Printf("Error loading analytics: %v", err)
	}
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

	// Extract IP from RemoteAddr (format: "IP:port")
	if idx := strings.LastIndex(r.RemoteAddr, ":"); idx != -1 {
		return r.RemoteAddr[:idx]
	}
	return r.RemoteAddr
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

	// Initialize analytics
	analytics = NewAnalytics("analytics.json")
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

	// Save analytics on shutdown
	defer analytics.save()

	log.Println("Server starting on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	analytics.Track("home", r)

	tmpl, err := template.ParseFS(content, "templates/layout.html", "templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, map[string]interface{}{
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

	tmpl, err := template.ParseFS(content, "templates/layout.html", "templates/post.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, map[string]interface{}{
		"Post":     post,
		"Now":      time.Now(),
		"Visitors": analytics.GetUniqueVisitors("post:" + slug),
	})
}

func handleAbout(w http.ResponseWriter, r *http.Request) {
	analytics.Track("about", r)

	tmpl, err := template.ParseFS(content, "templates/layout.html", "templates/about.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, map[string]interface{}{
		"Now":      time.Now(),
		"Visitors": analytics.GetUniqueVisitors("about"),
	})
}
