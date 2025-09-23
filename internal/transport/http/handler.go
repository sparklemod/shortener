package http

import (
	"net/http"

	"shortener/internal/usecase/dto"

	"github.com/gin-gonic/gin"
)

func (srv *HttpServer) ParseData(c *gin.Context, data any) error {
	if err := c.ShouldBindJSON(data); err != nil {
		c.JSON(http.StatusBadRequest, RequestStatus{Status: http.StatusBadRequest, Desc: "no correct body"})
		return err
	}
	if err := srv.validate.Struct(data); err != nil {
		c.JSON(http.StatusBadRequest, RequestStatus{Status: http.StatusBadRequest, Desc: err.Error()})
		return err
	}
	return nil
}

func (srv *HttpServer) userIDFromHeader(c *gin.Context) (string, bool) {
	// Prefer userID set by JWT middleware (subject claim)
	if uid := c.GetString("userID"); uid != "" {
		return uid, true
	}
	// Fallback for backward compatibility: X-User-ID header
	uid := c.GetHeader("X-User-ID")
	if uid == "" {
		c.JSON(http.StatusUnauthorized, RequestStatus{Status: http.StatusUnauthorized, Desc: "unauthorized: missing user id"})
		return "", false
	}
	return uid, true
}

func (srv *HttpServer) CreateLinkHandler(c *gin.Context) {
	var reqData CreateLinkRequest
	if err := srv.ParseData(c, &reqData); err != nil {
		return
	}
	//userID, ok := srv.userIDFromHeader(c)
	//if !ok {
	//	return
	//}
	id, err := srv.uc.CreateLink(c, dto.CreateLink{
		OriginalLink: reqData.Name,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, RequestStatus{Status: http.StatusInternalServerError, Desc: err.Error()})
		return
	}
	c.JSON(http.StatusOK, CreateLinkResponse{RequestStatus: RequestStatus{Status: http.StatusOK, Desc: "ok"}, Data: CreateLinkResponseData{ID: id}})
}

// TODO вынести в отд папки
type CreateLinkRequest struct {
	Name     string `json:"name"`
	Redirect string `json:"redirect,omitempty"`
}

type RequestStatus struct {
	Status int    `json:"status"`
	Desc   string `json:"desc"`
}

type CreateLinkResponseData struct {
	ID string `json:"id"`
}

type CreateLinkResponse struct {
	RequestStatus RequestStatus          `json:"status"`
	Data          CreateLinkResponseData `json:"data"`
}
