package main

import (
	"fmt"
	"net/http"

	_ "net/http/pprof"

	"github.com/Nickolasll/urlshortener/internal/app/config"
	"github.com/Nickolasll/urlshortener/internal/app/domain"
	"github.com/Nickolasll/urlshortener/internal/app/infrastructure/repositories"
	"github.com/Nickolasll/urlshortener/internal/app/presentation"
	"github.com/google/uuid"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
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
	var err error
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
	err = config.ParseFlags()
	if err != nil {
		panic(err)
	}
	go benchmark()
	mux := presentation.ChiFactory()
	if *config.EnableHTTPS {
		err = http.ListenAndServeTLS(*config.ServerEndpoint, "server.crt", "server.key", mux)
	} else {
		err = http.ListenAndServe(*config.ServerEndpoint, mux)
	}
	if err != nil {
		panic(err)
	}
}
