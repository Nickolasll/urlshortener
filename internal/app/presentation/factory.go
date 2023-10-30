package presentation

import (
	"net/http"
	"time"

	"github.com/Nickolasll/urlshortener/internal/app/config"
	"github.com/Nickolasll/urlshortener/internal/app/domain"
	"github.com/Nickolasll/urlshortener/internal/app/infrastructure/repositories"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()
var repository domain.ShortRepositoryInerface

func initRepository() domain.ShortRepositoryInerface {
	if *config.DatabaseDSN != "" {
		postgres := repositories.PostgresqlRepository{DSN: *config.DatabaseDSN, Timeout: 10 * time.Second}
		postgres.Init()
		return postgres
	} else if *config.FileStoragePath != "" {
		return repositories.FileRepository{
			Cache:    map[string]domain.Short{},
			FilePath: *config.FileStoragePath,
		}
	} else {
		return repositories.RAMRepository{
			Cache: map[string]domain.Short{},
		}
	}
}

func MuxFactory() *http.ServeMux {
	repository = initRepository()
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, mainPage)
	return mux
}

func ChiFactory() *chi.Mux {
	repository = initRepository()
	router := chi.NewRouter()
	router.Use(logging)
	router.Use(compress)

	router.Get("/{slug}", ExpandHandler)
	router.Get("/ping", PingHandler)
	router.Post("/", setCookie(PostHandler))
	router.Post("/api/shorten", setCookie(ShortenHandler))
	router.Post("/api/shorten/batch", setCookie(BatchShortenHandler))
	router.Get("/api/user/urls", authorize(FindURLs))
	router.Delete("/api/user/urls", authorize(Delete))
	return router
}
