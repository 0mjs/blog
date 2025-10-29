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

type VisitorInfo struct {
	IP        string    `json:"ip"`
	Hash      string    `json:"hash"`
	Country   string    `json:"country,omitempty"`
	City      string    `json:"city,omitempty"`
	LastVisit time.Time `json:"last_visit"`
}

type GeoLocation struct {
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	City        string  `json:"city"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Query       string  `json:"query"`
}

var (
	posts     []Post
	postsMap  = make(map[string]Post)
	analytics *Analytics
	ctx       = context.Background()
)

func NewAnalytics(redisURL string) *Analytics {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("Failed to parse Redis URL: %v", err)
	}

	client := redis.NewClient(opt)

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis successfully")

	a := &Analytics{
		client: client,
		ctx:    ctx,
	}

	// Cleanup old entries periodically
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			a.cleanup()
		}
	}()

	return a
}

// Create a unique visitor hash from IP
func (a *Analytics) getVisitorHash(ip string) string {
	hash := sha256.Sum256([]byte(ip))
	return fmt.Sprintf("%x", hash[:16])
}

// Get geolocation data for an IP address
func (a *Analytics) getGeoLocation(ip string) (*GeoLocation, error) {
	// Skip geolocation for localhost/private IPs
	if ip == "127.0.0.1" || ip == "::1" || ip == "localhost" || 
	   strings.HasPrefix(ip, "192.168.") || strings.HasPrefix(ip, "10.") || 
	   strings.HasPrefix(ip, "172.16.") || strings.HasPrefix(ip, "fe80:") {
		return &GeoLocation{Country: "Local", City: "Local"}, nil
	}

	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)
	
	// Create HTTP client with timeout
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("⚠️  Geolocation error for IP %s: %v", ip, err)
		return nil, err
	}
	defer resp.Body.Close()

	var geo GeoLocation
	if err := json.NewDecoder(resp.Body).Decode(&geo); err != nil {
		log.Printf("⚠️  Geolocation decode error for IP %s: %v", ip, err)
		return nil, err
	}

	log.Printf("🌍 IP %s → %s, %s", ip, geo.City, geo.Country)
	return &geo, nil
}

// Track a unique visitor to a specific page
func (a *Analytics) Track(page string, r *http.Request) {
	ip := getRealIP(r)
	log.Printf("📍 Headers: X-Forwarded-For=%s, X-Real-IP=%s, RemoteAddr=%s", 
		r.Header.Get("X-Forwarded-For"), r.Header.Get("X-Real-IP"), r.RemoteAddr)
	log.Printf("📍 Detected IP: %s for page: %s", ip, page)
	
	visitorHash := a.getVisitorHash(ip)
	
	// Check global visitor record (not page-specific)
	globalKey := fmt.Sprintf("visitor:%s", visitorHash)
	var visitor VisitorInfo
	data, err := a.client.Get(a.ctx, globalKey).Result()

	now := time.Now()
	shouldUpdate := false

	if err == redis.Nil {
		// Brand new visitor to the site
		shouldUpdate = true
		log.Printf("🆕 New visitor: %s", visitorHash)
	} else if err == nil {
		// Existing visitor - check if 24 hours have passed
		if err := json.Unmarshal([]byte(data), &visitor); err == nil {
			if now.Sub(visitor.LastVisit) > 24*time.Hour {
				shouldUpdate = true
				log.Printf("🔄 Returning visitor (>24h): %s", visitorHash)
			}
		}
	}

	if shouldUpdate {
		visitor.IP = ip
		visitor.Hash = visitorHash
		visitor.LastVisit = now

		// Get geolocation asynchronously (don't block page load)
		go func() {
			if visitor.Country == "" {
				if geo, err := a.getGeoLocation(ip); err == nil {
					visitor.Country = geo.Country
					visitor.City = geo.City
				} else {
					visitor.Country = "Unknown"
					visitor.City = "Unknown"
				}
			}

			// Save with geolocation data
			visitorJSON, _ := json.Marshal(visitor)
			a.client.Set(a.ctx, globalKey, visitorJSON, 60*24*time.Hour)
		}()

		// Save basic visitor record immediately (without blocking for geo)
		visitorJSON, _ := json.Marshal(visitor)
		a.client.Set(a.ctx, globalKey, visitorJSON, 60*24*time.Hour)

		// Add to global visitors set
		a.client.SAdd(a.ctx, "global:visitors", visitorHash)
		a.client.Expire(a.ctx, "global:visitors", 60*24*time.Hour)
	}

	// Always track which pages this visitor has seen
	pageKey := fmt.Sprintf("visitor:%s:pages", visitorHash)
	a.client.SAdd(a.ctx, pageKey, page)
	a.client.Expire(a.ctx, pageKey, 60*24*time.Hour)
}

// Get total unique visitors across the entire site
func (a *Analytics) GetTotalUniqueVisitors() int {
	count, err := a.client.SCard(a.ctx, "global:visitors").Result()
	if err != nil {
		return 0
	}
	return int(count)
}

// Get unique visitor count for a specific page (unused now, keeping for compatibility)
func (a *Analytics) GetUniqueVisitors(page string) int {
	return a.GetTotalUniqueVisitors()
}

// Get all analytics data (for admin view)
func (a *Analytics) GetAllAnalytics() map[string][]VisitorInfo {
	result := make(map[string][]VisitorInfo)

	// Get all visitor hashes
	visitorHashes, err := a.client.SMembers(a.ctx, "global:visitors").Result()
	if err != nil {
		log.Printf("Error getting visitor hashes: %v", err)
		return result
	}

	for _, hash := range visitorHashes {
		// Get visitor info
		visitorKey := fmt.Sprintf("visitor:%s", hash)
		data, err := a.client.Get(a.ctx, visitorKey).Result()
		if err != nil {
			continue
		}

		var visitor VisitorInfo
		if err := json.Unmarshal([]byte(data), &visitor); err != nil {
			continue
		}

		// Get pages this visitor has seen
		pagesKey := fmt.Sprintf("visitor:%s:pages", hash)
		pages, err := a.client.SMembers(a.ctx, pagesKey).Result()
		if err != nil {
			continue
		}

		// Add visitor to each page they visited
		for _, page := range pages {
			result[page] = append(result[page], visitor)
		}
	}

	return result
}

// Cleanup old entries (older than 60 days)
func (a *Analytics) cleanup() {
	// Redis handles expiration automatically with TTL
	// This is just a placeholder for any additional cleanup logic
	log.Println("Running analytics cleanup (Redis auto-expires old entries)")
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
	http.HandleFunc("/admin/analytics", handleAdminAnalytics)
	http.HandleFunc("/admin/analytics/clear", handleClearAnalytics)

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
		"Visitors": analytics.GetTotalUniqueVisitors(),
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
		"Visitors": analytics.GetTotalUniqueVisitors(),
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
		"Visitors": analytics.GetTotalUniqueVisitors(),
	})
}

func checkAdminAuth(r *http.Request) bool {
	adminSecret := os.Getenv("ADMIN_SECRET")
	if adminSecret == "" {
		adminSecret = "change-me-in-production"
	}

	providedSecret := r.URL.Query().Get("secret")
	if providedSecret == "" {
		providedSecret = r.Header.Get("X-Admin-Secret")
	}

	return providedSecret == adminSecret
}

func handleAdminAnalytics(w http.ResponseWriter, r *http.Request) {
	if !checkAdminAuth(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get all analytics data
	data := analytics.GetAllAnalytics()

	// Calculate totals
	totalVisitors := 0
	pageStats := make(map[string]map[string]int)
	countriesMap := make(map[string]bool)

	for page, visitors := range data {
		totalVisitors += len(visitors)

		// Count visitors by country
		countryCount := make(map[string]int)
		for _, visitor := range visitors {
			country := visitor.Country
			if country == "" {
				country = "Unknown"
			}
			countryCount[country]++
			countriesMap[country] = true
		}
		pageStats[page] = countryCount
	}

	// Check if JSON format is requested
	format := r.URL.Query().Get("format")
	acceptHeader := r.Header.Get("Accept")

	if format == "json" || strings.Contains(acceptHeader, "application/json") {
		// Return JSON response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"total_visitors":  totalVisitors,
			"total_countries": len(countriesMap),
			"pages":           data,
			"page_stats":      pageStats,
			"timestamp":       time.Now(),
		})
		return
	}

	// Return HTML dashboard
	tmpl, err := template.ParseFS(content, "templates/admin.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get admin secret from query or header
	adminSecret := r.URL.Query().Get("secret")
	if adminSecret == "" {
		adminSecret = r.Header.Get("X-Admin-Secret")
	}

	tmpl.Execute(w, map[string]interface{}{
		"TotalVisitors":  totalVisitors,
		"TotalCountries": len(countriesMap),
		"Pages":          data,
		"PageStats":      pageStats,
		"Timestamp":      time.Now(),
		"AdminSecret":    adminSecret,
	})
}

func handleClearAnalytics(w http.ResponseWriter, r *http.Request) {
	if !checkAdminAuth(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Delete all analytics keys
	patterns := []string{"visitor:*", "global:visitors"}
	totalDeleted := 0

	for _, pattern := range patterns {
		keys, err := analytics.client.Keys(analytics.ctx, pattern).Result()
		if err != nil {
			continue
		}

		if len(keys) > 0 {
			if err := analytics.client.Del(analytics.ctx, keys...).Err(); err == nil {
				totalDeleted += len(keys)
			}
		}
	}

	log.Printf("🗑️  Cleared %d analytics keys", totalDeleted)

	// Redirect back to analytics page
	secret := r.URL.Query().Get("secret")
	if secret == "" {
		secret = r.Header.Get("X-Admin-Secret")
	}
	http.Redirect(w, r, "/admin/analytics?secret="+secret, http.StatusSeeOther)
}
