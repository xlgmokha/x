package main

import (
	"bytes"
	"cmp"
	_ "embed"
	"flag"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/xlgmokha/x/pkg/x"
	"github.com/xlgmokha/x/pkg/xlog"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

//go:embed index.html
var source string

const csp = "default-src 'none'; style-src 'unsafe-inline'; img-src 'self'; script-src https://cdn.jsdelivr.net 'unsafe-inline'"

var (
	page     = template.Must(template.New("index").Parse(source))
	markdown = goldmark.New(goldmark.WithExtensions(extension.GFM))
)

type node struct {
	Name, Href          string
	Dir, Open, Selected bool
	Children            []node
}

type view struct {
	Path    string
	Tree    []node
	HTML    template.HTML
	Text    string
	Mermaid bool
}

func tree(fsys fs.FS, dir, selected string) []node {
	entries := x.Must(fs.ReadDir(fsys, cmp.Or(dir, ".")))
	slices.SortFunc(entries, func(a, b fs.DirEntry) int {
		if a.IsDir() != b.IsDir() {
			if a.IsDir() {
				return -1
			}
			return 1
		}
		return strings.Compare(a.Name(), b.Name())
	})

	nodes := make([]node, 0, len(entries))
	for _, entry := range entries {
		item := path.Join(dir, entry.Name())
		child := node{
			Name:     entry.Name(),
			Href:     (&url.URL{Path: "/" + item}).String(),
			Dir:      entry.IsDir(),
			Selected: item == selected,
			Open:     strings.HasPrefix(selected, item+"/"),
		}
		if child.Dir {
			child.Children = tree(fsys, item, selected)
		}
		nodes = append(nodes, child)
	}
	return nodes
}

func binary(name string, size int64) bool {
	switch strings.ToLower(path.Ext(name)) {
	case ".png", ".jpg", ".jpeg", ".gif", ".webp", ".svg", ".ico", ".pdf":
		return true
	}
	return size > 1<<20
}

func render(name string, data []byte) view {
	if ext := strings.ToLower(path.Ext(name)); ext == ".md" || ext == ".markdown" {
		var buffer bytes.Buffer
		x.Check(markdown.Convert(data, &buffer))
		html := buffer.String()
		return view{HTML: template.HTML(html), Mermaid: strings.Contains(html, `class="language-mermaid"`)}
	}
	return view{Text: string(data)}
}

func main() {
	root := flag.String("root", ".", "directory to serve")
	addr := flag.String("addr", "127.0.0.1:8080", "address to bind to")
	flag.Parse()

	fsys := x.Must(os.OpenRoot(*root)).FS()
	logger := xlog.New(os.Stdout, xlog.Fields{})

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := cmp.Or(strings.Trim(path.Clean(r.URL.Path), "/"), ".")
		info, err := fs.Stat(fsys, name)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		var body view
		switch {
		case info.IsDir():
			if readme, err := fs.ReadFile(fsys, path.Join(name, "README.md")); err == nil {
				body = render("README.md", readme)
			}
		case binary(name, info.Size()):
			http.ServeFileFS(w, r, fsys, name)
			return
		default:
			data := x.Must(fs.ReadFile(fsys, name))
			if bytes.IndexByte(data, 0) >= 0 {
				http.ServeFileFS(w, r, fsys, name)
				return
			}
			body = render(name, data)
		}

		body.Path = name
		body.Tree = tree(fsys, "", name)
		w.Header().Set("Content-Security-Policy", csp)
		x.Check(page.Execute(w, body))
	})

	logger.Info("listening", slog.String("address", *addr), slog.String("root", *root))
	x.Check(http.ListenAndServe(*addr, x.Middleware[http.Handler](handler, xlog.HTTP(logger))))
}
