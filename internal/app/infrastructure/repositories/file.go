package repositories

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/Nickolasll/urlshortener/internal/app/domain"
)

type FileRepository struct {
	FilePath string
	Cache    map[string]string
}

func (r FileRepository) cache(short domain.Short) {
	r.Cache[short.ShortURL] = short.OriginalURL
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

func (r FileRepository) Get(slug string) (string, bool, error) {
	value, ok := r.Cache[slug]
	if !ok {
		file, _ := os.OpenFile(r.FilePath, os.O_RDONLY|os.O_CREATE, 0666)
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			var short domain.Short
			json.Unmarshal(scanner.Bytes(), &short)
			r.Cache[short.ShortURL] = short.OriginalURL
		}
		value, ok := r.Cache[slug]
		return value, ok, nil
	}
	return value, ok, nil
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