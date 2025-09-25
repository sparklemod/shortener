package http

import (
	"context"
	"shortener/internal/usecase/dto"
)

type UseCase interface {
	CreateLink(ctx context.Context, link dto.CreateLink) (id string, err error)
}
