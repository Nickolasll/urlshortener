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

// PostgresqlRepository - Имплементация postgres репозитория
type PostgresqlRepository struct {
	// DSN - имя источника данных
	DSN string
	// Timeout - время ожидания операции
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

// Ping - Проверка работоспособности
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

// Init - Инициализация репозитория
func (r PostgresqlRepository) Init() error {
	db, ctx, cancel, err := r.openConn()
	if err != nil {
		return err
	}
	defer db.Close()
	defer cancel()
	query := `
		CREATE TABLE IF NOT EXISTS shortener (
		id uuid NOT NULL PRIMARY KEY
		, short_url varchar(10) NOT NULL
		, original_url varchar(100) NOT NULL UNIQUE
		, user_id uuid NOT NULL
		, deleted boolean NOT NULL
	)`
	if _, err = db.ExecContext(ctx, query); err != nil {
		return err
	}
	query = "CREATE INDEX short_url_idx on shortener(short_url)"
	if _, err = db.ExecContext(ctx, query); err != nil {
		return err
	}
	query = "CREATE INDEX user_id_idx on shortener(user_id)"
	if _, err = db.ExecContext(ctx, query); err != nil {
		return err
	}
	return nil
}

// Save - Сохранить сокращенную ссылку
func (r PostgresqlRepository) Save(short domain.Short) error {
	db, ctx, cancel, err := r.openConn()
	if err != nil {
		return err
	}
	defer db.Close()
	defer cancel()
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
	_, err = db.ExecContext(
		ctx,
		query,
		short.UUID,
		short.ShortURL,
		short.OriginalURL,
		short.UserID,
		short.Deleted,
	)
	return err
}

// GetByShortURL - Получить объект сокращенной ссылки по значению
func (r PostgresqlRepository) GetByShortURL(slug string) (domain.Short, error) {
	var short domain.Short
	db, ctx, cancel, err := r.openConn()
	if err != nil {
		return short, err
	}
	defer db.Close()
	defer cancel()
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
	row := db.QueryRowContext(ctx, query, slug)
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

// GetShortURL - Получить сокращенную ссылку по несокращенному значению
func (r PostgresqlRepository) GetShortURL(originalURL string) (string, error) {
	var short string
	db, ctx, cancel, err := r.openConn()
	if err != nil {
		return short, err
	}
	defer db.Close()
	defer cancel()
	query := `
		SELECT
			shortener.short_url	
		FROM
			shortener
		WHERE
			shortener.original_url = $1::TEXT
		;`
	row := db.QueryRowContext(ctx, query, originalURL)
	err = row.Scan(&short)
	if err != nil {
		return "", err
	}
	return short, nil
}

// BulkSave - Сохранить пачку сокращенных ссылок
func (r PostgresqlRepository) BulkSave(shorts []domain.Short) error {
	db, ctx, cancel, err := r.openConn()
	if err != nil {
		return err
	}
	defer db.Close()
	defer cancel()
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
	_, err = db.ExecContext(ctx, query, vals...)
	return err
}

// FindByUserID - Получить список сокращенных ссылок по идентификатору пользователя
func (r PostgresqlRepository) FindByUserID(userID string) ([]domain.Short, error) {
	shorts := []domain.Short{}
	db, ctx, cancel, err := r.openConn()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	defer cancel()
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

	rows, err := db.QueryContext(ctx, query, userID)
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

// BulkDelete - Удалить пачку сокращенных ссылок
func (r PostgresqlRepository) BulkDelete(shortURLs []string, userID string) error {
	db, ctx, cancel, err := r.openConn()
	if err != nil {
		return err
	}
	defer db.Close()
	defer cancel()
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
	_, err = db.ExecContext(ctx, query)
	return err
}

// GetStats получение количества пользователей и сокращенных URL
func (r PostgresqlRepository) GetStats() (int, int, error) {
	var uniqueUserIds []string
	shortCount := 0
	db, ctx, cancel, err := r.openConn()
	if err != nil {
		return len(uniqueUserIds), shortCount, err
	}
	defer db.Close()
	defer cancel()
	query := `
		SELECT
			shortener.id
			, shortener.short_url
			, shortener.original_url
			, shortener.user_id
			, shortener.deleted
		FROM
			shortener
		;`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return len(uniqueUserIds), shortCount, err
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
			shortCount += 1
			if !contains(uniqueUserIds, short.UserID) {
				uniqueUserIds = append(uniqueUserIds, short.UserID)
			}
		}
	}
	return len(uniqueUserIds), shortCount, nil
}
