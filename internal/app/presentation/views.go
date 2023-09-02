package presentation

import (
	"io"
	"math/rand"
	"net/http"
)

var urlShortenerMap map[string]string = map[string]string{}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(size int) string {
	result := make([]byte, size)
	for i := range result {
		result[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(result)
}

func mainPage(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		body, _ := io.ReadAll(req.Body)
		slug := "/" + randStringBytes(8)
		urlShortenerMap[slug] = string(body)
		res.Header().Set("content-type", "text/plain")
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte("http://" + req.Host + slug))
		return
	} else if req.Method == http.MethodGet {
		value, ok := urlShortenerMap[req.URL.Path]
		if ok {
			res.Header().Add("Location", value)
			res.WriteHeader(http.StatusTemporaryRedirect)
			return
		}
	}
	res.WriteHeader(http.StatusBadRequest)
}
