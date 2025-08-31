package usecase

import (
	"context"
	"shortener/internal/repository/domain"
)

type Repository interface {
	CreateLink(ctx context.Context, link domain.CreateLink) (string, error)
}
