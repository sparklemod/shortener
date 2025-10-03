package model

import (
	"errors"
	"time"
)

var (
	ErrorNotFound          = errors.New("not found")
	ErrorNonUniq           = errors.New("redirect link already exists")
	ErrorAttemptsExhausted = errors.New("failed to create uniq link")
)

type Link struct {
	Id          int       `db:"id"`
	OriginalUrl string    `db:"original_url"`
	RedirectUrl string    `db:"redirect_url"`
	IsActive    bool      `db:"is_active"`
	Visits      int       `db:"visits"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type CreateLinkRequest struct {
	OriginalUrl string
}
