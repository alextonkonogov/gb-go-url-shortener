package pgstore

import (
	"context"
	"database/sql"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/entities/statistics"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/usecases/app/repos/repoStatistics"

	_ "github.com/jackc/pgx/v4/stdlib" // Postgresql driver
)

var _ repoStatistics.StatisticsStore = &Statistics{}

type DBPgStatistics struct {
	ID        int64
	CreatedAt string
	Long      string
	Short     string
	Admin     string
}

type Statistics struct {
	db *sql.DB
}

func NewStatistics(dsn string) (*Statistics, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS statistics (
		id bigint primary key generated always as identity,
		viewed_at timestamptz NOT NULL,
		ip  varchar NOT NULL,
		url_id 	integer not null,
		constraint url_id_usages_fk foreign key (url_id) references urls (id)
	)`)
	if err != nil {
		db.Close()
		return nil, err
	}
	st := &Statistics{
		db: db,
	}
	return st, nil
}

func (st *Statistics) Close() {
	st.db.Close()
}

func (st *Statistics) Read(ctx context.Context, s statistics.Statistics) (*statistics.Statistics, error) {
	dbu := &DBPgStatistics{
		Admin: s.Admin,
	}

	err := st.db.QueryRowContext(ctx, `SELECT long, short FROM urls WHERE admin = $1`, dbu.Admin).Scan(&s.Long, &s.Short)
	if err != nil {
		return nil, err
	}

	return &s, nil
}
