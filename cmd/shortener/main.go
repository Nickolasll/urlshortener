package main

import (
	"net/http"

	_ "net/http/pprof"

	"github.com/Nickolasll/urlshortener/internal/app/config"
	"github.com/Nickolasll/urlshortener/internal/app/domain"
	"github.com/Nickolasll/urlshortener/internal/app/infrastructure/repositories"
	"github.com/Nickolasll/urlshortener/internal/app/presentation"
	"github.com/google/uuid"
)

func benchmark() {
	for {
		URL := "www.testlongurl.com"
		userID := uuid.New().String()
		repository := repositories.RAMRepository{
			OriginalToShorts: map[string]domain.Short{},
			Cache:            map[string]string{},
		}
		for i := 0; i < 10000000; i++ {
			short := domain.Shorten(URL, userID)
			repository.Save(short)
		}
	}
}

func main() {
	config.ParseFlags()
	go benchmark()
	mux := presentation.ChiFactory()
	err := http.ListenAndServe(*config.ServerEndpoint, mux)
	if err != nil {
		panic(err)
	}
}
