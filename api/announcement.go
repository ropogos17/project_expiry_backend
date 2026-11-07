package api

import (
	"database/sql"
	"errors"
	"net/http"
	"net/url"
	db "project_expiry_backend/db/sqlc"
	"strings"

	"github.com/gin-gonic/gin"
)

// Request structs for announcements

type createAnnouncementRequest struct {
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description"`
	ExternalLink string `json:"external_link"`
}

type getAnnouncementRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

type listAnnouncementsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1,max=1000000"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=20"`
}

type updateAnnouncementUriRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

type updateAnnouncementBodyRequest struct {
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description"`
	ExternalLink string `json:"external_link"`
}

type deleteAnnouncementRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

type restoreAnnouncementRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func nullableString(value string) sql.NullString {
	value = strings.TrimSpace(value)

	return sql.NullString{
		String: value,
		Valid:  value != "",
	}
}

func validateExternalLink(value string) error {
	value = strings.TrimSpace(value)

	if value == "" {
		return nil
	}

	parsed, err := url.ParseRequestURI(value)
	if err != nil {
		return errors.New("invalid external link")
	}

	scheme := strings.ToLower(parsed.Scheme)

	if (scheme != "http" && scheme != "https") || parsed.Host == "" {
		return errors.New("external link must be a valid HTTP or HTTPS URL")
	}

	return nil
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

	announcements, err := server.store.ListAnnouncements(c.Request.Context(), arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if announcements == nil {
		announcements = []db.Announcement{}
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

	announcement, err := server.store.GetAnnouncementByID(c.Request.Context(), req.ID)
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

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "announcement name cannot be empty"})
		return
	}
	req.ExternalLink = strings.TrimSpace(req.ExternalLink)
	if err := validateExternalLink(req.ExternalLink); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authUser, status, err := server.getAuthenticatedUser(c)
	if err != nil {
		c.JSON(status, errorResponse(err))
		return
	}

	arg := db.CreateAnnouncementParams{
		Name:         req.Name,
		Description:  nullableString(req.Description),
		ExternalLink: nullableString(req.ExternalLink),
		CreatedBy:    authUser.ID,
	}

	announcement, err := server.store.CreateAnnouncement(c.Request.Context(), arg)
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

	bodyReq.Name = strings.TrimSpace(bodyReq.Name)
	if bodyReq.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "announcement name cannot be empty"})
	}
	bodyReq.ExternalLink = strings.TrimSpace(bodyReq.ExternalLink)
	if err := validateExternalLink(bodyReq.ExternalLink); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	arg := db.UpdateAnnouncementParams{
		ID:           uriReq.ID,
		Name:         bodyReq.Name,
		Description:  nullableString(bodyReq.Description),
		ExternalLink: nullableString(bodyReq.ExternalLink),
	}

	updatedAnnouncement, err := server.store.UpdateAnnouncement(c.Request.Context(), arg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "announcement not found"})
			return
		}
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

	err := server.store.DeleteAnnouncement(c.Request.Context(), req.ID)
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

	err := server.store.RestoreAnnouncement(c.Request.Context(), req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "announcement restored successfully"})
}

// hardDeleteAnnouncement handles DELETE /announcements/:id/hard - requires admin role (permanent delete)
func (server *Server) hardDeleteAnnouncement(c *gin.Context) {
	var req deleteAnnouncementRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	rows, err := server.store.HardDeleteAnnouncement(c.Request.Context(), req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "announcement not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "announcement permanently deleted"})
}
