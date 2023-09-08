package presentation

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_mainPage(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		body        string
	}
	tests := []struct {
		name   string
		body   string
		method string
		want   want
	}{
		{
			name:   "success short and unshort",
			body:   "https://practicum.yandex.ru/",
			method: http.MethodPost,
			want: want{
				statusCode:  201,
				contentType: "text/plain",
				body:        "http://localhost:8080/ImePAqgP",
			},
		},
		{
			name:   "bad request",
			body:   "https://practicum.yandex.ru/",
			method: http.MethodPut,
			want: want{
				statusCode:  400,
				contentType: "",
				body:        "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyReader := strings.NewReader(tt.body)
			request := httptest.NewRequest(tt.method, "/", bodyReader)
			recorder := httptest.NewRecorder()
			handler := http.HandlerFunc(mainPage)
			handler(recorder, request)

			result := recorder.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			defer result.Body.Close()
			resBody, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			assert.Equal(t, len(tt.want.body), len(resBody))
			err = result.Body.Close()
			require.NoError(t, err)
			if tt.want.statusCode == 201 {
				bodyReader := strings.NewReader(tt.body)
				request := httptest.NewRequest(http.MethodGet, string(resBody), bodyReader)
				recorder := httptest.NewRecorder()
				handler := http.HandlerFunc(mainPage)
				handler(recorder, request)

				result := recorder.Result()
				assert.Equal(t, http.StatusTemporaryRedirect, result.StatusCode)
				assert.Equal(t, tt.body, result.Header.Get("Location"))
				err = result.Body.Close()
				require.NoError(t, err)
			}
		})
	}
}
