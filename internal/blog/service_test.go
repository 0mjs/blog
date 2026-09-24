package blog

import (
	"strings"
	"testing"
	"testing/fstest"
)

func TestCodeBlocks(t *testing.T) {
	post := "---\n{\"title\": \"Code\", \"date\": \"2026-01-01\"}\n---\n\n```go\nfunc main() { fmt.Println(\"hi\") }\n```\n\n```\n<raw> & plain\n```\n"
	service, err := New(fstest.MapFS{"blog/code.md": {Data: []byte(post)}})
	if err != nil {
		t.Fatal(err)
	}
	html := service.Posts()[0].HTML
	for _, want := range []string{
		`<div class="code" data-lang="go"><pre class="chroma">`,
		`<span class="kd">func</span>`,
		`<span class="s">&#34;hi&#34;</span>`,
		`<div class="code"><pre><code>&lt;raw&gt; &amp; plain`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("rendered HTML missing %q:\n%s", want, html)
		}
	}
}
