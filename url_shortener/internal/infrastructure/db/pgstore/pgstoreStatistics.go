package pgstore

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/entities/statistics"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/usecases/app/repos/repoStatistics"

	_ "github.com/jackc/pgx/v4/stdlib" // Postgresql driver
)

var _ repoStatistics.StatisticsStore = &Statistics{}

type DBPgStatistics struct {
	URLID    int64
	ViewedAt string
	IP       string
	Count    int64
	Long     string
	Short    string
	Admin    string
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
		viewed_at varchar DEFAULT '' NOT NULL,
		ip  varchar DEFAULT '' NOT NULL,
		count integer DEFAULT 0 NOT NULL,
		url_id 	integer NOT NULL,
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

func (st *Statistics) Create(ctx context.Context, URLID int64) error {
	_, err := st.db.ExecContext(ctx, `INSERT INTO statistics (url_id) values ($1)`, URLID)
	if err != nil {
		return err
	}

	return nil
}

func (st *Statistics) Read(ctx context.Context, s statistics.Statistics) (*statistics.Statistics, error) {
	dbs := &DBPgStatistics{
		Admin: s.Admin,
	}

	err := st.db.QueryRowContext(ctx, `
	SELECT 
	       st.ip, 
	       st.count, 
	       st.viewed_at, 
	       urls.short,
	       urls.long 
	FROM urls 
	    LEFT JOIN statistics AS st 
	    ON urls.id = st.url_id 
	WHERE urls.admin = $1`,
		dbs.Admin,
	).Scan(&s.IP, &s.Count, &s.Viewed, &s.Short, &s.Long)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &s, nil
}

func (st *Statistics) Update(ctx context.Context, s statistics.Statistics, URLID int64) (*statistics.Statistics, error) {
	dbs := &DBPgStatistics{
		IP:       s.IP,
		ViewedAt: s.Viewed,
		URLID:    URLID,
	}

	err := st.db.QueryRowContext(ctx, `
	UPDATE statistics SET 
		ip = $1, 
		count = count +1, 
		viewed_at = $2 
	WHERE url_id = $3 RETURNING count`,
		dbs.IP,
		dbs.ViewedAt,
		dbs.URLID,
	).Scan(&dbs.Count)

	s.Count = dbs.Count

	if err != nil {
		return nil, err
	}

	return &s, nil
}
