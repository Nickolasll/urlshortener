package presentation

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Nickolasll/urlshortener/internal/app/config"
	"github.com/Nickolasll/urlshortener/internal/app/domain"
)

func GetHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		slug := req.URL.Path
		value, ok := Repository.Get(slug)
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
		short := domain.Shorten(string(body))
		Repository.Save(short)
		res.Header().Set("content-type", "text/plain")
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(*config.SlugEndpoint + short.ShortURL))
		return
	}
}

type Input struct {
	URL string `json:"url"`
}

type Output struct {
	Result string `json:"result"`
}

func ShortenHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		res.Header().Set("Content-Type", "application/json")
		body, _ := io.ReadAll(req.Body)
		var input Input
		json.Unmarshal(body, &input)
		if input.URL == "" {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		short := domain.Shorten(input.URL)
		Repository.Save(short)
		resp, _ := json.Marshal(Output{Result: *config.SlugEndpoint + short.ShortURL})
		res.WriteHeader(http.StatusCreated)
		res.Write(resp)
		return
	}
}
