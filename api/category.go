package api

import (
	"database/sql"
	"errors"
	"net/http"
	db "project_expiry_backend/db/sqlc"
	"strings"

	"github.com/gin-gonic/gin"
)

type createCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

func (server *Server) createCategory(c *gin.Context) {
	var req createCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category name cannot be empty"})
		return
	}
	category, err := server.store.CreateCategory(c.Request.Context(), req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, category)
}

type getCategoryRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getCategory(c *gin.Context) {
	var req getCategoryRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	category, err := server.store.GetCategoryByID(c.Request.Context(), req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, category)
}

type listCategoriesRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1,max=1000000"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=20"`
}

func (server *Server) listCategories(c *gin.Context) {
	var req listCategoriesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListCategoriesParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	categories, err := server.store.ListCategories(c.Request.Context(), arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if categories == nil {
		categories = []db.Category{}
	}
	c.JSON(http.StatusOK, categories)
}

type deleteCategoryRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) hardDeleteCategory(c *gin.Context) {
	var req deleteCategoryRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	rows, err := server.store.HardDeleteCategory(c.Request.Context(), req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "category permanently deleted"})
}

type updateCategoryUriRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}
type updateCategoryBodyRequest struct {
	Name string `json:"name" binding:"required"`
}

func (server *Server) updateCategory(c *gin.Context) {
	var uriReq updateCategoryUriRequest
	if err := c.ShouldBindUri(&uriReq); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var bodyReq updateCategoryBodyRequest
	if err := c.ShouldBindJSON(&bodyReq); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	bodyReq.Name = strings.TrimSpace(bodyReq.Name)
	if bodyReq.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category name cannot be empty"})
		return
	}

	arg := db.UpdateCategoryParams{
		ID:   uriReq.ID,
		Name: bodyReq.Name,
	}

	updatedCategory, err := server.store.UpdateCategory(c.Request.Context(), arg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, updatedCategory)
}

type categoryExistsByNameRequest struct {
	Name string `uri:"name" binding:"required,min=1"`
}

func (server *Server) categoryExistsByName(c *gin.Context) {
	var req categoryExistsByNameRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	exists, err := server.store.CategoryExistsByName(c.Request.Context(), req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, exists)
}
