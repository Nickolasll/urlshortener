package presentation

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/Nickolasll/urlshortener/internal/app/config"
	"github.com/Nickolasll/urlshortener/internal/app/domain"
)

func getUserID(con context.Context) string {
	userID := con.Value(userIDKey)
	if userID != nil {
		return userID.(string)
	} else {
		return ""
	}
}

func ExpandHandler(res http.ResponseWriter, req *http.Request) {
	slug := req.URL.Path
	value, err := repository.GetByShortURL(slug)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	if value.Deleted {
		res.WriteHeader(http.StatusGone)
		return
	}

	res.Header().Add("Location", value.OriginalURL)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func PostHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "text/plain")
	body, _ := io.ReadAll(req.Body)
	userID := getUserID(req.Context())
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
	userID := getUserID(req.Context())
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
	var shorts []domain.Short
	res.Header().Set("Content-Type", "application/json")
	body, _ := io.ReadAll(req.Body)
	json.Unmarshal(body, &batchInput)
	batchOutput := []BatchOutput{}
	userID := getUserID(req.Context())
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
	userID := getUserID(req.Context())
	shorts, _ := repository.FindByUserID(userID)
	if len(shorts) == 0 {
		res.WriteHeader(http.StatusNoContent)
		return
	}
	for _, short := range shorts {
		result := FindURLsResult{
			ShortURL:    *config.SlugEndpoint + short.ShortURL,
			OriginalURL: short.OriginalURL,
		}
		URLResults = append(URLResults, result)
	}
	resp, _ := json.Marshal(URLResults)
	res.WriteHeader(http.StatusOK)
	res.Write(resp)
}

func Delete(res http.ResponseWriter, req *http.Request) {
	var shortURLs []string
	userID := getUserID(req.Context())
	body, _ := io.ReadAll(req.Body)
	json.Unmarshal(body, &shortURLs)
	for index, shortURL := range shortURLs {
		shortURL = "/" + shortURL
		shortURLs[index] = shortURL
	}
	res.WriteHeader(http.StatusAccepted)
	go repository.BulkDelete(shortURLs, userID)
}
