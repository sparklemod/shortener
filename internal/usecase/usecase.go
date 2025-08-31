package usecase

import (
	"context"

	"shortener/internal/repository/domain"
	"shortener/internal/usecase/dto"
	"shortener/utils"
)

type Usecase struct {
	repo Repository
}

func New(repo Repository) *Usecase {
	return &Usecase{repo: repo}
}

func (uc *Usecase) CreateLink(ctx context.Context, in dto.CreateLink) (string, error) {
	code, err := utils.GenerateBase58String(8)
	if err != nil {
		return "", err
	}
	link := domain.CreateLink{
		OriginalLink: in.OriginalLink,
		RedirectLink: code,
	}

	id, err := uc.repo.CreateLink(ctx, link)
	if err != nil {
		return "", err
	}
	return id, nil
}
