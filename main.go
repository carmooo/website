package main

import (
	"embed"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"net/http"
	"os"
	"website/handler"
)

//go:embed public
var FS embed.FS

func main() {
	if err := initEverything(); err != nil {
		log.Fatal(err)
	}

	router := chi.NewMux()

	router.Handle("/*", http.StripPrefix("/", http.FileServer(http.FS(FS))))
	router.Get("/", handler.MakeHandler(handler.HandleLandingIndex))
	router.Get("/about", handler.MakeHandler(handler.HandleAboutIndex))

	http_base_url := os.Getenv("HTTP_BASE_URL")
	slog.Info("Application running at", "url", http_base_url)
	log.Fatal(http.ListenAndServe(http_base_url, router))
}

func initEverything() error {
	return godotenv.Load()
}
