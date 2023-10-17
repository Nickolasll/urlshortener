package presentation

import (
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
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

func WithLogging(handler http.Handler) http.Handler {
	logFn := func(res http.ResponseWriter, req *http.Request) {

		recorder := &ResponseRecorder{
			ResponseWriter: res,
			Status:         200,
			ContentLength:  0,
		}
		start := time.Now()
		uri := req.RequestURI
		method := req.Method

		handler.ServeHTTP(recorder, req)

		duration := time.Since(start)

		log.WithFields(logrus.Fields{
			"uri":      uri,
			"method":   method,
			"duration": duration,
		}).Info("Request info")

		log.WithFields(logrus.Fields{
			"status":         recorder.Status,
			"content length": recorder.ContentLength,
		}).Info("Response info")

	}
	return http.HandlerFunc(logFn)
}

func gzipMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, reader *http.Request) {
		originalWriter := writer

		acceptEncoding := reader.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			compressWriter := newCompressWriter(writer)
			originalWriter = compressWriter
			defer compressWriter.Close()
		}

		contentEncoding := reader.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := newCompressReader(reader.Body)
			if err != nil {
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}
			reader.Body = cr
			defer cr.Close()
		}

		handler.ServeHTTP(originalWriter, reader)
	})
}
