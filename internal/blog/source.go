package blog

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Meta is a post's front matter.
type Meta struct {
	Title    string   `yaml:"title" json:"title"`
	Subtitle string   `yaml:"subtitle" json:"subtitle"`
	Date     string   `yaml:"date" json:"date"`
	Draft    bool     `yaml:"draft" json:"draft"`
	Tags     []string `yaml:"tags" json:"tags"`
	ReadTime int      `yaml:"read_time" json:"read_time"`
}

// Split separates a post's source into its front matter and Markdown body.
func Split(source string) (Meta, string, error) {
	var meta Meta
	head, body, ok := cutFrontMatter(source)
	if !ok {
		return meta, "", errors.New("no front matter")
	}
	if err := yaml.Unmarshal([]byte(head), &meta); err != nil {
		return meta, "", err
	}
	return meta, strings.TrimLeft(body, "\n"), nil
}

// Join writes a post's source in the layout the existing posts use.
func Join(meta Meta, body string) string {
	tags := make([]string, len(meta.Tags))
	for i, tag := range meta.Tags {
		tags[i] = quote(tag)
	}
	var b strings.Builder
	b.WriteString("---\n{\n")
	fmt.Fprintf(&b, "  \"title\": %s,\n", quote(meta.Title))
	fmt.Fprintf(&b, "  \"subtitle\": %s,\n", quote(meta.Subtitle))
	fmt.Fprintf(&b, "  \"date\": %s,\n", quote(meta.Date))
	if meta.ReadTime > 0 {
		fmt.Fprintf(&b, "  \"read_time\": %d,\n", meta.ReadTime)
	}
	fmt.Fprintf(&b, "  \"draft\": %t,\n", meta.Draft)
	fmt.Fprintf(&b, "  \"tags\": [%s],\n", strings.Join(tags, ", "))
	b.WriteString("}\n---\n\n")
	b.WriteString(strings.TrimSpace(body))
	b.WriteString("\n")
	return b.String()
}

// cutFrontMatter finds the block between the opening "---" line and the next "---" line.
func cutFrontMatter(source string) (head, body string, ok bool) {
	rest, ok := strings.CutPrefix(source, "---\n")
	if !ok {
		return "", "", false
	}
	if head, ok := strings.CutSuffix(rest, "\n---"); ok {
		return head, "", true
	}
	return strings.Cut(rest, "\n---\n")
}

func quote(value string) string {
	var b bytes.Buffer
	encoder := json.NewEncoder(&b)
	encoder.SetEscapeHTML(false)
	_ = encoder.Encode(value)
	return strings.TrimSpace(b.String())
}
