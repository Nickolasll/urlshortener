package presentation

import (
	"net/http"

	"github.com/Nickolasll/urlshortener/internal/app/config"
	"github.com/Nickolasll/urlshortener/internal/app/domain"
	"github.com/Nickolasll/urlshortener/internal/app/infrastructure/repositories"
	"github.com/go-chi/chi/v5"
)

func initRepository() domain.ShortRepositoryInerface {
	if *config.DatabaseDSN != "" {
		postgres := repositories.PostgresqlRepository{DSN: *config.DatabaseDSN}
		postgres.Init()
		return postgres
	} else if *config.FileStoragePath != "" {
		return repositories.FileRepository{
			Cache:    map[string]string{},
			FilePath: *config.FileStoragePath,
		}
	} else {
		return repositories.RAMRepository{
			Cache: map[string]string{},
		}
	}
}

var repository domain.ShortRepositoryInerface

func MuxFactory() *http.ServeMux {
	repository = initRepository()
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, mainPage)
	return mux
}

func ChiFactory() *chi.Mux {
	repository = initRepository()
	router := chi.NewRouter()
	router.Get("/{slug}", WithLogging(gzipMiddleware(GetHandler)))
	router.Get("/ping", WithLogging(gzipMiddleware(PingHandler)))
	router.Post("/", WithLogging(gzipMiddleware(PostHandler)))
	router.Post("/api/shorten", WithLogging(gzipMiddleware(ShortenHandler)))
	router.Post("/api/shorten/batch", WithLogging(gzipMiddleware(BatchShortenHandler)))
	return router
}
