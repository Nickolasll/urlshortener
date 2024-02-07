package presentation

import (
	"github.com/Nickolasll/urlshortener/internal/app/config"
	"github.com/Nickolasll/urlshortener/internal/app/domain"
)

func expand(slug string) (domain.Short, error) {
	value, err := repository.GetByShortURL(slug)
	if err != nil {
		return value, err
	}
	return value, nil
}

func shorten(url string, userID string) (domain.Short, error) {
	short := domain.Shorten(url, userID)
	err := repository.Save(short)
	return short, err
}

func batchShorten(batchInput []BatchInput, userID string) []BatchOutput {
	var shorts []domain.Short
	batchOutput := []BatchOutput{}
	for _, batch := range batchInput {
		short := domain.Shorten(batch.OriginalURL, userID)
		shorts = append(shorts, short)
		output := BatchOutput{
			CorrelationID: batch.CorrelationID,
			ShortURL:      *config.SlugEndpoint + short.ShortURL,
		}
		batchOutput = append(batchOutput, output)
	}
	repository.BulkSave(shorts)
	return batchOutput
}

func findURLs(userID string) []FindURLsResult {
	var URLResults []FindURLsResult
	shorts, _ := repository.FindByUserID(userID)
	for _, short := range shorts {
		result := FindURLsResult{
			ShortURL:    *config.SlugEndpoint + short.ShortURL,
			OriginalURL: short.OriginalURL,
		}
		URLResults = append(URLResults, result)
	}
	return URLResults
}

func bulkDelete(shortURLs []string, userID string) {
	for index, shortURL := range shortURLs {
		shortURL = "/" + shortURL
		shortURLs[index] = shortURL
	}
	repository.BulkDelete(shortURLs, userID)
}

func getInternalStats() (GetInternalStatsResult, error) {
	users, urls, err := repository.GetStats()
	results := GetInternalStatsResult{
		URLs:  urls,
		Users: users,
	}
	return results, err
}
