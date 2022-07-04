package pgstore

import (
	"context"
	"database/sql"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/entities/url"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/usecases/app/repos/URLrepo"

	_ "github.com/jackc/pgx/v4/stdlib" // Postgresql driver
)

var _ URLrepo.URLStore = &URL{}

type DBPgURL struct {
	ID        int64
	CreatedAt string
	Long      string
	Short     string
	Admin     string
}

type URL struct {
	db *sql.DB
}

func NewURL(dsn string) (*URL, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS urls (
		id bigint primary key generated always as identity,
		created_at timestamptz NOT NULL,
		long varchar NOT NULL,
		short varchar NOT NULL,
		admin varchar NOT NULL
	)`)
	if err != nil {
		db.Close()
		return nil, err
	}
	us := &URL{
		db: db,
	}
	return us, nil
}

func (us *URL) Close() {
	us.db.Close()
}

func (us *URL) Create(ctx context.Context, u url.URL) (*int64, error) {
	dbu := &DBPgURL{
		ID:        u.ID,
		CreatedAt: u.Created,
		Short:     u.Short,
		Long:      u.Long,
		Admin:     u.Admin,
	}

	err := us.db.QueryRowContext(ctx, `INSERT INTO urls (created_at, long, short, admin) 
		values ($1, $2, $3, $4) RETURNING id`,
		dbu.CreatedAt,
		dbu.Long,
		dbu.Short,
		dbu.Admin,
	).Scan(&dbu.ID)

	if err != nil {
		return nil, err
	}

	return &dbu.ID, nil
}

func (us *URL) Read(ctx context.Context, u url.URL) (*string, error) {
	dbu := &DBPgURL{
		Short: u.Short,
	}

	err := us.db.QueryRowContext(ctx, `SELECT long FROM urls WHERE short = $1`, dbu.Short).Scan(&dbu.Long)
	if err != nil {
		return nil, err
	}

	return &dbu.Long, nil
}
