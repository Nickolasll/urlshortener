package presentation

import (
	"io"
	"net/http"

	"github.com/Nickolasll/urlshortener/internal/app/config"
	"github.com/Nickolasll/urlshortener/internal/app/domain"
	"github.com/Nickolasll/urlshortener/internal/app/infrastructure"
)

func GetHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		slug := req.URL.Path
		value, ok := infrastructure.RAMRepository.Get(slug)
		if ok {
			res.Header().Add("Location", value)
			res.WriteHeader(http.StatusTemporaryRedirect)
			return
		}
	}
}

func PostHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		body, _ := io.ReadAll(req.Body)
		slug := domain.GenerateSlug(8)
		infrastructure.RAMRepository.Save(string(body), slug)
		res.Header().Set("content-type", "text/plain")
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(*config.SlugEndpoint + slug))
		return
	}
}
