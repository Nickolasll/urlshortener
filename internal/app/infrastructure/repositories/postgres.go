package repositories

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
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
		, original_url varchar(100) NOT NULL UNIQUE
	)`
	if _, err = db.ExecContext(context.Background(), query); err != nil {
		return err
	}
	query = "CREATE INDEX short_url_idx on shortener(short_url)"
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

func (r PostgresqlRepository) GetOriginalURL(slug string) (string, error) {
	var originalURL string
	db, err := sql.Open("pgx", r.DSN)
	if err != nil {
		return "", err
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
	err = row.Scan(&originalURL)
	if err != nil {
		return "", err
	}
	return originalURL, nil
}

func (r PostgresqlRepository) GetShortURL(originalURL string) (string, error) {
	var short string
	db, err := sql.Open("pgx", r.DSN)
	if err != nil {
		return "", err
	}
	defer db.Close()
	query := `
		SELECT
			shortener.short_url	
		FROM
			shortener
		WHERE
			shortener.original_url = $1::TEXT
		;`
	row := db.QueryRowContext(context.Background(), query, originalURL)
	err = row.Scan(&short)
	if err != nil {
		return "", err
	}
	return short, nil
}

func (r PostgresqlRepository) BulkSave(shorts []domain.Short) error {
	db, err := sql.Open("pgx", r.DSN)
	if err != nil {
		return err
	}
	defer db.Close()
	query := `
		INSERT INTO shortener(
			id
			, short_url
			, original_url
		) 
		VALUES `
	vals := []interface{}{}
	counter := 1
	// В самом pgx есть функция CopyFrom и синтаксис мне нравится больше
	// Для единообразия все сделано через sql.Open
	for _, short := range shorts {
		query += "(" +
			"$" + strconv.Itoa(counter) + "::UUID" +
			", $" + strconv.Itoa(counter+1) + "::TEXT" +
			", $" + strconv.Itoa(counter+2) + "::TEXT" +
			"),"
		vals = append(vals, short.UUID, short.ShortURL, short.OriginalURL)
		counter += 3
	}
	query = strings.TrimSuffix(query, ",")
	statement, err := db.Prepare(query)
	if err != nil {
		return err
	}
	if _, err := statement.Exec(vals...); err != nil {
		return err
	}

	return nil
}
