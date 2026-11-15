package blog

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestSplitJoinRoundTrip(t *testing.T) {
	meta := Meta{Title: `A "quoted" --- title`, Subtitle: "Sub: title", Date: "2026-11-14", Draft: true, Tags: []string{"go", "life"}, ReadTime: 4, Mood: "cloud"}
	body := "First paragraph.\n\n---\n\nAfter a rule."
	gotMeta, gotBody, err := Split(Join(meta, body))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(gotMeta, meta) {
		t.Fatalf("meta = %+v, want %+v", gotMeta, meta)
	}
	if gotBody != body+"\n" {
		t.Fatalf("body = %q", gotBody)
	}
}

// Posts the editor rewrites must still render the same.
func TestJoinKeepsExistingPosts(t *testing.T) {
	files, err := filepath.Glob("../../content/blog/*.md")
	if err != nil || len(files) == 0 {
		t.Fatalf("no posts found: %v", err)
	}
	service, err := New(nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		source, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		meta, body, err := Split(string(source))
		if err != nil {
			t.Fatalf("%s: %v", file, err)
		}
		before, err := service.Render("post", source)
		if err != nil {
			t.Fatal(err)
		}
		after, err := service.Render("post", []byte(Join(meta, body)))
		if err != nil {
			t.Fatalf("%s: %v", file, err)
		}
		if !reflect.DeepEqual(before, after) {
			t.Errorf("%s renders differently after a rewrite", file)
		}
	}
}
