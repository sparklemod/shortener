package model

import (
	"errors"
	"time"
)

var (
	ErrorNotFound          = errors.New("not found")
	ErrorNonUniq           = errors.New("shorten link already exists")
	ErrorAttemptsExhausted = errors.New("failed to create uniq link")
	ErrorIncrementVisits   = errors.New("failed to increment visits")
)

type Link struct {
	Id          int       `db:"id"`
	OriginalUrl string    `db:"original_url"`
	ShortenUrl  string    `db:"shorten_url"`
	IsActive    bool      `db:"is_active"`
	Visits      int       `db:"visits"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type CreateLinkInput struct {
	OriginalUrl string
}

type FilterLinksInput struct {
	IsActive  *bool `db:"is_active"`
	SortBy    string
	SortOrder string
	Limit     int
	Offset    int
}
