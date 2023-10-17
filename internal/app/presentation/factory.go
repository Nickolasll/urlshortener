package presentation

import (
	"net/http"

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

func MuxFactory() *http.ServeMux {
	repository = initRepository()
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, mainPage)
	return mux
}

func ChiFactory() *chi.Mux {
	repository = initRepository()
	router := chi.NewRouter()
	router.Use(WithLogging)
	router.Use(gzipMiddleware)

	router.Get("/{slug}", GetHandler)
	router.Get("/ping", PingHandler)
	router.Post("/", PostHandler)
	router.Post("/api/shorten", ShortenHandler)
	router.Post("/api/shorten/batch", BatchShortenHandler)
	return router
}
