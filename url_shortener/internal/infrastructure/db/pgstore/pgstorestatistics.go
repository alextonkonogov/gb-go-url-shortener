package pgstore

import (
	"context"
	"database/sql"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/entities/statistics"
	"github.com/alextonkonogov/gb-go-url-shortener/url_shortener/internal/usecases/app/repos/repostatistics"
	"github.com/sirupsen/logrus"

	_ "github.com/jackc/pgx/v4/stdlib" // Postgresql driver
)

var _ repostatistics.StatisticsStore = &Statistics{}

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
	db  *sql.DB
	log *logrus.Logger
}

func NewStatistics(dsn string, log *logrus.Logger) (*Statistics, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Error("err when opening connection: ", err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		log.Error("err when pinging: ", err)
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
		log.Error("err when creating table: ", err)
		return nil, err
	}
	st := &Statistics{
		db:  db,
		log: log,
	}
	return st, nil
}

func (st *Statistics) Close() {
	st.db.Close()
}

func (st *Statistics) Create(ctx context.Context, urlID int64) error {
	_, err := st.db.ExecContext(ctx, `INSERT INTO statistics (url_id) values ($1)`, urlID)
	if err != nil {
		st.log.Error("err when inserting: ", err)
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
		st.log.Error("err when selecting: ", err)
		return nil, err
	}

	return &s, nil
}

func (st *Statistics) Update(ctx context.Context, s statistics.Statistics, urlID int64) (*statistics.Statistics, error) {
	dbs := &DBPgStatistics{
		IP:       s.IP,
		ViewedAt: s.Viewed,
		URLID:    urlID,
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
		st.log.Error("err when updating: ", err)
		return nil, err
	}

	return &s, nil
}
