package domain

import (
	"math/rand"
	"time"

	"github.com/google/uuid"

	"github.com/Nickolasll/urlshortener/internal/app/config"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GenerateSlug(size int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, size)
	for i := range result {
		result[i] = letterBytes[r.Intn(len(letterBytes))]
	}
	return "/" + string(result)
}

func Shorten(url string, userID string) Short {
	return Short{
		UUID:        uuid.New().String(),
		ShortURL:    GenerateSlug(config.SlugSize),
		OriginalURL: url,
		UserID:      userID,
		Deleted:     false,
	}
}
