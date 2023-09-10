package main

import (
	"net/http"

	"github.com/Nickolasll/urlshortener/internal/app/config"
	"github.com/Nickolasll/urlshortener/internal/app/presentation"
)

func main() {
	config.ParseFlags()
	mux := presentation.ChiFactory()
	err := http.ListenAndServe(*config.ServerEndpoint, mux)
	if err != nil {
		panic(err)
	}
}
