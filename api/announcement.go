package api

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	db "project_expiry_backend/db/sqlc"
)

type createAnnouncementRequest struct {
	Name         string  `json:"name" binding:"required"`
	Description  *string `json:"description"`
	ExternalLink *string `json:"external_link"`
	HasRead      bool    `json:"has_read"`
}

func (server *Server) createAnnouncement(c *gin.Context) {
	var req createAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateAnnouncementParams{
		Name:    req.Name,
		HasRead: req.HasRead,
	}

	if req.Description != nil {
		arg.Description = sql.NullString{String: *req.Description, Valid: true}
	}

	if req.ExternalLink != nil {
		arg.ExternalLink = sql.NullString{String: *req.ExternalLink, Valid: true}
	}

	announcement, err := server.store.CreateAnnouncement(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, announcement)
}

type getAnnouncementRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAnnouncement(c *gin.Context) {
	var req getAnnouncementRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	announcement, err := server.store.GetAnnouncement(c, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, announcement)
}

type listAnnouncementsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=20"`
}

func (server *Server) listAnnouncements(c *gin.Context) {
	var req listAnnouncementsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListAnnouncementsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	announcements, err := server.store.ListAnnouncements(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, announcements)
}

type deleteAnnouncementRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteAnnouncement(c *gin.Context) {
	var req deleteAnnouncementRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := server.store.DeleteAnnouncement(c, req.ID); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, struct{}{})
}

type updateAnnouncementUriRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

type updateAnnouncementBodyRequest struct {
	Name         string  `json:"name" binding:"required"`
	Description  *string `json:"description"`
	ExternalLink *string `json:"external_link"`
	HasRead      bool    `json:"has_read"`
}

func (server *Server) updateAnnouncement(c *gin.Context) {
	var uriReq updateAnnouncementUriRequest
	if err := c.ShouldBindUri(&uriReq); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var bodyReq updateAnnouncementBodyRequest
	if err := c.ShouldBindBodyWith(&bodyReq, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.store.GetAnnouncement(c, uriReq.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "announcement not found"})
		} else {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	arg := db.UpdateAnnouncementParams{
		ID:      uriReq.ID,
		Name:    bodyReq.Name,
		HasRead: bodyReq.HasRead,
	}

	if bodyReq.Description != nil {
		arg.Description = sql.NullString{String: *bodyReq.Description, Valid: true}
	}

	if bodyReq.ExternalLink != nil {
		arg.ExternalLink = sql.NullString{String: *bodyReq.ExternalLink, Valid: true}
	}

	updatedAnnouncement, err := server.store.UpdateAnnouncement(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, updatedAnnouncement)
}