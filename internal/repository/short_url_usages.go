package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

func (r *Repository) NewShortURLUsage(ctx context.Context, dbpool *pgxpool.Pool, IP string, shortURLID string) (ID int, err error) {
	query := `INSERT INTO short_url_usages (date, IP, short_url_id) VALUES ($1, $2, $3) RETURNING ID`
	err = dbpool.QueryRow(ctx, query, time.Now().UTC(), IP, shortURLID).Scan(&ID)
	if err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
		return
	}

	return
}

func (r *Repository) GetLongURLCountByAdminIDAndCode(ctx context.Context, dbpool *pgxpool.Pool, adminURLID, adminURLCode string) (count int, err error) {
	row := dbpool.QueryRow(ctx, `
		SELECT count(*)
		FROM short_url_usages 
		LEFT JOIN short_urls
		ON short_url_usages.short_url_id = short_urls.id
		WHERE short_urls.id = $1 AND short_urls.admin_url_code = $2`,
		adminURLID, adminURLCode)
	if err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
		return
	}

	err = row.Scan(&count)
	if err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
		return
	}

	return
}
