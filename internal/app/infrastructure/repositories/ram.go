package repositories

import (
	"github.com/Nickolasll/urlshortener/internal/app/domain"
)

type RAMRepository struct {
	OriginalToShorts map[string]domain.Short
	Cache            map[string]string
}

func (r RAMRepository) Save(short domain.Short) error {
	r.OriginalToShorts[short.OriginalURL] = short
	r.Cache[short.ShortURL] = short.OriginalURL
	return nil
}

func (r RAMRepository) GetByShortURL(shortURL string) (domain.Short, error) {
	originalURL := r.Cache[shortURL]
	short := r.OriginalToShorts[originalURL]
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
	short := r.OriginalToShorts[originalURL]
	return short.OriginalURL, nil
}

func (r RAMRepository) FindByUserID(userID string) ([]domain.Short, error) {
	shorts := []domain.Short{}
	for _, short := range r.OriginalToShorts {
		if short.UserID == userID {
			shorts = append(shorts, short)
		}
	}
	return shorts, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (r RAMRepository) BulkDelete(shortURLs []string, userID string) error {
	for key, short := range r.OriginalToShorts {
		if contains(shortURLs, short.ShortURL) && short.UserID == userID {
			short.Deleted = true
			r.OriginalToShorts[key] = short
		}
	}
	return nil
}
