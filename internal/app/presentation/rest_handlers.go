package presentation

import (
	"context"
	"encoding/json"
	"io"

	"net/http"

	"github.com/Nickolasll/urlshortener/internal/app/config"
)

func getUserID(con context.Context) string {
	userID := con.Value(userIDKey)
	if userID != nil {
		return userID.(string)
	} else {
		return ""
	}
}

func expandHandler(res http.ResponseWriter, req *http.Request) {
	slug := req.URL.Path
	value, err := expand(slug)
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

func postHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "text/plain")
	body, _ := io.ReadAll(req.Body)
	userID := getUserID(req.Context())
	short, err := shorten(string(body), userID)
	if err != nil {
		slug, _ := getShortURLByOriginalURL(short.OriginalURL)
		res.WriteHeader(http.StatusConflict)
		res.Write([]byte(*config.SlugEndpoint + slug))
		return
	}
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(*config.SlugEndpoint + short.ShortURL))
}

func shortenHandler(res http.ResponseWriter, req *http.Request) {
	var input Input
	res.Header().Set("Content-Type", "application/json")
	body, _ := io.ReadAll(req.Body)
	err := json.Unmarshal(body, &input)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		log.Info(err)
		return
	}
	if input.URL == "" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	userID := getUserID(req.Context())
	short, err := shorten(string(body), userID)
	if err != nil {
		slug, _ := getShortURLByOriginalURL(short.OriginalURL)
		resp, _ := json.Marshal(Output{Result: *config.SlugEndpoint + slug})
		res.WriteHeader(http.StatusConflict)
		res.Write(resp)
		return
	}
	resp, _ := json.Marshal(Output{Result: *config.SlugEndpoint + short.ShortURL})
	res.WriteHeader(http.StatusCreated)
	res.Write(resp)
}

func pingHandler(res http.ResponseWriter, req *http.Request) {
	if ping() != nil {
		res.WriteHeader(http.StatusInternalServerError)
	}
}

func batchShortenHandler(res http.ResponseWriter, req *http.Request) {
	var batchInput []BatchInput
	res.Header().Set("Content-Type", "application/json")
	body, _ := io.ReadAll(req.Body)
	err := json.Unmarshal(body, &batchInput)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		log.Info(err)
		return
	}
	userID := getUserID(req.Context())
	batchOutput := batchShorten(batchInput, userID)
	resp, _ := json.Marshal(batchOutput)
	res.WriteHeader(http.StatusCreated)
	res.Write(resp)
}

func findURLsHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	userID := getUserID(req.Context())
	URLs := findURLs(userID)
	if len(URLs) == 0 {
		res.WriteHeader(http.StatusNoContent)
		return
	}
	resp, _ := json.Marshal(URLs)
	res.WriteHeader(http.StatusOK)
	res.Write(resp)
}

func deleteHandler(res http.ResponseWriter, req *http.Request) {
	var shortURLs []string
	userID := getUserID(req.Context())
	body, _ := io.ReadAll(req.Body)
	err := json.Unmarshal(body, &shortURLs)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		log.Info(err)
		return
	}
	res.WriteHeader(http.StatusAccepted)
	go bulkDelete(shortURLs, userID)
}

func getInternalStatsHandler(res http.ResponseWriter, req *http.Request) {
	stats, err := getInternalStats()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Info(err)
		return
	}
	resp, _ := json.Marshal(stats)
	res.WriteHeader(http.StatusOK)
	res.Write(resp)
}
