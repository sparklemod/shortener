package usecase

import (
	"context"
	"shortener/internal/repository/domain"
)

type Repository interface {
	CreateLink(ctx context.Context, link domain.CreateLink) (string, error)
	CheckLinkIfExist(ctx context.Context, RedirectLink string) (isExists bool, err error)
}
