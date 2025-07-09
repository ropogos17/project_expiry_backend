package api

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"time"
	db "project_expiry_backend/db/sqlc"
)

type createItemRequest struct {
	Name         string    `json:"name" binding:"required"`
	CategoryID   int32     `json:"category_id" binding:"required,min=1"`
	Description  *string   `json:"description"`
	Quantity     *string   `json:"quantity"`
	HasImage     bool      `json:"has_image"`
	ExpiresAt    time.Time `json:"expires_at" binding:"required"`
	StorageID    int32     `json:"storage_id" binding:"required,min=1"`
	Shelf        *int32    `json:"shelf"`
	LastModified int32     `json:"last_modified" binding:"required,min=1"`
}

func (server *Server) createItem(c *gin.Context) {
	var req createItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateItemParams{
		Name:         req.Name,
		CategoryID:   req.CategoryID,
		HasImage:     req.HasImage,
		ExpiresAt:    req.ExpiresAt,
		StorageID:    req.StorageID,
		LastModified: req.LastModified,
	}

	if req.Description != nil {
		arg.Description = sql.NullString{String: *req.Description, Valid: true}
	}

	if req.Quantity != nil {
		arg.Quantity = sql.NullString{String: *req.Quantity, Valid: true}
	}

	if req.Shelf != nil {
		arg.Shelf = sql.NullInt32{Int32: *req.Shelf, Valid: true}
	}

	item, err := server.store.CreateItem(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, item)
}

type getItemRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getItem(c *gin.Context) {
	var req getItemRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	item, err := server.store.GetItem(c, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, item)
}

type listItemsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=20"`
}

func (server *Server) listItems(c *gin.Context) {
	var req listItemsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListItemsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	items, err := server.store.ListItems(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, items)
}

type deleteItemRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteItem(c *gin.Context) {
	var req deleteItemRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := server.store.DeleteItem(c, req.ID); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, struct{}{})
}

type updateItemUriRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

type updateItemBodyRequest struct {
	Name         string    `json:"name" binding:"required"`
	CategoryID   int32     `json:"category_id" binding:"required,min=1"`
	Description  *string   `json:"description"`
	Quantity     *string   `json:"quantity"`
	HasImage     bool      `json:"has_image"`
	IsDeleted    bool      `json:"is_deleted"`
	ExpiresAt    time.Time `json:"expires_at" binding:"required"`
	StorageID    int32     `json:"storage_id" binding:"required,min=1"`
	Shelf        *int32    `json:"shelf"`
	LastModified int32     `json:"last_modified" binding:"required,min=1"`
}

func (server *Server) updateItem(c *gin.Context) {
	var uriReq updateItemUriRequest
	if err := c.ShouldBindUri(&uriReq); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var bodyReq updateItemBodyRequest
	if err := c.ShouldBindBodyWith(&bodyReq, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.store.GetItem(c, uriReq.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		} else {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	arg := db.UpdateItemParams{
		ID:           uriReq.ID,
		Name:         bodyReq.Name,
		CategoryID:   bodyReq.CategoryID,
		HasImage:     bodyReq.HasImage,
		IsDeleted:    bodyReq.IsDeleted,
		UpdatedAt:    time.Now(),
		ExpiresAt:    bodyReq.ExpiresAt,
		StorageID:    bodyReq.StorageID,
		LastModified: bodyReq.LastModified,
	}

	if bodyReq.Description != nil {
		arg.Description = sql.NullString{String: *bodyReq.Description, Valid: true}
	}

	if bodyReq.Quantity != nil {
		arg.Quantity = sql.NullString{String: *bodyReq.Quantity, Valid: true}
	}

	if bodyReq.Shelf != nil {
		arg.Shelf = sql.NullInt32{Int32: *bodyReq.Shelf, Valid: true}
	}

	updatedItem, err := server.store.UpdateItem(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, updatedItem)
}