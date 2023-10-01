package infrastructure

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/Nickolasll/urlshortener/internal/app/config"
	"github.com/Nickolasll/urlshortener/internal/app/domain"
)

type RAMRepository struct {
	urlShortenerMap map[string]string
}

func (r RAMRepository) Save(short domain.Short) {
	r.urlShortenerMap[short.ShortURL] = short.OriginalURL
}

func (r RAMRepository) Get(slug string) (string, bool) {
	value, ok := r.urlShortenerMap[slug]
	return value, ok
}

type FileRepository struct {
	filePath string
	cache    map[string]string
}

func (r FileRepository) Save(short domain.Short) {
	r.cache[short.ShortURL] = short.OriginalURL
	file, _ := os.OpenFile(r.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	defer file.Close()
	data, _ := json.Marshal(short)
	data = append(data, '\n')
	file.Write(data)
}

func (r FileRepository) Get(slug string) (string, bool) {
	value, ok := r.cache[slug]
	if !ok {
		file, _ := os.OpenFile(r.filePath, os.O_RDONLY|os.O_CREATE, 0666)
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			var short domain.Short
			json.Unmarshal(scanner.Bytes(), &short)
			r.cache[short.ShortURL] = short.OriginalURL
		}
		value, ok := r.cache[slug]
		return value, ok
	}
	return value, ok
}

var Repository domain.ShortRepositoryInerface = RAMRepository{
	urlShortenerMap: map[string]string{},
}

func InitRepository() {
	if *config.FileStoragePath != "" {
		Repository = FileRepository{
			cache:    map[string]string{},
			filePath: *config.FileStoragePath,
		}
	}
}
