package infrastructure

import (
	"testing"

	"github.com/Nickolasll/urlshortener/internal/app/domain"
	"github.com/stretchr/testify/assert"
)

func TestRAMRepository(t *testing.T) {
	type args struct {
		short domain.Short
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "save and get",
			args: args{
				short: domain.Short{
					UUID:        "1",
					OriginalURL: "http://test.com",
					ShortURL:    "/ABCDE",
				},
			},
			want: "http://test.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := RAMRepository{
				urlShortenerMap: map[string]string{},
			}
			repository.Save(tt.args.short)
			got, ok := repository.Get(tt.args.short.ShortURL)
			assert.Equal(t, ok, true)
			assert.Equal(t, got, tt.want)
		})
	}
}
