package presentation

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Nickolasll/urlshortener/internal/app/auth"
	"github.com/Nickolasll/urlshortener/internal/app/config"
	"github.com/sirupsen/logrus"
)

type key int

type ResponseRecorder struct {
	http.ResponseWriter
	Status        int
	ContentLength int
}

const userIDKey key = 0

func (r *ResponseRecorder) Write(buf []byte) (int, error) {
	r.ContentLength = len(buf)
	return r.ResponseWriter.Write(buf)
}

func (r *ResponseRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func logging(handler http.Handler) http.Handler {
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

func compress(handler http.Handler) http.Handler {
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

func setCookie(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, reader *http.Request) {
		var cookie *http.Cookie
		cookie, err := reader.Cookie("authorization")
		if err != nil || !auth.IsValid(cookie.Value) {
			token, err := auth.IssueToken()
			if err != nil {
				writer.WriteHeader(http.StatusInternalServerError)
				log.Info(err)
				return
			}
			cookie = &http.Cookie{
				Name:   "authorization",
				Value:  token,
				MaxAge: config.TokenExp,
				Path:   "/",
			}
			http.SetCookie(writer, cookie)
		}
		UserID := auth.GetUserID(cookie.Value)
		ctx := context.WithValue(reader.Context(), userIDKey, UserID)
		handler.ServeHTTP(writer, reader.WithContext(ctx))
	})
}

func authorize(handlerFn http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, reader *http.Request) {
		cookie, err := reader.Cookie("authorization")
		if err != nil || !auth.IsValid(cookie.Value) {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		UserID := auth.GetUserID(cookie.Value)
		ctx := context.WithValue(reader.Context(), userIDKey, UserID)
		handlerFn.ServeHTTP(writer, reader.WithContext(ctx))
	}
}
