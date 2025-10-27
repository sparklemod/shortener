package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	adapterpg "shortener/internal/adapter/postgres"
	"shortener/internal/model"
	"strings"

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

func (p *Postgres) Get(ctx context.Context, shortenUrl string) (*model.Link, error) {
	const query = `
		SELECT * FROM links 
		WHERE shorten_url = $1
	`
	var result model.Link
	err := pgxscan.Get(ctx, p.conn.Pool, &result, query, shortenUrl)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.ErrorNonUniq
	}

	if err != nil {
		return nil, err
	}

	return &result, nil
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

func (p *Postgres) FilterLinks(ctx context.Context, f model.FilterLinksInput) ([]model.Link, error) {
	query := `
		SELECT id, original_url, shorten_url, visits, is_active
		FROM links
	`
	where := []string{}
	args := []interface{}{}
	i := 1

	if f.IsActive != nil {
		where = append(where, fmt.Sprintf("is_active = $%d", i))
		args = append(args, *f.IsActive)
		i++
	}

	if f.ShortenUrl != nil {
		where = append(where, fmt.Sprintf("shorten_url = $%d", i))
		args = append(args, *f.ShortenUrl)
		i++
	}

	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}

	sortField := "original_url"
	switch f.SortBy {
	case "visits":
		sortField = "visits"
	}

	sortOrder := "ASC"
	if strings.ToUpper(f.SortOrder) == "DESC" {
		sortOrder = "DESC"
	}

	query += fmt.Sprintf(" ORDER BY %s %s", sortField, sortOrder)

	if f.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", i)
		args = append(args, f.Limit)
		i++
	}
	if f.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", i)
		args = append(args, f.Offset)
		i++
	}

	rows, err := p.conn.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	var links []model.Link
	for rows.Next() {
		var l model.Link
		if err := rows.Scan(&l.Id, &l.OriginalUrl, &l.ShortenUrl, &l.Visits, &l.IsActive); err != nil {
			return nil, err
		}
		links = append(links, l)
	}

	return links, rows.Err()
}
