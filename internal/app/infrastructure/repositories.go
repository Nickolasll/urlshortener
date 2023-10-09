package infrastructure

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/Nickolasll/urlshortener/internal/app/config"
	"github.com/Nickolasll/urlshortener/internal/app/domain"
)

type RAMRepository struct {
	urlShortenerMap map[string]string
}

func (r RAMRepository) Save(short domain.Short) error {
	r.urlShortenerMap[short.ShortURL] = short.OriginalURL
	return nil
}

func (r RAMRepository) Get(slug string) (string, bool) {
	value, ok := r.urlShortenerMap[slug]
	return value, ok
}

func (r RAMRepository) Ping() error {
	return nil
}

type FileRepository struct {
	filePath string
	cache    map[string]string
}

func (r FileRepository) Save(short domain.Short) error {
	r.cache[short.ShortURL] = short.OriginalURL
	file, err := os.OpenFile(r.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
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
	return nil
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

func (r FileRepository) Ping() error {
	return nil
}

type PostgresqlRepository struct {
	dsn string
}

func (r PostgresqlRepository) Ping() error {
	db, err := sql.Open("pgx", r.dsn)
	if err != nil {
		return err
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

var Postgres = PostgresqlRepository{
	dsn: *config.DatabaseDSN,
}

func GetRepository() domain.ShortRepositoryInerface {
	if *config.FileStoragePath == "" {
		return RAMRepository{
			urlShortenerMap: map[string]string{},
		}
	} else {
		return FileRepository{
			cache:    map[string]string{},
			filePath: *config.FileStoragePath,
		}
	}
}
