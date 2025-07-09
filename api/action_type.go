package api

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	db "project_expiry_backend/db/sqlc"
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

	actionType, err := server.store.CreateActionType(c, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, actionType)
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

	actionType, err := server.store.GetActionType(c, req.ID)
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
	PageID   int32 `form:"page_id" binding:"required,min=1"`
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

	actionTypes, err := server.store.ListActionTypes(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, actionTypes)
}

type deleteActionTypeRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteActionType(c *gin.Context) {
	var req deleteActionTypeRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := server.store.DeleteActionType(c, req.ID); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, struct{}{})
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

	currentActionType, err := server.store.GetActionType(c, uriReq.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "action type not found"})
		} else {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	if currentActionType.Name == bodyReq.Name {
		c.JSON(http.StatusBadRequest, gin.H{"error": "action type name cannot be the same"})
		return
	}

	arg := db.UpdateActionTypeParams{
		ID:   uriReq.ID,
		Name: bodyReq.Name,
	}

	actionType, err := server.store.UpdateActionType(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, actionType)
}
