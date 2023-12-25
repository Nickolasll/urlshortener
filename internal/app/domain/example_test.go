package domain_test

import "github.com/Nickolasll/urlshortener/internal/app/domain"

func ExampleShorten() {
	URL := "https://yandex.ru"
	UserID := "4a7878fb-d657-40d9-a6e2-6c4f167ca0ce"
	short := domain.Shorten(URL, UserID)
	if short.OriginalURL == URL {
		// Здесь записывается url до сокращения
	}
	if short.ShortURL != "" {
		// Здесь генерируется значение сокращенного url
	}
}
