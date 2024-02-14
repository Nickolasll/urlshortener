package presentation

// Input - Модель для десериализации OriginalURL
type Input struct {
	URL string `json:"url"`
}

// Output - Модель для сериализации ShortURL
type Output struct {
	Result string `json:"result"`
}

// BatchInput - Модель для десериализации пачки OriginalURL
type BatchInput struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// BatchOutput - Модель для сериализации результатов сокращения пачки OriginalURL
type BatchOutput struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// FindURLsResult - Модель для сериализации всех сокращенных пользователем ссылок
type FindURLsResult struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// GetInternalStatsResult - Модель для сериализации количества пользователей и сокращенных URL
type GetInternalStatsResult struct {
	URLs  int `json:"urls"`
	Users int `json:"users"`
}
