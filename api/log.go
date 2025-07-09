package api

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	db "project_expiry_backend/db/sqlc"
)

type createLogRequest struct {
	ActionType int32   `json:"action_type" binding:"required,min=1"`
	Message    *string `json:"message"`
	Data       *string `json:"data"`
	UserID     int32   `json:"user_id" binding:"required,min=1"`
}

func (server *Server) createLog(c *gin.Context) {
	var req createLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateLogParams{
		ActionType: req.ActionType,
		UserID:     req.UserID,
	}

	if req.Message != nil {
		arg.Message = sql.NullString{String: *req.Message, Valid: true}
	}

	if req.Data != nil {
		arg.Data = sql.NullString{String: *req.Data, Valid: true}
	}

	log, err := server.store.CreateLog(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, log)
}

type getLogRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getLog(c *gin.Context) {
	var req getLogRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	log, err := server.store.GetLog(c, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, log)
}

type listLogsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=20"`
}

func (server *Server) listLogs(c *gin.Context) {
	var req listLogsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListLogsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	logs, err := server.store.ListLogs(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, logs)
}

type deleteLogRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteLog(c *gin.Context) {
	var req deleteLogRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := server.store.DeleteLog(c, req.ID); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, struct{}{})
}

type updateLogUriRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

type updateLogBodyRequest struct {
	ActionType int32   `json:"action_type" binding:"required,min=1"`
	Message    *string `json:"message"`
	Data       *string `json:"data"`
	UserID     int32   `json:"user_id" binding:"required,min=1"`
}

func (server *Server) updateLog(c *gin.Context) {
	var uriReq updateLogUriRequest
	if err := c.ShouldBindUri(&uriReq); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var bodyReq updateLogBodyRequest
	if err := c.ShouldBindBodyWith(&bodyReq, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.store.GetLog(c, uriReq.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "log not found"})
		} else {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	arg := db.UpdateLogParams{
		ID:         uriReq.ID,
		ActionType: bodyReq.ActionType,
		UserID:     bodyReq.UserID,
	}

	if bodyReq.Message != nil {
		arg.Message = sql.NullString{String: *bodyReq.Message, Valid: true}
	}

	if bodyReq.Data != nil {
		arg.Data = sql.NullString{String: *bodyReq.Data, Valid: true}
	}

	updatedLog, err := server.store.UpdateLog(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, updatedLog)
}