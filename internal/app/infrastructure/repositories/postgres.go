package repositories

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/Nickolasll/urlshortener/internal/app/domain"
)

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
