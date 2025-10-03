package postgres

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Postgres struct {
	urlConnect string

	*pgxpool.Pool
}

func New(urlConnect string) *Postgres {
	return &Postgres{urlConnect: urlConnect}
}

func (pg *Postgres) Connect(ctx context.Context) error {
	cfg, err := pgxpool.ParseConfig(pg.urlConnect)
	if err != nil {
		return err
	}

	pool, err := pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		return err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return err
	}

	pg.Pool = pool
	return nil
}

func (pg *Postgres) Close() {
	if pg.Pool != nil {
		pg.Pool.Close()
	}
}
