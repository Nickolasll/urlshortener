package presentation

type Input struct {
	URL string `json:"url"`
}

type Output struct {
	Result string `json:"result"`
}

type BatchInput struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchOutput struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type FindURLsResult struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
