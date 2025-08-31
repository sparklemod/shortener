package domain

import "time"

type Link struct {
	ID           string
	OriginalLink string
	RedirectLink string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type CreateLink struct {
	OriginalLink string
	RedirectLink string
}
