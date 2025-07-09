package api

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	db "project_expiry_backend/db/sqlc"
)

type createUserStatusRequest struct {
	Name string `json:"name" binding:"required"`
}

func (server *Server) createUserStatus(c *gin.Context) {
	var req createUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userStatus, err := server.store.CreateStatus(c, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, userStatus)
}

type getUserStatusRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getUserStatus(c *gin.Context) {
	var req getUserStatusRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userStatus, err := server.store.GetStatus(c, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, userStatus)
}

type listUserStatusesRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=20"`
}

func (server *Server) listUserStatuses(c *gin.Context) {
	var req listUserStatusesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListStatusesParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	userStatuses, err := server.store.ListStatuses(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, userStatuses)
}

type deleteUserStatusRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteUserStatus(c *gin.Context) {
	var req deleteUserStatusRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := server.store.DeleteStatus(c, req.ID); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, struct{}{})
}

type updateUserStatusUriRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

type updateUserStatusBodyRequest struct {
	Name string `json:"name" binding:"required"`
}

func (server *Server) updateUserStatus(c *gin.Context) {
	var uriReq updateUserStatusUriRequest
	if err := c.ShouldBindUri(&uriReq); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var bodyReq updateUserStatusBodyRequest
	if err := c.ShouldBindBodyWith(&bodyReq, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	currentUserStatus, err := server.store.GetStatus(c, uriReq.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user status not found"})
		} else {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	if currentUserStatus.Name == bodyReq.Name {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user status name cannot be the same"})
		return
	}

	arg := db.UpdateStatusParams{
		ID:   uriReq.ID,
		Name: bodyReq.Name,
	}

	updatedUserStatus, err := server.store.UpdateStatus(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, updatedUserStatus)
}