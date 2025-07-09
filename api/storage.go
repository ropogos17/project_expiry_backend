package api

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	db "project_expiry_backend/db/sqlc"
)

type createStorageRequest struct {
	Name        string `json:"name" binding:"required"`
	Description *string `json:"description"`
	Location    *string `json:"location"`
	Shelves     *int32  `json:"shelves"`
	HasImage    bool    `json:"has_image"`
	Owner       int32   `json:"owner" binding:"required,min=1"`
}

func (server *Server) createStorage(c *gin.Context) {
	var req createStorageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateStorageParams{
		Name:     req.Name,
		HasImage: req.HasImage,
		Owner:    req.Owner,
	}

	if req.Description != nil {
		arg.Description = sql.NullString{String: *req.Description, Valid: true}
	}

	if req.Location != nil {
		arg.Location = sql.NullString{String: *req.Location, Valid: true}
	}

	if req.Shelves != nil {
		arg.Shelves = sql.NullInt32{Int32: *req.Shelves, Valid: true}
	}

	storage, err := server.store.CreateStorage(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, storage)
}

type getStorageRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getStorage(c *gin.Context) {
	var req getStorageRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	storage, err := server.store.GetStorage(c, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, storage)
}

type listStoragesRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=20"`
}

func (server *Server) listStorages(c *gin.Context) {
	var req listStoragesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListStoragesParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	storages, err := server.store.ListStorages(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, storages)
}

type deleteStorageRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteStorage(c *gin.Context) {
	var req deleteStorageRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := server.store.DeleteStorage(c, req.ID); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, struct{}{})
}

type updateStorageUriRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

type updateStorageBodyRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
	Location    *string `json:"location"`
	Shelves     *int32  `json:"shelves"`
	HasImage    bool    `json:"has_image"`
	Owner       int32   `json:"owner" binding:"required,min=1"`
}

func (server *Server) updateStorage(c *gin.Context) {
	var uriReq updateStorageUriRequest
	if err := c.ShouldBindUri(&uriReq); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var bodyReq updateStorageBodyRequest
	if err := c.ShouldBindBodyWith(&bodyReq, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.store.GetStorage(c, uriReq.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "storage not found"})
		} else {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	arg := db.UpdateStorageParams{
		ID:       uriReq.ID,
		Name:     bodyReq.Name,
		HasImage: bodyReq.HasImage,
		Owner:    bodyReq.Owner,
	}

	if bodyReq.Description != nil {
		arg.Description = sql.NullString{String: *bodyReq.Description, Valid: true}
	}

	if bodyReq.Location != nil {
		arg.Location = sql.NullString{String: *bodyReq.Location, Valid: true}
	}

	if bodyReq.Shelves != nil {
		arg.Shelves = sql.NullInt32{Int32: *bodyReq.Shelves, Valid: true}
	}

	updatedStorage, err := server.store.UpdateStorage(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, updatedStorage)
}