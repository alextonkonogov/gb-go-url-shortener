package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

type longURL struct {
	Id         int    `json:"id" db:"id"`
	LongUrl    string `json:"long_url" db:"long_url"`
	ShortURLid int    `json:"short_url_id" db:"short_url_id"`
}

func (r *Repository) NewLongUrl(ctx context.Context, dbpool *pgxpool.Pool, longUrl string, shortUrlId int) (id int, err error) {
	query := `INSERT INTO long_urls (long_url, short_url_id) VALUES ($1, $2) RETURNING id`
	err = dbpool.QueryRow(ctx, query, longUrl, shortUrlId).Scan(&id)
	if err != nil {
		err = fmt.Errorf("LONG failed to query data: %w", err)
		return
	}

	return
}

func (r *Repository) GetLongUrlByShortIdAndCode(ctx context.Context, dbpool *pgxpool.Pool, shortUrlId, shortUrlCode string) (longUrl longURL, err error) {
	row := dbpool.QueryRow(ctx, `
		SELECT long_urls.id, long_urls.long_url
		FROM long_urls 
		LEFT JOIN short_urls
		ON long_urls.short_url_id = short_urls.id
		WHERE short_urls.id = $1 AND short_urls.short_url_code = $2`,
		shortUrlId, shortUrlCode)
	if err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
		return
	}

	err = row.Scan(&longUrl.Id, &longUrl.LongUrl)
	if err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
		return
	}

	return
}
