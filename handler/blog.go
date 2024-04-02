package handler

import (
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
