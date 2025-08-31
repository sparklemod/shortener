package postgres

import (
	"context"
	adapterpg "shortener/internal/adapter/postgres"
)

type Postgres struct {
	conn *adapterpg.Postgres
}

func NewPostgres(conn *adapterpg.Postgres) *Postgres {
	return &Postgres{conn: conn}
}

func (p *Postgres) CreateLink(ctx context.Context, originalLink string) error {

}
