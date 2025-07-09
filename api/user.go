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

type createUserRequest struct {
	Name        string `json:"name" binding:"required"`
	Username    string `json:"username" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Phone       string `json:"phone" binding:"required"`
	HasAvatar   bool   `json:"has_avatar"`
	Password    string `json:"password" binding:"required,min=6"`
	Salt        string `json:"salt" binding:"required"`
	RoleID      int32  `json:"role_id" binding:"required,min=1"`
	IsSuperuser bool   `json:"is_superuser"`
	StatusID    int32  `json:"status_id" binding:"required,min=1"`
}

func (server *Server) createUser(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Name:        req.Name,
		Username:    req.Username,
		Email:       req.Email,
		Phone:       req.Phone,
		HasAvatar:   req.HasAvatar,
		Password:    req.Password,
		Salt:        req.Salt,
		RoleID:      req.RoleID,
		IsSuperuser: req.IsSuperuser,
		StatusID:    req.StatusID,
	}

	user, err := server.store.CreateUser(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, user)
}

type getUserRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getUser(c *gin.Context) {
	var req getUserRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(c, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, user)
}

type listUsersRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=20"`
}

func (server *Server) listUsers(c *gin.Context) {
	var req listUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListUsersParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	users, err := server.store.ListUsers(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, users)
}

type deleteUserRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteUser(c *gin.Context) {
	var req deleteUserRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := server.store.DeleteUser(c, req.ID); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, struct{}{})
}

type updateUserUriRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

type updateUserBodyRequest struct {
	Name              string  `json:"name" binding:"required"`
	Username          string  `json:"username" binding:"required"`
	Email             string  `json:"email" binding:"required,email"`
	Phone             string  `json:"phone" binding:"required"`
	HasAvatar         bool    `json:"has_avatar"`
	Password          string  `json:"password" binding:"required,min=6"`
	Salt              string  `json:"salt" binding:"required"`
	EmailVerified     bool    `json:"email_verified"`
	PhoneVerified     bool    `json:"phone_verified"`
	TwoFactorEnabled  bool    `json:"two_factor_enabled"`
	TwoFactorSecret   *string `json:"two_factor_secret"`
	LoginAttempts     int32   `json:"login_attempts"`
	PasswordChangedAt time.Time `json:"password_changed_at" binding:"required"`
	RoleID            int32   `json:"role_id" binding:"required,min=1"`
	IsSuperuser       bool    `json:"is_superuser"`
	StatusID          int32   `json:"status_id" binding:"required,min=1"`
}

func (server *Server) updateUser(c *gin.Context) {
	var uriReq updateUserUriRequest
	if err := c.ShouldBindUri(&uriReq); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var bodyReq updateUserBodyRequest
	if err := c.ShouldBindBodyWith(&bodyReq, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.store.GetUser(c, uriReq.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		} else {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	arg := db.UpdateUserParams{
		ID:                uriReq.ID,
		Name:              bodyReq.Name,
		Username:          bodyReq.Username,
		Email:             bodyReq.Email,
		Phone:             bodyReq.Phone,
		HasAvatar:         bodyReq.HasAvatar,
		Password:          bodyReq.Password,
		Salt:              bodyReq.Salt,
		EmailVerified:     bodyReq.EmailVerified,
		PhoneVerified:     bodyReq.PhoneVerified,
		TwoFactorEnabled:  bodyReq.TwoFactorEnabled,
		LoginAttempts:     bodyReq.LoginAttempts,
		PasswordChangedAt: bodyReq.PasswordChangedAt,
		RoleID:            bodyReq.RoleID,
		IsSuperuser:       bodyReq.IsSuperuser,
		StatusID:          bodyReq.StatusID,
		UpdatedAt:         time.Now(),
	}

	if bodyReq.TwoFactorSecret != nil {
		arg.TwoFactorSecret = sql.NullString{String: *bodyReq.TwoFactorSecret, Valid: true}
	}

	updatedUser, err := server.store.UpdateUser(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, updatedUser)
}