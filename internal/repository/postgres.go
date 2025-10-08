package repository

import (
	"context"
	"errors"
	"log"
	adapterpg "shortener/internal/adapter/postgres"
	"shortener/internal/model"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
)

type Postgres struct {
	conn *adapterpg.Postgres
}

func NewPostgres(conn *adapterpg.Postgres) *Postgres {
	return &Postgres{conn: conn}
}

func (p *Postgres) Post(ctx context.Context, link model.Link) (*model.Link, error) {
	const query = `
		INSERT INTO links (original_url, shorten_url) 
		VALUES ($1, $2)
		ON CONFLICT (shorten_url) DO NOTHING
		RETURNING *;
	`
	var result model.Link
	err := pgxscan.Get(ctx, p.conn.Pool, &result, query, link.OriginalUrl, link.ShortenUrl)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.ErrorNonUniq
	}

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (p *Postgres) Get(ctx context.Context, shortenUrl string) (string, error) {
	const query = `
		SELECT original_url FROM links 
		WHERE shorten_url = $1
	`
	var result string
	err := p.conn.Pool.QueryRow(ctx, query, shortenUrl).Scan(&result)

	if errors.Is(err, pgx.ErrNoRows) {
		return "", model.ErrorNotFound
	}

	if err != nil {
		return "", err
	}

	return result, nil
}

func (p *Postgres) IncrementVisits(ctx context.Context, shortenUrl string) error {
	const query = `
		UPDATE links 
		SET visits = visits + 1 
		WHERE shorten_url = $1;
	`

	_, err := p.conn.Pool.Exec(ctx, query, shortenUrl)
	if err != nil {
		log.Printf("failed to increment visits for %q: %v", shortenUrl, err)
		return model.ErrorIncrementVisits
	}

	return nil
}
