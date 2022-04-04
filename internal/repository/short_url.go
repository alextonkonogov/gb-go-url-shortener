package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type shortURL struct {
	ID           int    `json:"id" db:"id"`
	ShortURLCode string `json:"short_url_code" db:"short_url_code"`
	AdminURLCode string `json:"admin_url_code" db:"admin_url_code"`
}

func (r *Repository) NewShortURL(ctx context.Context, dbpool *pgxpool.Pool, shortUrlCode, adminUrlCode string) (id int, err error) {
	query := `INSERT INTO short_urls (short_url_code, admin_url_code) VALUES ($1, $2) RETURNING id`
	err = dbpool.QueryRow(ctx, query, shortUrlCode, adminUrlCode).Scan(&id)
	if err != nil {
		err = fmt.Errorf("SHORT failed to query data: %w", err)
		return
	}

	return
}

func (r *Repository) GetShortURLByAdminIDAndCode(ctx context.Context, dbpool *pgxpool.Pool, adminUrlId, adminUrlCode string) (shortUrl shortURL, err error) {
	row := dbpool.QueryRow(ctx, `
		SELECT short_urls.id, short_urls.short_url_code
		FROM short_urls
		WHERE short_urls.id = $1 AND short_urls.admin_url_code = $2`,
		adminUrlId, adminUrlCode)
	if err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
		return
	}

	err = row.Scan(&shortUrl.ID, &shortUrl.ShortURLCode)
	if err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
		return
	}

	return
}
