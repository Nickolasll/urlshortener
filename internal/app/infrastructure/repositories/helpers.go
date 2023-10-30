package repositories

import "github.com/Nickolasll/urlshortener/internal/app/domain"

func originalURLKeyMap(m map[string]domain.Short, value string) (key string, ok bool) {
	for k, v := range m {
		if v.OriginalURL == value {
			key = k
			ok = true
			return
		}
	}
	return
}
