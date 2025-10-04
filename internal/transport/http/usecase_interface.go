package http

import (
	"context"
	"shortener/internal/model"
)

type UseCase interface {
	CreateLink(ctx context.Context, link model.CreateLinkRequest) (link2 *model.Link, err error)
	Redirect(ctx context.Context, shortenUrl string) (string, error)
}
