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
