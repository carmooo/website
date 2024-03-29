package handler

import (
	"net/http"
	"website/view/profile"
)

func HandleProfileIndex(w http.ResponseWriter, r *http.Request) error {
	return profile.Index().Render(r.Context(), w)
}
