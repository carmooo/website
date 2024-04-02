package handler

import (
	"bytes"
	"github.com/adrg/frontmatter"
	"github.com/go-chi/chi/v5"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"website/types"
	"website/view/blog"
)

func HandleBlogIndex(w http.ResponseWriter, r *http.Request) error {
	files, err := os.ReadDir("./blog_posts")
	if err != nil {
		log.Fatal(err)
	}

	var postInfos []types.PostInfo
	for _, file := range files {
		md, err := os.ReadFile("./blog_posts/" + file.Name())
		if err != nil {
			log.Fatal(err)
		}

		var postInfo types.PostInfo
		_, err = frontmatter.MustParse(bytes.NewReader(md), &postInfo)
		if err != nil {
			log.Fatal(err)
		}

		postInfo.FileName = strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))

		postInfos = append(postInfos, postInfo)
	}

	sort.Slice(postInfos, func(i, j int) bool {
		return postInfos[i].Date.After(postInfos[j].Date)
	})

	return blog.Index(postInfos).Render(r.Context(), w)
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
	md, err := os.ReadFile("blog_posts/" + postName + ".md")
	if err != nil {
		log.Fatal(err)
	}

	// This is just to remove the frontmatter before going ahead and converting to HTML
	var emptyStruct struct{}
	md, err = frontmatter.MustParse(bytes.NewReader(md), &emptyStruct)
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

	return blog.BlogPost(html).Render(r.Context(), w)
}
