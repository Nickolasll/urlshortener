package domain

type IShortRepository interface {
	Save(short Short) error
	GetByShortURL(shortURL string) (Short, error)
	Ping() error
	BulkSave(shorts []Short) error
	GetShortURL(originalURL string) (string, error)
	FindByUserID(userID string) ([]Short, error)
	BulkDelete(shortURLs []string, userID string) error
}
