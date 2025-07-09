package api

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	db "project_expiry_backend/db/sqlc"
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

	role, err := server.store.CreateRole(c, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, role)
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

	role, err := server.store.GetRole(c, req.ID)
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
	PageSize int32 `form:"page_size" binding:"required,min=5,max=20"`
}

func (server *Server) listRoles(c *gin.Context) {
	var req listRolesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListRolesParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	roles, err := server.store.ListRoles(c, arg)
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

	if err := server.store.DeleteRole(c, req.ID); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, struct{}{})
}

type updateRoleUriRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}
type updateRoleBodyRequest struct {
	Name string `json:"name" binding:"required"`
}

func (server *Server) updateRole(c *gin.Context) {
	var uriReq updateRoleUriRequest
	if err := c.ShouldBindUri(&uriReq); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var bodyReq updateRoleBodyRequest
	if err := c.ShouldBindBodyWith(&bodyReq, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	currentRole, err := server.store.GetRole(c, uriReq.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
		} else {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	if currentRole.Name == bodyReq.Name {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role name cannot be the same"})
		return
	}

	arg := db.UpdateRoleParams{
		ID:   uriReq.ID,
		Name: bodyReq.Name,
	}

	updatedRole, err := server.store.UpdateRole(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, updatedRole)
}
