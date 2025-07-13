package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "project_expiry_backend/db/sqlc"
	"project_expiry_backend/token"
	"project_expiry_backend/utils"
)

type createUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required,alphanum"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type userResponse struct {
	ID                int32      `json:"id"`
	Name              string     `json:"name"`
	Username          string     `json:"username"`
	Email             string     `json:"email"`
	Phone             string     `json:"phone"`
	HasAvatar         bool       `json:"has_avatar"`
	EmailVerified     bool       `json:"email_verified"`
	PhoneVerified     bool       `json:"phone_verified"`
	TwoFactorEnabled  bool       `json:"two_factor_enabled"`
	LastLogin         *time.Time `json:"last_login"`
	LoginAttempts     int32      `json:"login_attempts"`
	LockedUntil       *time.Time `json:"locked_until"`
	PasswordChangedAt time.Time  `json:"password_changed_at"`
	RoleID            int32      `json:"role_id"`
	IsSuperuser       bool       `json:"is_superuser"`
	StatusID          int32      `json:"status_id"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at"`
}

func newUserResponse(user db.User) userResponse {
	var lastLogin *time.Time
	if user.LastLogin.Valid {
		lastLogin = &user.LastLogin.Time
	}

	var lockedUntil *time.Time
	if user.LockedUntil.Valid {
		lockedUntil = &user.LockedUntil.Time
	}

	var deletedAt *time.Time
	if user.DeletedAt.Valid {
		deletedAt = &user.DeletedAt.Time
	}

	return userResponse{
		ID:                user.ID,
		Name:              user.Name,
		Username:          user.Username,
		Email:             user.Email,
		Phone:             user.Phone,
		HasAvatar:         user.HasAvatar,
		EmailVerified:     user.EmailVerified,
		PhoneVerified:     user.PhoneVerified,
		TwoFactorEnabled:  user.TwoFactorEnabled,
		LastLogin:         lastLogin,
		LoginAttempts:     user.LoginAttempts,
		LockedUntil:       lockedUntil,
		PasswordChangedAt: user.PasswordChangedAt,
		RoleID:            user.RoleID,
		IsSuperuser:       user.IsSuperuser,
		StatusID:          user.StatusID,
		CreatedAt:         user.CreatedAt,
		UpdatedAt:         user.UpdatedAt,
		DeletedAt:         deletedAt,
	}
}

func (server *Server) registerUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Name:        req.Name,
		Username:    req.Username,
		Email:       req.Email,
		Phone:       "",
		HasAvatar:   false,
		Password:    hashedPassword,
		RoleID:      0, // TODO: change this
		IsSuperuser: false,
		StatusID:    0, // TODO: change this
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newUserResponse(user)
	ctx.JSON(http.StatusOK, rsp)
}

type loginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("invalid email")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = utils.CheckPassword(req.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid password")))
		return
	}

	accessToken, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := loginUserResponse{
		AccessToken: accessToken,
		User:        newUserResponse(user),
	}
	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) getCurrentUser(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	user, err := server.store.GetUserByUsername(ctx, authPayload.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("user not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newUserResponse(user)
	ctx.JSON(http.StatusOK, rsp)
}

type updateProfileRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Phone string `json:"phone" binding:"required"`
}

func (server *Server) updateProfile(ctx *gin.Context) {
	var req updateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	user, err := server.store.GetUserByUsername(ctx, authPayload.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("user not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.UpdateUserParams{
		ID:                user.ID,
		Name:              req.Name,
		Username:          user.Username,
		Email:             req.Email,
		Phone:             req.Phone,
		HasAvatar:         user.HasAvatar,
		Password:          user.Password,
		EmailVerified:     user.EmailVerified,
		PhoneVerified:     user.PhoneVerified,
		TwoFactorEnabled:  user.TwoFactorEnabled,
		TwoFactorSecret:   user.TwoFactorSecret,
		LastLogin:         user.LastLogin,
		LoginAttempts:     user.LoginAttempts,
		LockedUntil:       user.LockedUntil,
		PasswordChangedAt: user.PasswordChangedAt,
		RoleID:            user.RoleID,
		IsSuperuser:       user.IsSuperuser,
		StatusID:          user.StatusID,
		DeletedAt:         user.DeletedAt,
		UpdatedAt:         time.Now(),
	}

	updatedUser, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newUserResponse(updatedUser)
	ctx.JSON(http.StatusOK, rsp)
}

type changePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

func (server *Server) changePassword(ctx *gin.Context) {
	var req changePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	user, err := server.store.GetUserByUsername(ctx, authPayload.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("user not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = utils.CheckPassword(req.OldPassword, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid old password")))
		return
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	now := time.Now()
	arg := db.UpdateUserParams{
		ID:                user.ID,
		Name:              user.Name,
		Username:          user.Username,
		Email:             user.Email,
		Phone:             user.Phone,
		HasAvatar:         user.HasAvatar,
		Password:          hashedPassword,
		EmailVerified:     user.EmailVerified,
		PhoneVerified:     user.PhoneVerified,
		TwoFactorEnabled:  user.TwoFactorEnabled,
		TwoFactorSecret:   user.TwoFactorSecret,
		LastLogin:         user.LastLogin,
		LoginAttempts:     user.LoginAttempts,
		LockedUntil:       user.LockedUntil,
		PasswordChangedAt: now,
		RoleID:            user.RoleID,
		IsSuperuser:       user.IsSuperuser,
		StatusID:          user.StatusID,
		DeletedAt:         user.DeletedAt,
		UpdatedAt:         now,
	}

	_, err = server.store.UpdateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "password changed successfully"})
}

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type refreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}

func (server *Server) refreshToken(ctx *gin.Context) {
	var req refreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, err := server.tokenMaker.CreateToken(
		payload.Username,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := refreshTokenResponse{
		AccessToken: accessToken,
	}
	ctx.JSON(http.StatusOK, rsp)
}
