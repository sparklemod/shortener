package http

import (
	"context"
	"shortener/internal/model"
)

type UseCase interface {
	CreateLink(ctx context.Context, link model.CreateLinkInput) (link2 *model.Link, err error)
	Redirect(ctx context.Context, shortenUrl string) (string, error)
	FilterLinks(ctx context.Context, filters model.FilterLinksInput) ([]model.Link, error)
}
