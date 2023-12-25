package domain

import (
	"math/rand"
	"time"

	"github.com/google/uuid"

	"github.com/Nickolasll/urlshortener/internal/app/config"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Генерирует случайный набор символов длиной size
func generateSlug(size int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, size)
	for i := range result {
		result[i] = letterBytes[r.Intn(len(letterBytes))]
	}
	return "/" + string(result)
}

// Создает сокращенную ссылку
func Shorten(url string, userID string) Short {
	return Short{
		UUID:        uuid.New().String(),
		ShortURL:    generateSlug(config.SlugSize),
		OriginalURL: url,
		UserID:      "4a7878fb-d657-40d9-a6e2-6c4f167ca0ce",
		Deleted:     false,
	}
}
