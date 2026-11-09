package api

/*
import (
	"database/sql"
	"net/http"
	db "project_expiry_backend/db/sqlc"
	"strings"

	"github.com/gin-gonic/gin"
)

type CreateStorageRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
	Location    *string `json:"location"`
	ShelfCount  int32   `json:"shelf_count" binding:"required,min=1,max=1000000"`
	HasImage    bool    `json:"has_image"`
}

func (server *Server) createStorage(c *gin.Context) {
	var req CreateStorageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "storage name cannot be empty"})
		return
	}
	req.Description = trimStringPtr(req.Description)
	req.Location = trimStringPtr(req.Location)

	authUser, status, err := server.getAuthenticatedUser(c)
	if err != nil {
		c.JSON(status, errorResponse(err))
		return
	}

	arg := db.CreateStorageParams{
		Name:        req.Name,
		Description: nullableString(valueOrEmpty(req.Description)),
		Location:    nullableString(valueOrEmpty(req.Location)),
		ShelfCount:  req.ShelfCount,
		HasImage:    req.HasImage,
		Owner:       authUser.ID,
	}

	storage, err := server.store.CreateStorage(c.Request.Context(), arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusCreated, storage)
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

	storage, err := server.store.GetStorageByID(c.Request.Context(), req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "storage not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, storage)
}

type listStorageRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1,max=100"`
}

func (server *Server) listStorage(c *gin.Context) {
	var req listStorageRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListStoragesParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	storageList, err := server.store.ListStorages(c.Request.Context(), arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if storageList == nil {
		storageList = []db.Storage{}
	}
	c.JSON(http.StatusOK, storageList)
}

type listStorageByOwnerRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1,max=100"`
	OwnerID  int32 `form:"owner" binding:"required,min=1,max=100000"`
}

func (server *Server) listStorageByOwner(c *gin.Context) {
	var req listStorageByOwnerRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListStoragesByOwnerParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
		Owner:  req.OwnerID,
	}

	storageList, err := server.store.ListStoragesByOwner(c.Request.Context(), arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if storageList == nil {
		storageList = []db.Storage{}
	}
	c.JSON(http.StatusOK, storageList)
}

type deleteStorageRequest struct {
	ID int32 `uri:"id" binding:"required,min=1,max=100000"`
}

func (server *Server) deleteStorage(c *gin.Context) {
	var req deleteStorageRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteStorage(c.Request.Context(), req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "storage deleted successfully"})
}

type restoreStorageRequest struct {
	ID int32 `uri:"id" binding:"required,min=1,max=100000"`
}

func (server *Server) restoreStorage(c *gin.Context) {
	var req restoreStorageRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.RestoreStorage(c.Request.Context(), req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "storage restored successfully"})
}

func (server *Server) hardDeleteStorage(c *gin.Context) {
	var req deleteStorageRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	rows, err := server.store.HardDeleteStorage(c.Request.Context(), req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "storage not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "storage permanently deleted"})
}
*/
