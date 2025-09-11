package main

import (
	"embed"
	"encoding/json"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/alecthomas/chroma"
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
	Slug    string        `json:"slug"`
	Content string        `json:"content"`
	HTML    template.HTML `json:"-"`
}

var (
	posts        []Post
	postsMap     = make(map[string]Post)
	weeklyStatus = "Abandoning more side-projects"
)

// Custom template functions
var templateFuncs = template.FuncMap{
	"truncate": func(s string, n interface{}) string {
		var limit int
		switch v := n.(type) {
		case int:
			limit = v
		case int64:
			limit = int(v)
		case float64:
			limit = int(v)
		case string:
			// Try to convert string to int
			parsed, err := strconv.Atoi(v)
			if err != nil {
				return s // Return original string if conversion fails
			}
			limit = parsed
		default:
			return s // Return original string for any other type
		}

		if len(s) <= limit {
			return s
		}
		for !utf8.ValidString(s[:limit]) {
			limit--
		}
		return s[:limit] + "..."
	},
	"preview": func(s string) string {
		// Create a simple preview of around 150 characters
		s = strings.TrimSpace(s)
		if len(s) <= 150 {
			return s
		}
		for i := 150; i > 0; i-- {
			if utf8.ValidString(s[:i]) {
				// Try to end at a space or punctuation
				lastPeriod := strings.LastIndex(s[:i], ".")
				lastQuestion := strings.LastIndex(s[:i], "?")
				lastExclamation := strings.LastIndex(s[:i], "!")
				lastSpace := strings.LastIndex(s[:i], " ")

				// Use the closest sentence-ending punctuation or space
				endPos := i
				if lastPeriod > endPos-30 && lastPeriod > 0 {
					endPos = lastPeriod + 1
				} else if lastQuestion > endPos-30 && lastQuestion > 0 {
					endPos = lastQuestion + 1
				} else if lastExclamation > endPos-30 && lastExclamation > 0 {
					endPos = lastExclamation + 1
				} else if lastSpace > 0 {
					endPos = lastSpace
				}

				return s[:endPos] + "..."
			}
		}
		return s[:100] + "..." // Fallback
	},
	"countWords": func(s string) int {
		return len(strings.Fields(s))
	},
	"readingTime": func(words int) int {
		// Assuming average reading speed of 200 words per minute
		minutes := words / 200
		if minutes < 1 {
			return 1
		}
		return minutes
	},
}

func init() {
	// Load posts from content directory
	entries, err := os.ReadDir("content")
	if err != nil {
		log.Printf("Warning: Could not read content directory: %v", err)
		// Initialize with empty posts rather than crashing
		posts = []Post{}
		return
	}

	// Initialize posts with an empty slice instead of nil
	posts = []Post{}
	// Initialize postsMap
	postsMap = make(map[string]Post)

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			content, err := os.ReadFile(filepath.Join("content", entry.Name()))
			if err != nil {
				log.Printf("Error reading %s: %v", entry.Name(), err)
				continue
			}

			// Parse front matter and content
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

			// Convert markdown to HTML
			extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.FencedCode
			p := parser.NewWithExtensions(extensions)
			doc := p.Parse([]byte(post.Content))

			// Create HTML renderer with flags
			htmlFlags := mdhtml.CommonFlags | mdhtml.HrefTargetBlank
			opts := mdhtml.RendererOptions{
				Flags: htmlFlags,
				RenderNodeHook: func(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
					if code, ok := node.(*ast.CodeBlock); ok && entering {
						lang := string(code.Info)
						if lang == "" {
							lang = "text"
						}

						// Get lexer for the language
						lexer := lexers.Get(lang)
						if lexer == nil {
							lexer = lexers.Fallback
						}
						lexer = chroma.Coalesce(lexer)

						// Use a custom style that matches our CSS
						style := styles.Get("monokai")
						if style == nil {
							style = styles.Fallback
						}

						formatter := html.New(
							html.WithClasses(true),
							html.TabWidth(4),
						)

						iterator, err := lexer.Tokenise(nil, string(code.Literal))
						if err != nil {
							io.WriteString(w, string(code.Literal))
							return ast.GoToNext, true
						}

						// Write the highlighted code
						io.WriteString(w, `<pre class="chroma"><code class="language-`+lang+`">`)
						formatter.Format(w, style, iterator)
						io.WriteString(w, "</code></pre>")
						return ast.GoToNext, true
					}
					return ast.GoToNext, false
				},
			}
			renderer := mdhtml.NewRenderer(opts)

			post.HTML = template.HTML(markdown.Render(doc, renderer))
			posts = append(posts, post)

			// Add to postsMap
			postsMap[post.Slug] = post
		}
	}

	// Only sort if we have posts
	if len(posts) > 0 {
		// Sort posts by date
		sort.Slice(posts, func(i, j int) bool {
			return posts[i].Date.After(posts[j].Date)
		})
	}
}

func main() {
	// Serve static files
	staticFS, _ := fs.Sub(content, "static")
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// Routes
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/post/", handlePost)
	http.HandleFunc("/about", handleAbout)

	// Start server
	log.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.New("layout.html").Funcs(templateFuncs).ParseFS(content, "templates/layout.html", "templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, map[string]interface{}{
		"Posts":        posts,
		"Now":          time.Now(),
		"CurrentRoute": "/",
		"WeeklyStatus": weeklyStatus,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/post/")

	post, exists := postsMap[slug]
	if !exists {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.New("layout.html").Funcs(templateFuncs).ParseFS(content, "templates/layout.html", "templates/post.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, map[string]interface{}{
		"Post":         post,
		"Now":          time.Now(),
		"CurrentRoute": "/post/" + slug,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleAbout(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("layout.html").Funcs(templateFuncs).ParseFS(content, "templates/layout.html", "templates/about.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, map[string]interface{}{
		"Now":          time.Now(),
		"CurrentRoute": "/about",
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
