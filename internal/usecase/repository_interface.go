package usecase

import (
	"context"
	"shortener/internal/model"
)

type Repository interface {
	Post(ctx context.Context, link model.Link) (*model.Link, error)
}
