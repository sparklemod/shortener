package http

import (
	"net/http"

	"shortener/internal/usecase/dto"

	"github.com/gin-gonic/gin"
)

func (srv *HttpServer) ParseData(c *gin.Context, data any) error {
	if err := c.ShouldBindJSON(data); err != nil {
		c.JSON(http.StatusBadRequest, RequestStatus{Status: http.StatusBadRequest, Description: "no correct body"})
		return err
	}
	if err := srv.validate.Struct(data); err != nil {
		c.JSON(http.StatusBadRequest, RequestStatus{Status: http.StatusBadRequest, Description: err.Error()})
		return err
	}
	return nil
}

func (srv *HttpServer) CreateLinkHandler(c *gin.Context) {
	var reqData CreateLinkRequest
	if err := srv.ParseData(c, &reqData); err != nil {
		return
	}

	id, err := srv.uc.CreateLink(c, dto.CreateLink{
		OriginalLink: reqData.Name,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			RequestStatus{Status: http.StatusInternalServerError, Description: err.Error()})
		return
	}

	c.JSON(http.StatusOK, CreateLinkResponse{
		RequestStatus: RequestStatus{Status: http.StatusOK, Description: "ok"},
		Data:          CreateLinkResponseData{ID: id},
	})
}
