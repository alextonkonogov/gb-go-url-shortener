package storage

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

func InitDBConn(ctx context.Context) (dbpool *pgxpool.Pool, err error) {
	url := "postgres://postgres:password@pg_db:5432/postgres?sslmode=disable"
	fmt.Println(url)

	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		err = fmt.Errorf("failed to parse pg config: %w", err)
		return
	}

	cfg.MaxConns = int32(5)
	cfg.MinConns = int32(1)
	cfg.HealthCheckPeriod = 1 * time.Minute
	cfg.MaxConnLifetime = 24 * time.Hour
	cfg.MaxConnIdleTime = 30 * time.Minute
	cfg.ConnConfig.ConnectTimeout = 1 * time.Second
	cfg.ConnConfig.DialFunc = (&net.Dialer{
		KeepAlive: cfg.HealthCheckPeriod,
		Timeout:   cfg.ConnConfig.ConnectTimeout,
	}).DialContext

	dbpool, err = pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		err = fmt.Errorf("failed to connect config: %w", err)
		return
	}

	return
}

func InitTables(ctx context.Context, dbpool *pgxpool.Pool) (err error) {
	query := `
				create table if not exists short_urls
				(
					id        		bigint primary key generated always as identity,
					short_url_code 	varchar(50) not null,
					admin_url_code 	varchar(50) not null
				);

				create table if not exists long_urls
				(
					id      		bigint primary key generated always as identity,
					long_url  		varchar(500) not null,
					short_url_id 	integer not null,
					constraint short_url_id_fk foreign key (short_url_id) references short_urls (id)
				);

				create table if not exists short_url_usages
				(
					id        		bigint primary key generated always as identity,
					date			timestamp with time zone not null,
					ip     			varchar(15) not null,
					short_url_id 	integer not null,
					constraint short_url_id_usages_fk foreign key (short_url_id) references short_urls (id)
				);
				`

	_, err = dbpool.Exec(ctx, query)

	return
}
