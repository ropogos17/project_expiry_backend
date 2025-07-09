package api

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	db "project_expiry_backend/db/sqlc"
)

type createShareRequest struct {
	SourceUserID int32 `json:"source_user_id" binding:"required,min=1"`
	StorageID    int32 `json:"storage_id" binding:"required,min=1"`
	TargetUserID int32 `json:"target_user_id" binding:"required,min=1"`
	CanShare     bool  `json:"can_share"`
}

func (server *Server) createShare(c *gin.Context) {
	var req createShareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateShareParams{
		SourceUserID: req.SourceUserID,
		StorageID:    req.StorageID,
		TargetUserID: req.TargetUserID,
		CanShare:     req.CanShare,
	}

	share, err := server.store.CreateShare(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, share)
}

type getShareRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getShare(c *gin.Context) {
	var req getShareRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	share, err := server.store.GetShare(c, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, share)
}

type listSharesRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=20"`
}

func (server *Server) listShares(c *gin.Context) {
	var req listSharesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListSharesParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	shares, err := server.store.ListShares(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, shares)
}

type deleteShareRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteShare(c *gin.Context) {
	var req deleteShareRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := server.store.DeleteShare(c, req.ID); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, struct{}{})
}

type updateShareUriRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

type updateShareBodyRequest struct {
	SourceUserID int32 `json:"source_user_id" binding:"required,min=1"`
	StorageID    int32 `json:"storage_id" binding:"required,min=1"`
	TargetUserID int32 `json:"target_user_id" binding:"required,min=1"`
	CanShare     bool  `json:"can_share"`
}

func (server *Server) updateShare(c *gin.Context) {
	var uriReq updateShareUriRequest
	if err := c.ShouldBindUri(&uriReq); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var bodyReq updateShareBodyRequest
	if err := c.ShouldBindBodyWith(&bodyReq, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.store.GetShare(c, uriReq.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "share not found"})
		} else {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	arg := db.UpdateShareParams{
		ID:           uriReq.ID,
		SourceUserID: bodyReq.SourceUserID,
		StorageID:    bodyReq.StorageID,
		TargetUserID: bodyReq.TargetUserID,
		CanShare:     bodyReq.CanShare,
	}

	updatedShare, err := server.store.UpdateShare(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, updatedShare)
}