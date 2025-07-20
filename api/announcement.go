package api

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	db "project_expiry_backend/db/sqlc"
	"project_expiry_backend/token"
)

// Request structs for announcements

type createAnnouncementRequest struct {
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description"`
	ExternalLink string `json:"external_link"`
	HasRead      bool   `json:"has_read"`
}

type getAnnouncementRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

type listAnnouncementsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=20"`
}

type updateAnnouncementUriRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

type updateAnnouncementBodyRequest struct {
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description"`
	ExternalLink string `json:"external_link"`
	HasRead      bool   `json:"has_read"`
}

type deleteAnnouncementRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

type restoreAnnouncementRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

// Role-based access control middleware
func requireAdminOrEditor() gin.HandlerFunc {
	return func(c *gin.Context) {
		authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)

		// Get user to check role and superuser status
		user, err := c.Keys["server"].(*Server).store.GetUserByUsername(c, authPayload.Username)
		if err != nil {
			c.JSON(http.StatusUnauthorized, errorResponse(errors.New("user not found")))
			c.Abort()
			return
		}

		// Check if user is superuser or has admin/editor role (assuming role IDs 1=admin, 2=editor)
		if !user.IsSuperuser && user.RoleID != 1 && user.RoleID != 2 {
			c.JSON(http.StatusForbidden, errorResponse(errors.New("insufficient permissions: admin or editor role required")))
			c.Abort()
			return
		}

		c.Next()
	}
}

// User endpoints (list and get announcements)

// listAnnouncements handles GET /announcements - accessible to all authenticated users
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

// getAnnouncement handles GET /announcements/:id - accessible to all authenticated users
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

	// Check if announcement is soft deleted (users should not see deleted announcements)
	if announcement.DeletedAt.Valid {
		c.JSON(http.StatusNotFound, errorResponse(errors.New("announcement not found")))
		return
	}

	c.JSON(http.StatusOK, announcement)
}

// Admin/Editor endpoints (create, update, delete, restore)

// createAnnouncement handles POST /announcements - requires admin/editor role
func (server *Server) createAnnouncement(c *gin.Context) {
	var req createAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateAnnouncementParams{
		Name: req.Name,
		Description: sql.NullString{
			String: req.Description,
			Valid:  req.Description != "",
		},
		ExternalLink: sql.NullString{
			String: req.ExternalLink,
			Valid:  req.ExternalLink != "",
		},
		HasRead: req.HasRead,
	}

	announcement, err := server.store.CreateAnnouncement(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, announcement)
}

// updateAnnouncement handles PUT /announcements/:id - requires admin/editor role
func (server *Server) updateAnnouncement(c *gin.Context) {
	var uriReq updateAnnouncementUriRequest
	if err := c.ShouldBindUri(&uriReq); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var bodyReq updateAnnouncementBodyRequest
	if err := c.ShouldBindJSON(&bodyReq); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Check if announcement exists and is not deleted
	currentAnnouncement, err := server.store.GetAnnouncement(c, uriReq.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "announcement not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if currentAnnouncement.DeletedAt.Valid {
		c.JSON(http.StatusNotFound, gin.H{"error": "announcement not found"})
		return
	}

	arg := db.UpdateAnnouncementParams{
		ID:   uriReq.ID,
		Name: bodyReq.Name,
		Description: sql.NullString{
			String: bodyReq.Description,
			Valid:  bodyReq.Description != "",
		},
		ExternalLink: sql.NullString{
			String: bodyReq.ExternalLink,
			Valid:  bodyReq.ExternalLink != "",
		},
		HasRead: bodyReq.HasRead,
	}

	updatedAnnouncement, err := server.store.UpdateAnnouncement(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, updatedAnnouncement)
}

// deleteAnnouncement handles DELETE /announcements/:id - requires admin/editor role (soft delete)
func (server *Server) deleteAnnouncement(c *gin.Context) {
	var req deleteAnnouncementRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Check if announcement exists
	_, err := server.store.GetAnnouncement(c, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "announcement not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.store.DeleteAnnouncement(c, req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "announcement deleted successfully"})
}

// restoreAnnouncement handles POST /announcements/:id/restore - requires admin/editor role
func (server *Server) restoreAnnouncement(c *gin.Context) {
	var req restoreAnnouncementRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.RestoreAnnouncement(c, req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "announcement restored successfully"})
}

// hardDeleteAnnouncement handles DELETE /announcements/:id/hard - requires admin/editor role (permanent delete)
func (server *Server) hardDeleteAnnouncement(c *gin.Context) {
	var req deleteAnnouncementRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.HardDeleteAnnouncement(c, req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "announcement permanently deleted"})
}
