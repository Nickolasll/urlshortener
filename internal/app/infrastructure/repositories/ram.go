package repositories

import (
	"github.com/Nickolasll/urlshortener/internal/app/domain"
)

type RAMRepository struct {
	Cache map[string]string
}

func (r RAMRepository) Save(short domain.Short) error {
	r.Cache[short.ShortURL] = short.OriginalURL
	return nil
}

func (r RAMRepository) Get(slug string) (string, bool, error) {
	value, ok := r.Cache[slug]
	return value, ok, nil
}

func (r RAMRepository) Ping() error {
	return nil
}

func (r RAMRepository) BulkSave(shorts []domain.Short) error {
	for _, short := range shorts {
		r.Save(short)
	}
	return nil
}
