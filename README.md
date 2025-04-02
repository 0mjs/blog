# 0mjs

A minimalist blog built with Go, inspired by [Ryan Dahl's blog](https://tinyclouds.org/). This blog focuses on software engineering, systems, and technology with a clean, modern aesthetic.

## Features

- Built with Go 1.22+ and the standard library
- Markdown-based content management
- Dark mode by default
- No JavaScript (except where absolutely necessary)
- Fast and lightweight
- Mobile-responsive design

## Getting Started

1. Clone the repository:
   ```bash
   git clone https://github.com/0mjs.git
   cd 0mjs
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run the server:
   ```bash
   go run main.go
   ```

4. Visit `http://localhost:8080` in your browser

## Adding Content

Blog posts are written in Markdown format and stored in the `content/` directory. Each post should have a front matter section with JSON metadata:

```markdown
---
{
    "title": "Your Post Title",
    "date": "2024-03-20T00:00:00Z"
}
---

Your content here...
```

## Development

The blog is built with:

- Go 1.22+ for the server
- HTML templates for rendering
- CSS for styling
- Markdown for content

## License

MIT License - feel free to use this code for your own projects. 