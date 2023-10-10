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

func (r RAMRepository) Get(slug string) (string, bool, error) {
	value, ok := r.urlShortenerMap[slug]
	return value, ok, nil
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

func (r FileRepository) Get(slug string) (string, bool, error) {
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
		return value, ok, nil
	}
	return value, ok, nil
}

func (r FileRepository) Ping() error {
	return nil
}

type PostgresqlRepository struct {
	DSN string
}

func (r PostgresqlRepository) Ping() error {
	db, err := sql.Open("pgx", r.DSN)
	if err != nil {
		return err
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

func (r PostgresqlRepository) Init() error {
	db, err := sql.Open("pgx", r.DSN)
	if err != nil {
		return err
	}
	defer db.Close()
	query := `
		CREATE TABLE IF NOT EXISTS shortener (
		id uuid NOT NULL PRIMARY KEY
		, short_url varchar(10) NOT NULL
		, original_url varchar(100) NOT NULL
	)`
	if _, err = db.ExecContext(context.Background(), query); err != nil {
		return err
	}
	return nil
}

func (r PostgresqlRepository) Save(short domain.Short) error {
	db, err := sql.Open("pgx", r.DSN)
	if err != nil {
		return err
	}
	defer db.Close()
	query := `
		INSERT INTO shortener (
			id
			, short_url
			, original_url
		) VALUES (
			$1::UUID
			, $2::TEXT
			, $3::TEXT
		)`
	_, err = db.ExecContext(
		context.Background(),
		query,
		short.UUID,
		short.ShortURL,
		short.OriginalURL,
	)
	return err
}

func (r PostgresqlRepository) Get(slug string) (string, bool, error) {
	db, err := sql.Open("pgx", r.DSN)
	if err != nil {
		return "", false, err
	}
	defer db.Close()
	query := `
		SELECT
			shortener.original_url	
		FROM
			shortener
		WHERE
			shortener.short_url = $1::TEXT
		;`
	row := db.QueryRowContext(context.Background(), query, slug)
	var originalURL string
	err = row.Scan(&originalURL)
	if err != nil {
		return "", false, err
	}
	return originalURL, true, nil
}

func GetRepository() domain.ShortRepositoryInerface {
	if *config.DatabaseDSN != "" {
		postgres := PostgresqlRepository{DSN: *config.DatabaseDSN}
		postgres.Init()
		return postgres
	} else if *config.FileStoragePath != "" {
		return FileRepository{
			cache:    map[string]string{},
			filePath: *config.FileStoragePath,
		}
	} else {
		return RAMRepository{
			urlShortenerMap: map[string]string{},
		}
	}
}
