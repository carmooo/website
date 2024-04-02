package handler

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"website/view/blog"
)

func HandleBlogIndex(w http.ResponseWriter, r *http.Request) error {
	files, err := os.ReadDir("./blog_posts")
	if err != nil {
		log.Fatal(err)
	}

	sort.Slice(files, func(i, j int) bool {
		fileI, err := files[i].Info()
		if err != nil {
			log.Fatal(err)
		}
		fileJ, err := files[j].Info()
		if err != nil {
			log.Fatal(err)
		}
		return fileI.ModTime().After(fileJ.ModTime())
	})

	return blog.Index(files).Render(r.Context(), w)
}

func renderTable(w io.Writer, p *ast.Table, entering bool) {
	if entering {
		io.WriteString(w, "<div class=\"overflow-x-auto\">\n  <table class=\"table\">")
	} else {
		io.WriteString(w, "</table>\n</div>")
	}
}

func tableRenderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	if table, ok := node.(*ast.Table); ok {
		renderTable(w, table, entering)
		return ast.GoToNext, true
	}
	return ast.GoToNext, false
}

func HandleBlogPost(w http.ResponseWriter, r *http.Request) error {
	postName := chi.URLParam(r, "postName")
	f, err := os.Open("blog_posts/" + postName + ".md")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	md, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{
		Flags:          htmlFlags,
		RenderNodeHook: tableRenderHook,
	}
	renderer := html.NewRenderer(opts)
	html := markdown.Render(doc, renderer)

	fmt.Printf("--- Markdown:\n%s\n\n--- HTML:\n%s\n", md, html)

	return blog.BlogPost(html).Render(r.Context(), w)
}
