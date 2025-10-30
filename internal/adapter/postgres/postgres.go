package postgres

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pressly/goose/v3"

	_ "github.com/jackc/pgx/v5/stdlib"
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

func (pg *Postgres) Close(ctx context.Context) error {
	if pg.Pool != nil {
		pg.Pool.Close()
	}

	return nil
}

func (pg *Postgres) RunMigrations(migrationsDir string) error {
	db, err := goose.OpenDBWithDriver("postgres", pg.urlConnect)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := goose.Up(db, migrationsDir); err != nil {
		return err
	}

	return nil
}
