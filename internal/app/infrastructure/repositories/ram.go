package repositories

import (
	"github.com/Nickolasll/urlshortener/internal/app/domain"
)

// RAMRepository - Имплементация репозитория в оперативной памяти
type RAMRepository struct {
	// OriginalToShorts - маппинг оригинальной ссылки к объекту сокращенной ссылки
	// Cache - маппинг сокращенной ссылки к оригинальной ссылке
	OriginalToShorts map[string]domain.Short
	Cache            map[string]string
}

// Save - Сохранить сокращенную ссылку
func (r RAMRepository) Save(short domain.Short) error {
	r.OriginalToShorts[short.OriginalURL] = short
	r.Cache[short.ShortURL] = short.OriginalURL
	return nil
}

// GetByShortURL - Получить объект сокращенной ссылки по значению
func (r RAMRepository) GetByShortURL(shortURL string) (domain.Short, error) {
	originalURL := r.Cache[shortURL]
	short := r.OriginalToShorts[originalURL]
	return short, nil
}

// Ping - Проверка работоспособности
func (r RAMRepository) Ping() error {
	return nil
}

// BulkSave - Сохранить пачку сокращенных ссылок
func (r RAMRepository) BulkSave(shorts []domain.Short) error {
	for _, short := range shorts {
		r.Save(short)
	}
	return nil
}

// GetShortURL - Получить сокращенную ссылку по несокращенному значению
func (r RAMRepository) GetShortURL(originalURL string) (string, error) {
	short := r.OriginalToShorts[originalURL]
	return short.OriginalURL, nil
}

// FindByUserID - Получить список сокращенных ссылок по идентификатору пользователя
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

// BulkDelete - Удалить пачку сокращенных ссылок
func (r RAMRepository) BulkDelete(shortURLs []string, userID string) error {
	for key, short := range r.OriginalToShorts {
		if contains(shortURLs, short.ShortURL) && short.UserID == userID {
			short.Deleted = true
			r.OriginalToShorts[key] = short
		}
	}
	return nil
}
