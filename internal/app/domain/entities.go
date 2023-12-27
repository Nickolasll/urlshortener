package domain

// Short - Сокращенная ссылка
type Short struct {
	// UUID - Уникальный идентификатор сокращенной ссылки
	UUID string `json:"uuid"`
	// ShortURL - Сгенерированный слаг сокращенной ссылки
	ShortURL string `json:"short_url"`
	// OriginalURL - Оригинальная ссылка
	OriginalURL string `json:"original_url"`
	// UserID - Идентификатор пользователя
	UserID string `json:"user_id"`
	// Deleted - Флаг удаленности
	Deleted bool `json:"deleted"`
}
