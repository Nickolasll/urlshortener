package domain

type ShortRepositoryInerface interface {
	Save(short Short) error
	GetOriginalURL(slug string) (string, error)
	Ping() error
	BulkSave(shorts []Short) error
	GetShortURL(originalURL string) (string, error)
}
