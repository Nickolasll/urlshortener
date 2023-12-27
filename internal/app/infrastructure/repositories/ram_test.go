package repositories

import (
	"strconv"
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
				OriginalToShorts: map[string]domain.Short{},
				Cache:            map[string]string{},
			}
			repository.Save(tt.args.short)
			got, _ := repository.GetByShortURL(tt.args.short.ShortURL)
			assert.Equal(t, got.OriginalURL, tt.want)
		})
	}
}

func BenchmarkRAMRepositorySave(b *testing.B) {
	short := domain.Short{
		UUID:        "1",
		OriginalURL: "http://test.com",
		ShortURL:    "",
	}
	repository := RAMRepository{
		OriginalToShorts: map[string]domain.Short{},
		Cache:            map[string]string{},
	}
	for i := 0; i < b.N; i++ {
		short.ShortURL = strconv.Itoa(i)
		repository.Save(short)
	}
}
