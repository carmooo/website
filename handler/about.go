package handler

import (
	"net/http"
	"website/view/about"
)

func HandleAboutIndex(w http.ResponseWriter, r *http.Request) error {
	return about.Index().Render(r.Context(), w)
}
