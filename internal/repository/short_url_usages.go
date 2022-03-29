package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type shortURLusages struct {
	Id         int       `json:"id" db:"id"`
	Date       time.Time `json:"date" db:"date"`
	Ip         string    `json:"ip" db:"ip"`
	ShortURLid int       `json:"short_url_id" db:"short_url_id"`
}

func (r *Repository) NewShortUrlUsage(ctx context.Context, dbpool *pgxpool.Pool, ip string, shortUrlId string) (id int, err error) {
	query := `INSERT INTO short_url_usages (date, ip, short_url_id) VALUES ($1, $2, $3) RETURNING id`
	err = dbpool.QueryRow(ctx, query, time.Now().UTC(), ip, shortUrlId).Scan(&id)
	if err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
		return
	}

	return
}

func (r *Repository) GetLongUrlCountByAdminIdAndCode(ctx context.Context, dbpool *pgxpool.Pool, adminUrlId, adminUrlCode string) (count int, err error) {
	row := dbpool.QueryRow(ctx, `
		SELECT count(*)
		FROM short_url_usages 
		LEFT JOIN short_urls
		ON short_url_usages.short_url_id = short_urls.id
		WHERE short_urls.id = $1 AND short_urls.admin_url_code = $2`,
		adminUrlId, adminUrlCode)
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
