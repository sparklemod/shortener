package usecase

import (
	"context"
	"errors"

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
	var code string
	var err error
	attempts := 3

	for {
		if attempts == 0 {
			return "", errors.New("не удалось создать ссылку")
		}

		code, err = utils.GenerateBase58String(8)
		if err != nil {
			return "", err
		}

		isExist, err := uc.repo.CheckLinkIfExist(ctx, code)
		if err != nil {
			return "", err
		}

		if !isExist {
			break
		}

		attempts--
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
