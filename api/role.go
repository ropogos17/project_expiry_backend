package api

import (
	"database/sql"
	"errors"
	"net/http"
	db "project_expiry_backend/db/sqlc"
	"strings"

	"github.com/gin-gonic/gin"
)

type createRoleRequest struct {
	Name string `json:"name" binding:"required"`
}

func (server *Server) createRole(c *gin.Context) {
	var req createRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role cannot be empty"})
	}
	role, err := server.store.CreateRole(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusCreated, role)
}

type getRoleRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getRole(c *gin.Context) {
	var req getRoleRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	role, err := server.store.GetRoleByID(c.Request.Context(), req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, role)
}

type listRolesRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1,max=20"`
}

func (server *Server) listRoles(c *gin.Context) {
	var req listRolesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListRolesParams{ // TODO: túlcsordulás
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	roles, err := server.store.ListRoles(c.Request.Context(), arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, roles)
}

type deleteRoleRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteRole(c *gin.Context) {
	var req deleteRoleRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	rows, err := server.store.HardDeleteRole(c.Request.Context(), req.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "role does not exist"})
		return
	}
	c.Status(http.StatusNoContent)
}

type updateRoleRequest struct {
	ID   int32  `uri:"id" binding:"required,min=1"`
	Name string `json:"name" binding:"required"`
}

func (server *Server) updateRole(c *gin.Context) {
	var req updateRoleRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role cannot be empty"})
		return
	}
	name = strings.ToLower(name)
	arg := db.UpdateRoleParams{
		ID:   req.ID,
		Name: name,
	}

	updatedRole, err := server.store.UpdateRole(c.Request.Context(), arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, updatedRole)
}
