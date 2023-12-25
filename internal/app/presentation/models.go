package presentation

// Модель для десериализации OriginalURL
type Input struct {
	URL string `json:"url"`
}

// Модель для сериализации ShortURL
type Output struct {
	Result string `json:"result"`
}

// Модель для десериализации пачки OriginalURL
type BatchInput struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// Модель для сериализации результатов сокращения пачки OriginalURL
type BatchOutput struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// Модель для сериализации всех сокращенных пользователем ссылок
type FindURLsResult struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
