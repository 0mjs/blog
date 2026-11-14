// Command editor is a local writing app for the blog. It edits the Markdown posts in content/blog with a live
// preview rendered by the site's own templates, saves images into public/image/blog, and publishes to Cloudflare.
// Run it from the repository root: make write
package main

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"mattjs.me/internal/blog"
	"mattjs.me/internal/siteinfo"
	"mattjs.me/views"
)

//go:embed index.html
var page []byte

const (
	postsDir  = "content/blog"
	imagesDir = "public/image/blog"
)

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

type editor struct {
	renderer   *blog.Service
	publishing sync.Mutex
}

type post struct {
	Slug string    `json:"slug"`
	Meta blog.Meta `json:"meta"`
	Body string    `json:"body"`
}

func main() {
	address := flag.String("addr", "127.0.0.1:4000", "address to listen on")
	open := flag.Bool("open", true, "open the editor in the browser")
	flag.Parse()

	if _, err := os.Stat(postsDir); err != nil {
		log.Fatal("run the editor from the repository root (content/blog not found)")
	}
	renderer, err := blog.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	e := &editor{renderer: renderer}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(page)
	})
	// The site's assets, read from disk so a new image shows up without a rebuild
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("public"))))
	mux.HandleFunc("GET /api/posts", e.list)
	mux.HandleFunc("POST /api/posts", e.create)
	mux.HandleFunc("GET /api/posts/{slug}", e.read)
	mux.HandleFunc("PUT /api/posts/{slug}", e.save)
	mux.HandleFunc("DELETE /api/posts/{slug}", e.remove)
	mux.HandleFunc("POST /api/preview", e.preview)
	mux.HandleFunc("POST /api/images", e.upload)
	mux.HandleFunc("GET /api/status", e.status)
	mux.HandleFunc("POST /api/publish", e.publish)

	listener, err := net.Listen("tcp", *address)
	if err != nil {
		log.Fatal(err)
	}
	url := "http://" + *address
	log.Printf("Editor running at %s (Ctrl+C to stop)", url)
	if *open {
		_ = exec.Command("open", url).Start()
	}
	log.Fatal(http.Serve(listener, guard(mux)))
}

// guard keeps other websites from driving the editor: it only answers on a local address, and anything
// that changes state needs a header a cross-site request can't send.
func guard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host, _, err := net.SplitHostPort(r.Host)
		if err != nil {
			host = r.Host
		}
		if host != "localhost" && host != "127.0.0.1" {
			http.Error(w, "local only", http.StatusForbidden)
			return
		}
		if r.Method != http.MethodGet && r.Header.Get("X-Editor") == "" {
			http.Error(w, "missing X-Editor header", http.StatusForbidden)
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

func postPath(slug string) string { return filepath.Join(postsDir, slug+".md") }

func readPost(slug string) (post, error) {
	source, err := os.ReadFile(postPath(slug))
	if err != nil {
		return post{}, err
	}
	meta, body, err := blog.Split(string(source))
	if err != nil {
		return post{}, fmt.Errorf("%s: %w", slug, err)
	}
	return post{Slug: slug, Meta: meta, Body: body}, nil
}

// writePost replaces the file in one step, so a crash mid-save can't leave half a post.
func writePost(p post) error {
	temporary := postPath(p.Slug) + ".tmp"
	if err := os.WriteFile(temporary, []byte(blog.Join(p.Meta, p.Body)), 0o644); err != nil {
		return err
	}
	return os.Rename(temporary, postPath(p.Slug))
}

func (e *editor) list(w http.ResponseWriter, r *http.Request) {
	files, err := filepath.Glob(filepath.Join(postsDir, "*.md"))
	if err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}
	type summary struct {
		Slug  string `json:"slug"`
		Title string `json:"title"`
		Date  string `json:"date"`
		Draft bool   `json:"draft"`
	}
	posts := []summary{}
	for _, file := range files {
		slug := strings.TrimSuffix(filepath.Base(file), ".md")
		p, err := readPost(slug)
		if err != nil {
			posts = append(posts, summary{Slug: slug, Title: slug + " (can't be read)"})
			continue
		}
		posts = append(posts, summary{Slug: slug, Title: p.Meta.Title, Date: p.Meta.Date, Draft: p.Meta.Draft})
	}
	// Drafts first, then newest first
	sort.SliceStable(posts, func(i, j int) bool {
		if posts[i].Draft != posts[j].Draft {
			return posts[i].Draft
		}
		return posts[i].Date > posts[j].Date
	})
	respond(w, posts)
}

func (e *editor) create(w http.ResponseWriter, r *http.Request) {
	var p post
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		fail(w, err, http.StatusBadRequest)
		return
	}
	if !slugPattern.MatchString(p.Slug) {
		fail(w, errors.New("the address can only use lowercase letters, numbers and dashes"), http.StatusBadRequest)
		return
	}
	if _, err := os.Stat(postPath(p.Slug)); err == nil {
		fail(w, errors.New("a post with that address already exists"), http.StatusConflict)
		return
	}
	p.Meta = blog.Meta{Title: p.Meta.Title, Date: time.Now().Format("2006-01-02"), Draft: true, Tags: []string{}}
	p.Body = ""
	if err := writePost(p); err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}
	respond(w, p)
}

func (e *editor) read(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if !slugPattern.MatchString(slug) {
		fail(w, errors.New("bad address"), http.StatusBadRequest)
		return
	}
	p, err := readPost(slug)
	if err != nil {
		fail(w, err, http.StatusNotFound)
		return
	}
	respond(w, p)
}

func (e *editor) save(w http.ResponseWriter, r *http.Request) {
	var p post
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		fail(w, err, http.StatusBadRequest)
		return
	}
	p.Slug = r.PathValue("slug")
	if !slugPattern.MatchString(p.Slug) {
		fail(w, errors.New("bad address"), http.StatusBadRequest)
		return
	}
	if _, err := os.Stat(postPath(p.Slug)); err != nil {
		fail(w, errors.New("that post no longer exists"), http.StatusNotFound)
		return
	}
	if err := writePost(p); err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}
	respond(w, map[string]bool{"saved": true})
}

// remove deletes a draft and its images. Published posts are unpublished instead, so nothing live vanishes by accident.
func (e *editor) remove(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if !slugPattern.MatchString(slug) {
		fail(w, errors.New("bad address"), http.StatusBadRequest)
		return
	}
	p, err := readPost(slug)
	if err != nil {
		fail(w, err, http.StatusNotFound)
		return
	}
	if !p.Meta.Draft {
		fail(w, errors.New("only drafts can be deleted; unpublish it first"), http.StatusConflict)
		return
	}
	if err := os.Remove(postPath(slug)); err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}
	_ = os.RemoveAll(filepath.Join(imagesDir, slug))
	respond(w, map[string]bool{"deleted": true})
}

// What the preview adds to the real page: links don't navigate away, and nothing waits to animate in.
const previewHead = `<style>.reveal{animation:none!important}</style>
<script>
document.documentElement.classList.add("no-entrance");
addEventListener("click", function (e) {
	var a = e.target.closest && e.target.closest("a");
	if (a && (a.getAttribute("href") || "").charAt(0) !== "#") e.preventDefault();
}, true);
</script>`

func (e *editor) preview(w http.ResponseWriter, r *http.Request) {
	var p post
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		fail(w, err, http.StatusBadRequest)
		return
	}
	if p.Slug == "" {
		p.Slug = "preview"
	}
	if p.Meta.Date == "" {
		p.Meta.Date = time.Now().Format("2006-01-02")
	}
	rendered, err := e.renderer.Render(p.Slug, []byte(blog.Join(p.Meta, p.Body)))
	if err != nil {
		fail(w, err, http.StatusUnprocessableEntity)
		return
	}
	var html bytes.Buffer
	if err := views.BlogPost(rendered).Render(r.Context(), &html); err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = io.WriteString(w, strings.Replace(html.String(), "</head>", previewHead+"</head>", 1))
}

var unsafeName = regexp.MustCompile(`[^a-z0-9]+`)

// Photos are shrunk and re-encoded, which also strips location and camera data. Anything else is kept as it is.
var (
	photos = map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".heic": true}
	asIs   = map[string]bool{".webp": true, ".gif": true, ".svg": true, ".avif": true}
)

func (e *editor) upload(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get("slug")
	if !slugPattern.MatchString(slug) {
		fail(w, errors.New("bad address"), http.StatusBadRequest)
		return
	}
	original := strings.ToLower(filepath.Base(r.URL.Query().Get("name")))
	extension := filepath.Ext(original)
	if !photos[extension] && !asIs[extension] {
		fail(w, fmt.Errorf("%q isn't an image type the blog can show", extension), http.StatusBadRequest)
		return
	}
	data, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 40<<20))
	if err != nil {
		fail(w, errors.New("that image is too big (40 MB limit)"), http.StatusRequestEntityTooLarge)
		return
	}
	name := strings.Trim(unsafeName.ReplaceAllString(strings.TrimSuffix(original, extension), "-"), "-")
	if name == "" {
		name = "image"
	}
	directory := filepath.Join(imagesDir, slug)
	if err := os.MkdirAll(directory, 0o755); err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}

	magick, _ := exec.LookPath("magick")
	convert := photos[extension] && magick != ""
	if extension == ".heic" && !convert {
		fail(w, errors.New("HEIC photos need ImageMagick installed (brew install imagemagick)"), http.StatusBadRequest)
		return
	}
	outputExtension := extension
	if convert {
		outputExtension = ".webp"
	}
	target := filepath.Join(directory, name+outputExtension)
	for n := 2; ; n++ {
		if _, err := os.Stat(target); err != nil {
			break
		}
		target = filepath.Join(directory, fmt.Sprintf("%s-%d%s", name, n, outputExtension))
	}

	if convert {
		temporary, err := os.CreateTemp("", "editor-*"+extension)
		if err != nil {
			fail(w, err, http.StatusInternalServerError)
			return
		}
		defer os.Remove(temporary.Name())
		_, _ = temporary.Write(data)
		_ = temporary.Close()
		output, err := exec.Command(magick, temporary.Name(), "-auto-orient", "-strip", "-resize", "1800x1800>", "-quality", "84", target).CombinedOutput()
		if err != nil {
			fail(w, fmt.Errorf("couldn't process the image: %s", strings.TrimSpace(string(output))), http.StatusUnprocessableEntity)
			return
		}
	} else if err := os.WriteFile(target, data, 0o644); err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}
	respond(w, map[string]string{"url": "/assets/image/blog/" + slug + "/" + filepath.Base(target)})
}

// status reports what a publish would send live besides this post.
func (e *editor) status(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get("slug")
	branch, _ := git("rev-parse", "--abbrev-ref", "HEAD")
	changes, _ := git("status", "--porcelain", "-uall")
	others := []string{}
	for _, line := range strings.Split(changes, "\n") {
		if len(line) < 4 {
			continue
		}
		path := line[3:]
		if _, renamed, ok := strings.Cut(path, " -> "); ok {
			path = renamed
		}
		path = strings.Trim(path, `"`)
		if path == filepath.ToSlash(postPath(slug)) || strings.HasPrefix(path, imagesDir+"/"+slug+"/") {
			continue
		}
		others = append(others, path)
	}
	respond(w, map[string]any{
		"branch": branch,
		"others": others,
		"url":    siteinfo.Current.BaseURL + "/blog/" + slug,
	})
}

func git(args ...string) (string, error) {
	output, err := exec.Command("git", args...).Output()
	return strings.TrimRight(string(output), "\n"), err
}

// stream sends each line to the browser as it's written. It never reports an error, so a closed tab can't
// stop a deploy halfway.
type stream struct {
	w http.ResponseWriter
}

func (s stream) Write(data []byte) (int, error) {
	_, _ = s.w.Write(data)
	if flusher, ok := s.w.(http.Flusher); ok {
		flusher.Flush()
	}
	return len(data), nil
}

// publish saves the post as published (or as a draft, to unpublish), commits it with its images, pushes, rebuilds
// the site and deploys it to Cloudflare. The log is streamed back; the last line starts with DONE or FAILED.
func (e *editor) publish(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Slug      string `json:"slug"`
		Unpublish bool   `json:"unpublish"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil || !slugPattern.MatchString(request.Slug) {
		fail(w, errors.New("bad request"), http.StatusBadRequest)
		return
	}
	if !e.publishing.TryLock() {
		fail(w, errors.New("a publish is already running"), http.StatusConflict)
		return
	}
	defer e.publishing.Unlock()

	before, err := os.ReadFile(postPath(request.Slug))
	if err != nil {
		fail(w, err, http.StatusNotFound)
		return
	}
	p, err := readPost(request.Slug)
	if err != nil {
		fail(w, err, http.StatusUnprocessableEntity)
		return
	}
	if !request.Unpublish && strings.TrimSpace(p.Meta.Title) == "" {
		fail(w, errors.New("give the post a title first"), http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	out := stream{w}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	run := func(title, name string, args ...string) error {
		fmt.Fprintf(out, "\n▸ %s\n", title)
		command := exec.CommandContext(ctx, name, args...)
		command.Stdout, command.Stderr = out, out
		return command.Run()
	}
	failed := func(step string, err error) {
		fmt.Fprintf(out, "\nFAILED: %s (%v)\n", step, err)
	}

	// Was this post already live? Decides the commit message.
	wasLive := false
	if committed, err := git("show", "HEAD:"+filepath.ToSlash(postPath(request.Slug))); err == nil {
		if meta, _, err := blog.Split(committed + "\n"); err == nil {
			wasLive = !meta.Draft
		}
	}
	p.Meta.Draft = request.Unpublish
	if p.Meta.Date == "" {
		p.Meta.Date = time.Now().Format("2006-01-02")
	}
	if err := writePost(p); err != nil {
		failed("saving the post", err)
		return
	}

	if err := run("Checking the site still builds and passes its tests", "go", "test", "./..."); err != nil {
		_ = os.WriteFile(postPath(request.Slug), before, 0o644)
		failed("tests didn't pass, so nothing was published and the post is back as it was", err)
		return
	}

	action := "Publish"
	switch {
	case request.Unpublish:
		action = "Unpublish"
	case wasLive:
		action = "Update post"
	}
	message := action + ": " + p.Meta.Title
	paths := []string{postPath(request.Slug)}
	if _, err := os.Stat(filepath.Join(imagesDir, request.Slug)); err == nil {
		paths = append(paths, filepath.Join(imagesDir, request.Slug))
	}
	if err := run("Saving to git", "git", append([]string{"add", "--"}, paths...)...); err != nil {
		failed("git add", err)
		return
	}
	if pending, _ := git(append([]string{"status", "--porcelain", "--"}, paths...)...); pending != "" {
		if err := run("Committing", "git", append([]string{"commit", "-m", message, "--"}, paths...)...); err != nil {
			failed("git commit", err)
			return
		}
	} else {
		fmt.Fprintln(out, "Nothing new to commit.")
	}
	if err := run("Pushing to GitHub", "git", "push", "origin", "HEAD"); err != nil {
		failed("git push", err)
		return
	}
	if err := run("Building the site", "npm", "run", "build:static"); err != nil {
		failed("building the site", err)
		return
	}
	commit, _ := git("rev-parse", "--short", "HEAD")
	if err := run("Deploying to Cloudflare", "npx", "wrangler", "deploy", "--message", fmt.Sprintf("%s (%s)", message, commit)); err != nil {
		failed("deploying", err)
		return
	}

	url := siteinfo.Current.BaseURL + "/blog/" + request.Slug
	want := http.StatusOK
	if request.Unpublish {
		want = http.StatusNotFound
	}
	fmt.Fprintf(out, "\n▸ Checking %s\n", url)
	if err := waitFor(ctx, url, want); err != nil {
		fmt.Fprintf(out, "Deployed, but the check didn't confirm it yet: %v\n", err)
	} else if request.Unpublish {
		fmt.Fprintln(out, "It's gone from the live site.")
	} else {
		fmt.Fprintln(out, "It's live.")
	}
	fmt.Fprintf(out, "\nDONE %s\n", url)
}

// waitFor polls the live page until it answers with the expected status (Cloudflare can take a few seconds).
func waitFor(ctx context.Context, url string, want int) error {
	client := &http.Client{Timeout: 10 * time.Second}
	var last error
	for attempt := 0; attempt < 8; attempt++ {
		request, _ := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s?v=%d", url, time.Now().UnixNano()), nil)
		response, err := client.Do(request)
		if err == nil {
			_ = response.Body.Close()
			if response.StatusCode == want {
				return nil
			}
			err = fmt.Errorf("got status %d, expected %d", response.StatusCode, want)
		}
		last = err
		time.Sleep(2 * time.Second)
	}
	return last
}

func respond(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(value)
}

func fail(w http.ResponseWriter, err error, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}
