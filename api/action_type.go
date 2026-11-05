package api

import (
	"database/sql"
	"errors"
	"net/http"
	db "project_expiry_backend/db/sqlc"
	"strings"

	"github.com/gin-gonic/gin"
)

type createActionTypeRequest struct {
	Name string `json:"name" binding:"required"`
}

func (server *Server) createActionType(c *gin.Context) {
	var req createActionTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "action type name cannot be empty"})
		return
	}
	actionType, err := server.store.CreateActionType(c.Request.Context(), req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusCreated, actionType)
}

type getActionTypeRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getActionType(c *gin.Context) {
	var req getActionTypeRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	actionType, err := server.store.GetActionTypeByID(c.Request.Context(), req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, actionType)
}

type listActionTypesRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1,max=1000000"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=20"`
}

func (server *Server) listActionTypes(c *gin.Context) {
	var req listActionTypesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListActionTypesParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	actionTypes, err := server.store.ListActionTypes(c.Request.Context(), arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if actionTypes == nil {
		actionTypes = []db.ActionType{}
	}
	c.JSON(http.StatusOK, actionTypes)
}

type deleteActionTypeRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) hardDeleteActionType(c *gin.Context) {
	var req deleteActionTypeRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	rows, err := server.store.HardDeleteActionType(c.Request.Context(), req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "action type not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "action type permanently deleted"})
}

type updateActionTypeUriRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}
type updateActionTypeBodyRequest struct {
	Name string `json:"name" binding:"required"`
}

func (server *Server) updateActionType(c *gin.Context) {
	var uriReq updateActionTypeUriRequest
	if err := c.ShouldBindUri(&uriReq); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var bodyReq updateActionTypeBodyRequest
	if err := c.ShouldBindJSON(&bodyReq); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	bodyReq.Name = strings.TrimSpace(bodyReq.Name)
	if bodyReq.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "announcement name cannot be empty"})
		return
	}

	arg := db.UpdateActionTypeParams{
		ID:   uriReq.ID,
		Name: bodyReq.Name,
	}

	updatedActionType, err := server.store.UpdateActionType(c.Request.Context(), arg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "action type not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, updatedActionType)
}

type actionTypeExistsByNameRequest struct {
	Name string `uri:"name" binding:"required,min=1"`
}

func (server *Server) actionTypeExistsByName(c *gin.Context) {
	var req actionTypeExistsByNameRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	exists, err := server.store.ActionTypeExistsByName(c.Request.Context(), req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, exists)
}
