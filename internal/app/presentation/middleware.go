package presentation

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type ResponseRecorder struct {
	http.ResponseWriter
	Status        int
	ContentLength int
}

func (r *ResponseRecorder) Write(buf []byte) (int, error) {
	r.ContentLength = len(buf)
	return r.ResponseWriter.Write(buf)
}

func (r *ResponseRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func WithLogging(handlerFn http.HandlerFunc) http.HandlerFunc {
	logFn := func(res http.ResponseWriter, req *http.Request) {

		recorder := &ResponseRecorder{
			ResponseWriter: res,
			Status:         200,
			ContentLength:  0,
		}
		start := time.Now()
		uri := req.RequestURI
		method := req.Method

		handlerFn.ServeHTTP(recorder, req)

		duration := time.Since(start)

		log.WithFields(log.Fields{
			"uri":      uri,
			"method":   method,
			"duration": duration,
		}).Info("Request info")

		log.WithFields(log.Fields{
			"status":         recorder.Status,
			"content length": recorder.ContentLength,
		}).Info("Response info")

	}
	return http.HandlerFunc(logFn)
}
