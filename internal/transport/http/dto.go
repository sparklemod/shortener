package http

type RequestStatus struct {
	Status      int    `json:"status"`
	Description string `json:"description"`
}

type CreateLinkRequest struct {
	Name     string `json:"name"`
	Redirect string `json:"redirect,omitempty"`
}

type CreateLinkResponseData struct {
	ID string `json:"id"`
}

type CreateLinkResponse struct {
	RequestStatus RequestStatus          `json:"status"`
	Data          CreateLinkResponseData `json:"data"`
}
