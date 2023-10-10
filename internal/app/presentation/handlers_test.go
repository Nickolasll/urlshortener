package presentation

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ShortenHandler(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name string
		body []byte
		want want
	}{
		{
			name: "success short",
			body: []byte(`{"url": "https://practicum.yandex.ru"}`),
			want: want{
				statusCode:  http.StatusCreated,
				contentType: "application/json",
			},
		},
		{
			name: "not a json",
			body: []byte(`not a json`),
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "application/json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository = initRepository()
			bodyReader := bytes.NewReader(tt.body)
			request := httptest.NewRequest(http.MethodPost, "/api/shorten", bodyReader)
			recorder := httptest.NewRecorder()
			ShortenHandler(recorder, request)

			result := recorder.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			defer result.Body.Close()

			if result.StatusCode == http.StatusOK {
				resBody, err := io.ReadAll(result.Body)
				require.NoError(t, err)
				var output Output
				json.Unmarshal(resBody, &output)
				assert.NotEmpty(t, output.Result)
				err = result.Body.Close()
				require.NoError(t, err)
			}
		})
	}
}
