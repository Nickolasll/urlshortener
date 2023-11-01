package repositories

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/Nickolasll/urlshortener/internal/app/domain"
)

type FileRepository struct {
	FilePath string
	Cache    map[string]domain.Short
}

func (r FileRepository) loadCache() error {
	file, err := os.OpenFile(r.FilePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var short domain.Short
		json.Unmarshal(scanner.Bytes(), &short)
		r.Cache[short.ShortURL] = short
	}
	return nil
}

func (r FileRepository) cache(short domain.Short) {
	r.Cache[short.ShortURL] = short
}

func (r FileRepository) Save(short domain.Short) error {
	file, err := os.OpenFile(r.FilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := json.Marshal(short)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	file.Write(data)
	r.cache(short)
	return nil
}

func (r FileRepository) GetByShortURL(shortURL string) (domain.Short, error) {
	short, ok := r.Cache[shortURL]
	if !ok {
		r.loadCache()
		short := r.Cache[shortURL]
		return short, nil
	}
	return short, nil
}

func (r FileRepository) Ping() error {
	return nil
}

func (r FileRepository) BulkSave(shorts []domain.Short) error {

	file, err := os.OpenFile(r.FilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	data := []byte{}
	for _, short := range shorts {
		serialized, err := json.Marshal(short)
		if err != nil {
			return err
		}
		data = append(data, append(serialized, '\n')...)
		r.cache(short)
	}
	file.Write(data)
	return nil

}

func (r FileRepository) GetShortURL(originalURL string) (string, error) {
	key, ok := originalURLKeyMap(r.Cache, originalURL)
	if !ok {
		r.loadCache()
		key, _ = originalURLKeyMap(r.Cache, originalURL)
	}
	return key, nil
}

func (r FileRepository) FindByUserID(userID string) ([]domain.Short, error) {
	shorts := []domain.Short{}
	for _, short := range r.Cache {
		if short.UserID == userID {
			shorts = append(shorts, short)
		}
	}
	return shorts, nil
}

func (r FileRepository) BulkDelete(shortURLs []string, userID string) error {
	for key, short := range r.Cache {
		if contains(shortURLs, short.ShortURL) && short.UserID == userID {
			short.Deleted = true
			r.Cache[key] = short
		}
	}
	file, err := os.OpenFile(r.FilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()
	data := []byte{}
	for _, short := range r.Cache {
		serialized, err := json.Marshal(short)
		if err != nil {
			return err
		}
		data = append(data, append(serialized, '\n')...)
	}
	file.Write(data)
	return nil
}
