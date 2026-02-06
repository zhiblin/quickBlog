package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

//go:embed templates/*
var templateFS embed.FS

var tmpl = template.Must(template.New("").Funcs(template.FuncMap{
	"formatDate": func(t time.Time) string {
		return t.Format("2006-01-02 15:04")
	},
}).ParseFS(templateFS, "templates/*.html"))

// PageInfo holds metadata for a static page file.
type PageInfo struct {
	Name    string
	File    string
	ModTime time.Time
}

var (
	pagesDir string
	mdParser goldmark.Markdown
)

func init() {
	mdParser = goldmark.New(
		goldmark.WithExtensions(extension.GFM),
	)
}

func main() {
	// Subcommand: convert
	convertCmd := flag.NewFlagSet("convert", flag.ExitOnError)
	convertMd := convertCmd.String("md", "posts", "directory containing .md files")
	convertOut := convertCmd.String("out", "pages", "output directory for .html files")

	// Subcommand: serve (default)
	serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)
	servePort := serveCmd.Int("port", 8080, "server port")
	serveDir := serveCmd.String("dir", "pages", "directory containing .html pages")
	serveMd := serveCmd.String("md", "", "if set, convert .md files from this directory before serving")

	if len(os.Args) > 1 && os.Args[1] == "convert" {
		convertCmd.Parse(os.Args[2:])
		n, err := convertAll(*convertMd, *convertOut)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Converted %d markdown files to %s", n, *convertOut)
		return
	}

	// Default: serve
	args := os.Args[1:]
	if len(os.Args) > 1 && os.Args[1] == "serve" {
		args = os.Args[2:]
	}
	serveCmd.Parse(args)
	pagesDir = *serveDir

	// Optional: convert markdown before serving
	if *serveMd != "" {
		n, err := convertAll(*serveMd, pagesDir)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Converted %d markdown files to %s", n, pagesDir)
	}

	abs, _ := filepath.Abs(pagesDir)
	log.Printf("Serving pages from: %s", abs)

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/page/", handlePage)

	addr := fmt.Sprintf(":%d", *servePort)
	log.Printf("Server started at http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// convertAll converts all .md files in srcDir to .html files in dstDir.
func convertAll(srcDir, dstDir string) (int, error) {
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return 0, fmt.Errorf("create output dir: %w", err)
	}

	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return 0, fmt.Errorf("read source dir: %w", err)
	}

	count := 0
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}

		src := filepath.Join(srcDir, e.Name())
		mdBytes, err := os.ReadFile(src)
		if err != nil {
			log.Printf("skip %s: %v", e.Name(), err)
			continue
		}

		var buf bytes.Buffer
		if err := mdParser.Convert(mdBytes, &buf); err != nil {
			log.Printf("skip %s: convert error: %v", e.Name(), err)
			continue
		}

		outName := strings.TrimSuffix(e.Name(), ".md") + ".html"
		dst := filepath.Join(dstDir, outName)
		if err := os.WriteFile(dst, buf.Bytes(), 0644); err != nil {
			log.Printf("skip %s: write error: %v", e.Name(), err)
			continue
		}

		log.Printf("  %s -> %s", e.Name(), outName)
		count++
	}

	return count, nil
}

// handleIndex lists all .html files sorted by modification time (newest first).
func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	pages, err := listPages()
	if err != nil {
		http.Error(w, "Failed to read pages directory", http.StatusInternalServerError)
		log.Printf("listPages error: %v", err)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.ExecuteTemplate(w, "index.html", pages)
}

// handlePage serves a single static HTML file from the pages directory.
func handlePage(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/page/")
	if name == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if strings.Contains(name, "..") || strings.Contains(name, "/") {
		http.NotFound(w, r)
		return
	}

	if !strings.HasSuffix(name, ".html") {
		name += ".html"
	}

	filePath := filepath.Join(pagesDir, name)
	info, err := os.Stat(filePath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	content, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, "Failed to read page", http.StatusInternalServerError)
		return
	}

	tmpl.ExecuteTemplate(w, "page.html", struct {
		Name    string
		ModTime time.Time
		Content template.HTML
	}{
		Name:    strings.TrimSuffix(info.Name(), ".html"),
		ModTime: info.ModTime(),
		Content: template.HTML(content),
	})
}

func listPages() ([]PageInfo, error) {
	entries, err := os.ReadDir(pagesDir)
	if err != nil {
		return nil, err
	}

	var pages []PageInfo
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".html") {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		pages = append(pages, PageInfo{
			Name:    strings.TrimSuffix(e.Name(), ".html"),
			File:    e.Name(),
			ModTime: info.ModTime(),
		})
	}

	sort.Slice(pages, func(i, j int) bool {
		return pages[i].ModTime.After(pages[j].ModTime)
	})

	return pages, nil
}
