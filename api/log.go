package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	db "project_expiry_backend/db/sqlc"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sqlc-dev/pqtype"
)

// Helper functions

func trimStringPtr(s *string) *string {
	if s != nil {
		if trimmed := strings.TrimSpace(*s); trimmed != "" {
			return &trimmed
		}
	}
	return nil
}

func valueOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// APIs
type createLogRequest struct {
	ActionTypeID int32           `json:"action_type_id" binding:"required"`
	Message      *string         `json:"message"`
	Data         json.RawMessage `json:"data"`
	UserID       int32           `json:"user_id" binding:"required"`
}

func (server *Server) createLog(c *gin.Context) {
	var req createLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	req.Message = trimStringPtr(req.Message)
	params := db.CreateLogParams{
		ActionTypeID: req.ActionTypeID,
		Message:      nullableString(valueOrEmpty(req.Message)),
		Data: pqtype.NullRawMessage{
			RawMessage: req.Data,
			Valid:      req.Data != nil,
		},
		UserID: req.UserID,
	}

	log, err := server.store.CreateLog(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusCreated, log)
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

	log, err := server.store.GetLogByID(c.Request.Context(), req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
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
	PageSize int32 `form:"page_size" binding:"required,min=1,max=100"`
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

	logs, err := server.store.ListLogs(c.Request.Context(), arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if logs == nil {
		logs = []db.Log{}
	}
	c.JSON(http.StatusOK, logs)
}

type deleteLogRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) hardDeleteLog(c *gin.Context) {
	var req deleteLogRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	rows, err := server.store.HardDeleteLog(c.Request.Context(), req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "log not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "log permanently deleted"})
}

type updateLogUriRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}
type updateLogBodyRequest struct {
	ActionTypeID int32           `json:"action_type_id" binding:"required,min=1"`
	Message      *string         `json:"message"`
	Data         json.RawMessage `json:"data"`
}

func (server *Server) updateLog(c *gin.Context) {
	var uriReq updateLogUriRequest
	if err := c.ShouldBindUri(&uriReq); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var bodyReq updateLogBodyRequest
	if err := c.ShouldBindJSON(&bodyReq); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	bodyReq.Message = trimStringPtr(bodyReq.Message)
	arg := db.UpdateLogParams{
		ID:           uriReq.ID,
		ActionTypeID: bodyReq.ActionTypeID,
		Message:      nullableString(valueOrEmpty(bodyReq.Message)),
		Data: pqtype.NullRawMessage{
			RawMessage: bodyReq.Data,
			Valid:      bodyReq.Data != nil,
		},
	}

	updatedLog, err := server.store.UpdateLog(c.Request.Context(), arg)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "log not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, updatedLog)

}
