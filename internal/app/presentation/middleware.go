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

// WithLogging добавляет дополнительный код для регистрации сведений о запросе
// и возвращает новый http.Handler.
func WithLogging(handlerFn http.HandlerFunc) http.HandlerFunc {
	logFn := func(res http.ResponseWriter, req *http.Request) {

		recorder := &ResponseRecorder{
			ResponseWriter: res,
			Status:         200,
			ContentLength:  0,
		}
		// функция Now() возвращает текущее время
		start := time.Now()

		// эндпоинт /ping
		uri := req.RequestURI
		// метод запроса
		method := req.Method

		// точка, где выполняется хендлер pingHandler
		handlerFn.ServeHTTP(recorder, req) // обслуживание оригинального запроса

		// Since возвращает разницу во времени между start
		// и моментом вызова Since. Таким образом можно посчитать
		// время выполнения запроса.
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
	// возвращаем функционально расширенный хендлер
	return http.HandlerFunc(logFn)
}
