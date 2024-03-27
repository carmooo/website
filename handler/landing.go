package handler

import (
	"net/http"
	"website/view/landing"
)

func HandleLandingIndex(w http.ResponseWriter, r *http.Request) error {
	return landing.Index().Render(r.Context(), w)
}
