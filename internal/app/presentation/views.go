package presentation

import (
	"io"
	"math/rand"
	"net/http"
)

var urlShortenerMap map[string][]byte = map[string][]byte{}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(size int) string {
	result := make([]byte, size)
	for i := range result {
		result[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(result)
}

func mainPage(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "text/plain")
	if req.Method == http.MethodPost {
		res.WriteHeader(http.StatusCreated)
		body, _ := io.ReadAll(req.Body)
		shorted_url := "/" + randStringBytes(8)
		urlShortenerMap[shorted_url] = body
		res.Write([]byte("http://" + req.Host + shorted_url))
		return
	} else if req.Method == http.MethodGet {
		value, ok := urlShortenerMap[req.URL.Path]
		if ok {
			res.WriteHeader(http.StatusTemporaryRedirect)
			res.Write(value)
			return
		}
	}
	res.WriteHeader(http.StatusBadRequest)
}
