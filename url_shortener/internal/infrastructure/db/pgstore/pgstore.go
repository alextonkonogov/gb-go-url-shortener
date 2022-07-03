package pgstore

import (
	"context"
	"database/sql"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/entities/url"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/usecases/app/repos/userrepo"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib" // Postgresql driver
)

var _ userrepo.URLStore = &URL{}

type DBPgURL struct {
	ID        int
	CreatedAt time.Time
	Long      string
	Short     string
}

type URL struct {
	db *sql.DB
}

func NewUsers(dsn string) (*URL, error) {
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
		id integer NOT NULL,
		created_at timestamptz NOT NULL,
		long varchar NOT NULL,
		short varchar NULL
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

func (us *URL) Create(ctx context.Context, u url.URL) (*int, error) {
	dbu := &DBPgURL{
		ID:        u.ID,
		CreatedAt: time.Now(),
		Short:     u.Short,
		Long:      u.Long,
	}

	_, err := us.db.ExecContext(ctx, `INSERT INTO urls (id, created_at, long, short)
	values ($1, $2, $3, $4)`,
		dbu.ID,
		dbu.CreatedAt,
		dbu.Long,
		dbu.Short,
	)
	if err != nil {
		return nil, err
	}

	return &u.ID, nil
}
