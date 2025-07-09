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

type createShoppingListItemRequest struct {
	Name         string  `json:"name" binding:"required"`
	Description  *string `json:"description"`
	Quantity     *string `json:"quantity"`
	IsChecked    bool    `json:"is_checked"`
	IsDeleted    bool    `json:"is_deleted"`
	LastModified int32   `json:"last_modified" binding:"required,min=1"`
}

func (server *Server) createShoppingListItem(c *gin.Context) {
	var req createShoppingListItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateShoppingListItemParams{
		Name:         req.Name,
		IsChecked:    req.IsChecked,
		IsDeleted:    req.IsDeleted,
		UpdatedAt:    time.Now(),
		LastModified: req.LastModified,
	}

	if req.Description != nil {
		arg.Description = sql.NullString{String: *req.Description, Valid: true}
	}

	if req.Quantity != nil {
		arg.Quantity = sql.NullString{String: *req.Quantity, Valid: true}
	}

	shoppingListItem, err := server.store.CreateShoppingListItem(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, shoppingListItem)
}

type getShoppingListItemRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getShoppingListItem(c *gin.Context) {
	var req getShoppingListItemRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	shoppingListItem, err := server.store.GetShoppingListItem(c, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, shoppingListItem)
}

type listShoppingListItemsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=20"`
}

func (server *Server) listShoppingListItems(c *gin.Context) {
	var req listShoppingListItemsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListShoppingListItemsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	shoppingListItems, err := server.store.ListShoppingListItems(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, shoppingListItems)
}

type deleteShoppingListItemRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteShoppingListItem(c *gin.Context) {
	var req deleteShoppingListItemRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := server.store.DeleteShoppingListItem(c, req.ID); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, struct{}{})
}

type updateShoppingListItemUriRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

type updateShoppingListItemBodyRequest struct {
	Name         string  `json:"name" binding:"required"`
	Description  *string `json:"description"`
	Quantity     *string `json:"quantity"`
	IsChecked    bool    `json:"is_checked"`
	IsDeleted    bool    `json:"is_deleted"`
	LastModified int32   `json:"last_modified" binding:"required,min=1"`
}

func (server *Server) updateShoppingListItem(c *gin.Context) {
	var uriReq updateShoppingListItemUriRequest
	if err := c.ShouldBindUri(&uriReq); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var bodyReq updateShoppingListItemBodyRequest
	if err := c.ShouldBindBodyWith(&bodyReq, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.store.GetShoppingListItem(c, uriReq.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "shopping list item not found"})
		} else {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	arg := db.UpdateShoppingListItemParams{
		ID:           uriReq.ID,
		Name:         bodyReq.Name,
		IsChecked:    bodyReq.IsChecked,
		IsDeleted:    bodyReq.IsDeleted,
		UpdatedAt:    time.Now(),
		LastModified: bodyReq.LastModified,
	}

	if bodyReq.Description != nil {
		arg.Description = sql.NullString{String: *bodyReq.Description, Valid: true}
	}

	if bodyReq.Quantity != nil {
		arg.Quantity = sql.NullString{String: *bodyReq.Quantity, Valid: true}
	}

	updatedShoppingListItem, err := server.store.UpdateShoppingListItem(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, updatedShoppingListItem)
}