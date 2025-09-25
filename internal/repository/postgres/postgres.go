package postgres

import (
	"context"
	adapterpg "shortener/internal/adapter/postgres"
	"shortener/internal/repository/domain"
)

type Postgres struct {
	conn *adapterpg.Postgres
}

func NewPostgres(conn *adapterpg.Postgres) *Postgres {
	return &Postgres{conn: conn}
}

func (p *Postgres) CreateLink(ctx context.Context, link domain.CreateLink) (string, error) {
	const query = `
		INSERT INTO links (original_link, redirect_link) 
		VALUES ($1, $2)
		RETURNING id;
	`
	var id string
	if err := p.conn.Pool.QueryRow(ctx, query, link.OriginalLink, link.RedirectLink).Scan(&id); err != nil {
		return "", err
	}

	return id, nil
}

func (p *Postgres) CheckLinkIfExist(ctx context.Context, redirectLink string) (bool, error) {
	const query = `
		SELECT EXISTS(SELECT 1 FROM links WHERE redirect_link = $1)
	`
	var isExists bool
	if err := p.conn.Pool.QueryRow(ctx, query, redirectLink).Scan(&isExists); err != nil {
		return false, err
	}

	return isExists, nil
}
