package repositories_test

import (
	"github.com/Nickolasll/urlshortener/internal/app/domain"
	"github.com/Nickolasll/urlshortener/internal/app/infrastructure/repositories"
)

func Example() {
	RAMRepo := repositories.RAMRepository{
		OriginalToShorts: map[string]domain.Short{},
		Cache:            map[string]string{},
	}
	short := domain.Short{
		UUID:        "d9e02c96-5b06-43c7-8746-92e94239ae55",
		ShortURL:    "ABCDE",
		OriginalURL: "https://yandex.ru",
		UserID:      "4a7878fb-d657-40d9-a6e2-6c4f167ca0ce",
		Deleted:     false,
	}
	RAMRepo.Save(short)
	shortURL, err := RAMRepo.GetShortURL("https://yandex.ru")
	if err != nil {
		// Запись не найдена
	}
	if shortURL != "" {
		// Получаем запись
	}
}
