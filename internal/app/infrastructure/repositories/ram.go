package repositories

import (
	"github.com/Nickolasll/urlshortener/internal/app/domain"
)

type RAMRepository struct {
	Cache     map[string]string
	ListCache []domain.Short
}

func (r RAMRepository) Save(short domain.Short) error {
	r.Cache[short.ShortURL] = short.OriginalURL
	r.ListCache = append(r.ListCache, short)
	return nil
}

func (r RAMRepository) GetOriginalURL(slug string) (string, error) {
	value := r.Cache[slug]
	return value, nil
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

func (r RAMRepository) GetShortURL(originalURL string) (string, error) {
	key, _ := mapkey(r.Cache, originalURL)
	return key, nil
}

func (r RAMRepository) FindByUserID(userID string) ([]domain.Short, error) {
	shorts := []domain.Short{}
	for _, short := range r.ListCache {
		if short.UserID == userID {
			shorts = append(shorts, short)
		}
	}
	return shorts, nil
}
