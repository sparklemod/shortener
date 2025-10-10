package http

import (
	"errors"
	"log"
	"net/http"
	"shortener/internal/model"

	"github.com/gin-gonic/gin"
)

func (srv *HttpServer) ParseQuery(c *gin.Context, data any) error {
	if err := c.ShouldBindQuery(data); err != nil {
		c.JSON(http.StatusBadRequest,
			RequestStatus{
				Status:      http.StatusBadRequest,
				Description: "incorrect query params",
			})
		return err
	}
	if err := srv.validate.Struct(data); err != nil {
		c.JSON(http.StatusBadRequest,
			RequestStatus{
				Status:      http.StatusBadRequest,
				Description: err.Error(),
			})
		return err
	}
	return nil
}

func (srv *HttpServer) ParseData(c *gin.Context, data any) error {
	if err := c.ShouldBindJSON(data); err != nil {
		c.JSON(http.StatusBadRequest,
			RequestStatus{
				Status:      http.StatusBadRequest,
				Description: "incorrect body request",
			})
		return err
	}
	if err := srv.validate.Struct(data); err != nil {
		c.JSON(http.StatusBadRequest,
			RequestStatus{
				Status:      http.StatusBadRequest,
				Description: err.Error(),
			})
		return err
	}
	return nil
}

func (srv *HttpServer) GetLinksHandler(c *gin.Context) {
	var filter FilterLinksRequest
	if err := srv.ParseQuery(c, &filter); err != nil {
		return
	}

	links, err := srv.uc.FilterLinks(c, model.FilterLinksInput{
		IsActive:  filter.IsActive,
		SortBy:    filter.SortBy,
		SortOrder: filter.SortOrder,
		Limit:     filter.Limit,
		Offset:    filter.Offset,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			RequestStatus{Status: http.StatusInternalServerError, Description: err.Error()})
		return
	}

	c.JSON(http.StatusOK, links)
}

func (srv *HttpServer) CreateLinkHandler(c *gin.Context) {
	var reqData CreateLinkRequest
	if err := srv.ParseData(c, &reqData); err != nil {
		return
	}

	link, err := srv.uc.CreateLink(c, model.CreateLinkInput{
		OriginalUrl: reqData.Url,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			RequestStatus{
				Status:      http.StatusInternalServerError,
				Description: err.Error(),
			})
		return
	}

	c.JSON(http.StatusOK, CreateLinkResponse{
		RequestStatus: RequestStatus{
			Status:      http.StatusOK,
			Description: "ok",
		},
		Data: CreateLinkResponseData{
			ShortenedUrl: link.ShortenUrl,
		},
	})
}

func (srv *HttpServer) RedirectHandler(c *gin.Context) {
	shortenUrl := c.Param("shorten-url")

	redirectUrl, err := srv.uc.Redirect(c, shortenUrl)
	if err != nil {
		if errors.Is(err, model.ErrorNotFound) {
			c.JSON(http.StatusNotFound,
				RequestStatus{
					Status:      http.StatusNotFound,
					Description: err.Error(),
				})
			return
		}

		log.Printf("error getting redirect url for %q: %v", shortenUrl, err)
		c.JSON(http.StatusInternalServerError,
			RequestStatus{
				Status:      http.StatusInternalServerError,
				Description: err.Error(),
			})
		return
	}

	c.Redirect(http.StatusMovedPermanently, redirectUrl)
}
