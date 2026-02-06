# QuickBlog

[中文文档](README_CN.md)

A lightweight static blog server written in Go. Drop `.html` or `.md` files into a directory and they're instantly served as a blog — no database, no config files, no framework.

## Features

- **Zero config** — just a single binary + a `pages/` directory
- **Markdown support** — write in Markdown, convert to HTML with one command (GFM: tables, task lists, code blocks, strikethrough)
- **Auto-sorted** — homepage lists all pages by modification time, newest first
- **Clean design** — warm minimalist template with Inter font and responsive layout
- **Easy deploy** — single binary with embedded templates, cross-compile to any platform

## Quick Start

```bash
# Build
go build -o quickblog .

# Add a markdown post
echo '# Hello World\n\nMy first post!' > posts/hello.md

# Convert markdown to HTML
./quickblog convert

# Start the server
./quickblog
```

Visit `http://localhost:8080` to see your blog.

## Usage

### Start the server

```bash
# Default: port 8080, serve from ./pages
./quickblog

# Custom port
./quickblog -port 3000

# Custom pages directory
./quickblog -dir /var/www/pages

# Convert markdown + start server in one step
./quickblog -md posts
```

### Convert Markdown to HTML

```bash
# Default: posts/ -> pages/
./quickblog convert

# Custom source and output directories
./quickblog convert -md src -out dist
```

## Project Structure

```
quickBlog/
├── main.go          # Server + Markdown converter (~230 lines)
├── go.mod
├── templates/       # Embedded HTML templates
│   ├── index.html   # Homepage (article list)
│   └── page.html    # Article detail page
├── posts/           # Markdown source files (.md)
└── pages/           # Published HTML files (.html)
```

## How It Works

1. Write `.md` files in `posts/` (or `.html` files directly in `pages/`)
2. Run `./quickblog convert` to generate HTML from Markdown
3. Run `./quickblog` to start the server
4. The homepage auto-lists all `.html` files in `pages/`, sorted by date

## Cross-Compile

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o quickblog .

# macOS ARM
GOOS=darwin GOARCH=arm64 go build -o quickblog .

# Windows
GOOS=windows GOARCH=amd64 go build -o quickblog.exe .
```

## Dependencies

- [Goldmark](https://github.com/yuin/goldmark) — CommonMark compliant Markdown parser with GFM extension

## License

MIT
