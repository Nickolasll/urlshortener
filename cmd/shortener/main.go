package main

import (
	"flag"
	"net/http"

	"github.com/Nickolasll/urlshortener/internal/app/config"
	"github.com/Nickolasll/urlshortener/internal/app/presentation"
)

func main() {
	flag.Parse()
	mux := presentation.ChiFactory()
	err := http.ListenAndServe(*config.ServerEndpoint, mux)
	if err != nil {
		panic(err)
	}
}
