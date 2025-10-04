package usecase

import (
	"context"
	"errors"
	"shortener/internal/model"
	"shortener/utils"
)

type Usecase struct {
	repo Repository
}

func New(repo Repository) *Usecase {
	return &Usecase{repo: repo}
}

func (uc *Usecase) CreateLink(ctx context.Context, in model.CreateLinkRequest) (*model.Link, error) {
	for attempts := 3; attempts > 0; attempts-- {
		shortenUrl, err := utils.GenerateShortLink(8)
		if err != nil {
			return nil, err
		}

		inputLink := model.Link{
			OriginalUrl: in.OriginalUrl,
			ShortenUrl:  shortenUrl,
		}

		link, err := uc.repo.Post(ctx, inputLink)
		if err == nil {
			return link, nil
		}

		if errors.Is(err, model.ErrorNonUniq) {
			continue
		}

		return nil, err
	}

	return nil, model.ErrorAttemptsExhausted
}
