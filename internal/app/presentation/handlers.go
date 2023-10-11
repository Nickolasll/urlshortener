package presentation

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Nickolasll/urlshortener/internal/app/config"
	"github.com/Nickolasll/urlshortener/internal/app/domain"
)

type Input struct {
	URL string `json:"url"`
}

type Output struct {
	Result string `json:"result"`
}

type BatchInput struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchOutput struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func GetHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		slug := req.URL.Path
		value, ok, _ := repository.Get(slug)
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
		repository.Save(short)
		res.Header().Set("content-type", "text/plain")
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(*config.SlugEndpoint + short.ShortURL))
		return
	}
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
		repository.Save(short)
		resp, _ := json.Marshal(Output{Result: *config.SlugEndpoint + short.ShortURL})
		res.WriteHeader(http.StatusCreated)
		res.Write(resp)
		return
	}
}

func PingHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		if repository.Ping() != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func BatchShortenHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		res.Header().Set("Content-Type", "application/json")
		body, _ := io.ReadAll(req.Body)
		var batchInput []BatchInput
		json.Unmarshal(body, &batchInput)
		shorts := []domain.Short{}
		batchOutput := []BatchOutput{}
		for _, batch := range batchInput {
			short := domain.Shorten(batch.OriginalURL)
			shorts = append(shorts, short)
			output := BatchOutput{
				CorrelationID: batch.CorrelationID,
				ShortURL:      short.ShortURL,
			}
			batchOutput = append(batchOutput, output)
		}
		repository.BulkSave(shorts)
		resp, _ := json.Marshal(batchOutput)
		res.WriteHeader(http.StatusCreated)
		res.Write(resp)
		return
	}
}
