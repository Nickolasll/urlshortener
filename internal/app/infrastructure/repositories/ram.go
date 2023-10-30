package repositories

import (
	"slices"

	"github.com/Nickolasll/urlshortener/internal/app/domain"
)

type RAMRepository struct {
	Cache map[string]domain.Short
}

func (r RAMRepository) Save(short domain.Short) error {
	r.Cache[short.ShortURL] = short
	return nil
}

func (r RAMRepository) GetByShortURL(shortURL string) (domain.Short, error) {
	short := r.Cache[shortURL]
	return short, nil
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
	key, _ := originalURLKeyMap(r.Cache, originalURL)
	return key, nil
}

func (r RAMRepository) FindByUserID(userID string) ([]domain.Short, error) {
	shorts := []domain.Short{}
	for _, short := range r.Cache {
		if short.UserID == userID {
			shorts = append(shorts, short)
		}
	}
	return shorts, nil
}

func (r RAMRepository) BulkDelete(shortURLs []string, userID string) error {
	for key, short := range r.Cache {
		if slices.Contains(shortURLs, short.ShortURL) && short.UserID == userID {
			delete(r.Cache, key)
		}
	}
	return nil
}
