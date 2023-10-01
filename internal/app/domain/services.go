package domain

import (
	"math/rand"

	"github.com/google/uuid"

	"github.com/Nickolasll/urlshortener/internal/app/config"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GenerateSlug(size int) string {
	result := make([]byte, size)
	for i := range result {
		result[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return "/" + string(result)
}

func Shorten(url string) Short {
	return Short{
		UUID:        uuid.New().String(),
		ShortURL:    GenerateSlug(config.SlugSize),
		OriginalURL: url,
	}
}
