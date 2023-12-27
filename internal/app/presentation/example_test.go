package presentation_test

import (
	"net/http"

	"github.com/Nickolasll/urlshortener/internal/app/config"
	"github.com/Nickolasll/urlshortener/internal/app/presentation"
)

func Example() {
	config.ParseFlags()
	mux := presentation.ChiFactory()
	err := http.ListenAndServe(*config.ServerEndpoint, mux)
	if err != nil {
		// Не удалось запустить сервер
	}
}
