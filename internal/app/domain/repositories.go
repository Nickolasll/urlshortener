package domain

// IShortRepository - Интерфейс репозитория
type IShortRepository interface {
	// Save - Сохранить сокращенную ссылку
	Save(short Short) error
	// GetByShortURL - Получить объект сокращенной ссылки по слагу
	GetByShortURL(shortURL string) (Short, error)
	// Ping - Проверка подключения к бд
	Ping() error
	// BulkSave - Операция по сохранению пачки сокращенных ссылок
	BulkSave(shorts []Short) error
	// GetShortURL - Получение сокращенного слага по полному адресу ссылки
	GetShortURL(originalURL string) (string, error)
	// FindByUserID - Поиск всех сокращенных ссылок по идентификатору пользователя
	FindByUserID(userID string) ([]Short, error)
	// BulkDelete - Удаление пачки сокращенных ссылок
	BulkDelete(shortURLs []string, userID string) error
	// GetStats получение количества пользователей и сокращенных URL
	GetStats() (int, int, error)
}
