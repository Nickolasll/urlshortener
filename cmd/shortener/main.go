package main

import (
	"net/http"

	"github.com/Nickolasll/urlshortener/internal/app/presentation"
)

func main() {
	mux := presentation.MuxFactory()
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
