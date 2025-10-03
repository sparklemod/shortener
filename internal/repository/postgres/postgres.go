package postgres

import (
	"context"
	"errors"
	adapterpg "shortener/internal/adapter/postgres"
	"shortener/internal/model"
	"time"

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
	//link.CreatedAt = time.Now().UTC()

	const query = `
		INSERT INTO links (original_url, redirect_url) 
		VALUES ($1, $2)
		ON CONFLICT (redirect_url) DO NOTHING
		RETURNING *;
	`
	var result model.Link
	err := pgxscan.Get(ctx, p.conn.Pool, &result, query, link.OriginalUrl, link.RedirectUrl)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.ErrorNonUniq
	}

	if err != nil {
		return nil, err
	}

	return &result, nil
}

type pgShortening struct {
	Identifier  int       `bson:"_id"`
	OriginalURL string    `bson:"original_url"`
	Visits      int       `bson:"visits"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
}

func pgShorteningFromModel(shortening model.Link) pgShortening {
	return pgShortening{
		Identifier:  shortening.Id,
		OriginalURL: shortening.OriginalUrl,
		Visits:      shortening.Visits,
		CreatedAt:   shortening.CreatedAt,
		UpdatedAt:   shortening.UpdatedAt,
	}
}

func modelShorteningFromPg(shortening pgShortening) *model.Link {
	return &model.Link{
		Id:          shortening.Identifier,
		OriginalUrl: shortening.OriginalURL,
		Visits:      shortening.Visits,
		CreatedAt:   shortening.CreatedAt,
		UpdatedAt:   shortening.UpdatedAt,
	}
}
