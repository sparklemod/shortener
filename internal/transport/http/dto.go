package http

type RequestStatus struct {
	Status      int    `json:"status"`
	Description string `json:"description"`
}

type CreateLinkRequest struct {
	Url string `json:"url" validate:"required,url"`
}

type CreateLinkResponseData struct {
	ShortenedUrl string `json:"shortened_url,omitempty"`
}

type CreateLinkResponse struct {
	RequestStatus RequestStatus          `json:"status"`
	Data          CreateLinkResponseData `json:"data"`
}

type FilterLinksRequest struct {
	IsActive   *bool   `form:"is_active"`
	ShortenUrl *string `form:"shorten_url"`
	SortBy     string  `form:"sort_by"`
	SortOrder  string  `form:"sort_order"`
	Limit      int     `form:"limit"`
	Offset     int     `form:"offset"`
}
