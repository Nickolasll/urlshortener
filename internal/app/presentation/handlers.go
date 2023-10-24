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

type FindURLsResult struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func GetHandler(res http.ResponseWriter, req *http.Request) {
	slug := req.URL.Path
	value, _ := repository.GetOriginalURL(slug)
	if value != "" {
		res.Header().Add("Location", value)
		res.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func PostHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "text/plain")
	body, _ := io.ReadAll(req.Body)
	userID := req.Context().Value(userIDKey).(string)
	short := domain.Shorten(string(body), userID)
	err := repository.Save(short)
	if err != nil {
		slug, _ := repository.GetShortURL(short.OriginalURL)
		res.WriteHeader(http.StatusConflict)
		res.Write([]byte(*config.SlugEndpoint + slug))
		return
	}
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(*config.SlugEndpoint + short.ShortURL))
}

func ShortenHandler(res http.ResponseWriter, req *http.Request) {
	var input Input
	res.Header().Set("Content-Type", "application/json")
	body, _ := io.ReadAll(req.Body)
	json.Unmarshal(body, &input)
	if input.URL == "" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	userID := req.Context().Value(userIDKey).(string)
	short := domain.Shorten(input.URL, userID)
	err := repository.Save(short)
	if err != nil {
		slug, _ := repository.GetShortURL(short.OriginalURL)
		resp, _ := json.Marshal(Output{Result: *config.SlugEndpoint + slug})
		res.WriteHeader(http.StatusConflict)
		res.Write(resp)
		return
	}
	resp, _ := json.Marshal(Output{Result: *config.SlugEndpoint + short.ShortURL})
	res.WriteHeader(http.StatusCreated)
	res.Write(resp)
}

func PingHandler(res http.ResponseWriter, req *http.Request) {
	if repository.Ping() != nil {
		res.WriteHeader(http.StatusInternalServerError)
	}
}

func BatchShortenHandler(res http.ResponseWriter, req *http.Request) {
	var batchInput []BatchInput
	res.Header().Set("Content-Type", "application/json")
	body, _ := io.ReadAll(req.Body)
	json.Unmarshal(body, &batchInput)
	shorts := []domain.Short{}
	batchOutput := []BatchOutput{}
	userID := req.Context().Value(userIDKey).(string)
	for _, batch := range batchInput {
		short := domain.Shorten(batch.OriginalURL, userID)
		shorts = append(shorts, short)
		output := BatchOutput{
			CorrelationID: batch.CorrelationID,
			ShortURL:      *config.SlugEndpoint + short.ShortURL,
		}
		batchOutput = append(batchOutput, output)
	}
	repository.BulkSave(shorts)
	resp, _ := json.Marshal(batchOutput)
	res.WriteHeader(http.StatusCreated)
	res.Write(resp)
}

func FindURLs(res http.ResponseWriter, req *http.Request) {
	var URLResults []FindURLsResult
	res.Header().Set("Content-Type", "application/json")
	userID := req.Context().Value(userIDKey).(string)
	shorts, _ := repository.FindByUserID(userID)
	if len(shorts) == 0 {
		res.WriteHeader(http.StatusNoContent)
		return
	}
	for _, short := range shorts {
		result := FindURLsResult{
			ShortURL:    short.ShortURL,
			OriginalURL: short.OriginalURL,
		}
		URLResults = append(URLResults, result)
	}
	resp, _ := json.Marshal(URLResults)
	res.WriteHeader(http.StatusOK)
	res.Write(resp)
}
