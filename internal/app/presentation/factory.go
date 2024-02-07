package presentation

import (
	"net/http"
	"time"

	"github.com/Nickolasll/urlshortener/internal/app/config"
	"github.com/Nickolasll/urlshortener/internal/app/domain"
	"github.com/Nickolasll/urlshortener/internal/app/infrastructure/repositories"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()
var repository domain.IShortRepository

func initRepository() domain.IShortRepository {
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
			OriginalToShorts: map[string]domain.Short{},
			Cache:            map[string]string{},
		}
	}
}

// MuxFactory - Фабрика приложения с использованием роутера mux
func MuxFactory() *http.ServeMux {
	repository = initRepository()
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, mainPage)
	return mux
}

// ChiFactory - Фабрика приложения с использованием роутера chi
func ChiFactory() *chi.Mux {
	repository = initRepository()
	router := chi.NewRouter()
	router.Use(logging)
	router.Use(compress)

	cookieSubRouter := chi.NewRouter()
	cookieSubRouter.Use(setCookie)
	cookieSubRouter.Post("/", postHandler)
	cookieSubRouter.Post("/api/shorten", shortenHandler)
	cookieSubRouter.Post("/api/shorten/batch", batchShortenHandler)

	router.Get("/{slug}", expandHandler)
	router.Get("/ping", pingHandler)
	router.Get("/api/internal/stats", trusted(getInternalStatsHandler))
	router.Get("/api/user/urls", authorize(findURLsHandler))
	router.Delete("/api/user/urls", authorize(deleteHandler))
	router.Mount("/", cookieSubRouter)
	router.Mount("/debug", middleware.Profiler())

	return router
}
