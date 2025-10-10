package usecase

import (
	"context"
	"shortener/internal/model"
)

type Repository interface {
	Post(ctx context.Context, link model.Link) (*model.Link, error)
	Get(ctx context.Context, shortenUrl string) (*model.Link, error)
	IncrementVisits(ctx context.Context, shortenUrl string) error
	FilterLinks(ctx context.Context, filters model.FilterLinksInput) ([]model.Link, error)
}
