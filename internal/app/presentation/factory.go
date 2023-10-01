package presentation

import (
	"net/http"

	"github.com/Nickolasll/urlshortener/internal/app/infrastructure"
	"github.com/go-chi/chi/v5"
)

func MuxFactory() *http.ServeMux {
	infrastructure.InitRepository()
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, mainPage)
	return mux
}

func ChiFactory() *chi.Mux {
	infrastructure.InitRepository()
	router := chi.NewRouter()
	router.Get("/{slug}", WithLogging(gzipMiddleware(GetHandler)))
	router.Post("/", WithLogging(gzipMiddleware(PostHandler)))
	router.Post("/api/shorten", WithLogging(gzipMiddleware(ShortenHandler)))
	return router
}
