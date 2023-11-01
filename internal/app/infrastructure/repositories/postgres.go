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
	DSN     string
	Timeout time.Duration
}

func (r PostgresqlRepository) openConn() (*sql.DB, context.Context, context.CancelFunc, error) {
	db, err := sql.Open("pgx", r.DSN)
	if err != nil {
		return nil, nil, nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), r.Timeout)
	return db, ctx, cancel, err
}

func (r PostgresqlRepository) exec(query string, args ...any) error {
	db, ctx, cancel, err := r.openConn()
	if err != nil {
		return err
	}
	defer db.Close()
	defer cancel()
	_, err = db.ExecContext(ctx, query, args...)
	return err
}

func (r PostgresqlRepository) queryRow(query string, args ...any) (*sql.Row, error) {
	db, ctx, cancel, err := r.openConn()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	defer cancel()
	row := db.QueryRowContext(ctx, query, args...)
	return row, nil
}

func (r PostgresqlRepository) query(query string, args ...any) (*sql.Rows, error) {
	db, ctx, cancel, err := r.openConn()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	defer cancel()
	return db.QueryContext(ctx, query, args...)
}

func (r PostgresqlRepository) Ping() error {
	db, ctx, cancel, err := r.openConn()
	if err != nil {
		return err
	}
	defer db.Close()
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

func (r PostgresqlRepository) Init() error {
	query := `
		CREATE TABLE IF NOT EXISTS shortener (
		id uuid NOT NULL PRIMARY KEY
		, short_url varchar(10) NOT NULL
		, original_url varchar(100) NOT NULL UNIQUE
		, user_id uuid NOT NULL
		, deleted boolean NOT NULL
	)`
	if err := r.exec(query); err != nil {
		return err
	}
	query = "CREATE INDEX short_url_idx on shortener(short_url)"
	if err := r.exec(query); err != nil {
		return err
	}
	query = "CREATE INDEX user_id_idx on shortener(user_id)"
	if err := r.exec(query); err != nil {
		return err
	}
	return nil
}

func (r PostgresqlRepository) Save(short domain.Short) error {
	query := `
		INSERT INTO shortener (
			id
			, short_url
			, original_url
			, user_id
			, deleted
		) VALUES (
			$1::UUID
			, $2::TEXT
			, $3::TEXT
			, $4::UUID
			, $5::BOOLEAN
		)`
	return r.exec(
		query,
		short.UUID,
		short.ShortURL,
		short.OriginalURL,
		short.UserID,
		short.Deleted,
	)
}

func (r PostgresqlRepository) GetByShortURL(slug string) (domain.Short, error) {
	var short domain.Short
	query := `
		SELECT
			shortener.id
			, shortener.short_url
			, shortener.original_url
			, shortener.user_id
			, shortener.deleted
		FROM
			shortener
		WHERE
			shortener.short_url = $1::TEXT
		;`
	row, err := r.queryRow(query, slug)
	if err != nil {
		return short, err
	}
	err = row.Scan(
		&short.UUID,
		&short.ShortURL,
		&short.OriginalURL,
		&short.UserID,
		&short.Deleted,
	)
	if err != nil {
		return short, err
	}
	return short, nil
}

func (r PostgresqlRepository) GetShortURL(originalURL string) (string, error) {
	var short string
	query := `
		SELECT
			shortener.short_url	
		FROM
			shortener
		WHERE
			shortener.original_url = $1::TEXT
		;`
	row, err := r.queryRow(query, originalURL)
	if err != nil {
		return "", err
	}
	err = row.Scan(&short)
	if err != nil {
		return "", err
	}
	return short, nil
}

func (r PostgresqlRepository) BulkSave(shorts []domain.Short) error {
	vals := []interface{}{}
	counter := 1
	query := `
		INSERT INTO shortener(
			id
			, short_url
			, original_url
			, user_id
			, deleted
		) 
		VALUES `
	for _, short := range shorts {
		query += "(" +
			"$" + strconv.Itoa(counter) + "::UUID" +
			", $" + strconv.Itoa(counter+1) + "::TEXT" +
			", $" + strconv.Itoa(counter+2) + "::TEXT" +
			", $" + strconv.Itoa(counter+3) + "::UUID" +
			", $" + strconv.Itoa(counter+4) + "::BOOLEAN" +
			"),"
		vals = append(
			vals,
			short.UUID,
			short.ShortURL,
			short.OriginalURL,
			short.UserID,
			short.Deleted,
		)
		counter += 5
	}
	query = strings.TrimSuffix(query, ",")
	return r.exec(query, vals...)
}

func (r PostgresqlRepository) FindByUserID(userID string) ([]domain.Short, error) {
	shorts := []domain.Short{}
	query := `
		SELECT
			shortener.id
			, shortener.short_url
			, shortener.original_url
			, shortener.user_id
			, shortener.deleted
		FROM
			shortener
		WHERE
			shortener.user_id = $1::UUID
		;`
	rows, err := r.query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var short domain.Short
		err = rows.Scan(
			&short.UUID,
			&short.ShortURL,
			&short.OriginalURL,
			&short.UserID,
			&short.Deleted,
		)
		if err == nil && rows.Err() == nil {
			shorts = append(shorts, short)
		}
	}
	return shorts, nil
}

func (r PostgresqlRepository) BulkDelete(shortURLs []string, userID string) error {
	query := `
		UPDATE shortener SET
			deleted = s.deleted
		FROM (VALUES `
	for _, shortUUID := range shortURLs {
		query += "(" +
			"'" + shortUUID + "'" +
			", true" +
			"),"
	}
	query = strings.TrimSuffix(query, ",")
	query += `
		) AS s(short_url, deleted)
		WHERE 
			s.short_url = shortener.short_url
			AND shortener.user_id = '` + userID + "'"
	return r.exec(query)
}
