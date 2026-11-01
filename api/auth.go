package api

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
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
	Password string `json:"password" binding:"required,min=15,max=128"`
}

type userResponse struct {
	ID               int32      `json:"id"`
	Name             string     `json:"name"`
	Username         string     `json:"username"`
	Email            string     `json:"email"`
	Phone            *string    `json:"phone"`
	HasAvatar        bool       `json:"has_avatar"`
	EmailVerified    bool       `json:"email_verified"`
	PhoneVerified    bool       `json:"phone_verified"`
	TwoFactorEnabled bool       `json:"two_factor_enabled"`
	LastLogin        *time.Time `json:"last_login"`
}

func newUserResponse(user db.User) userResponse {
	var lastLogin *time.Time
	if user.LastLogin.Valid {
		lastLogin = &user.LastLogin.Time
	}

	var phone *string
	if user.Phone.Valid {
		phone = &user.Phone.String
	}

	return userResponse{
		ID:               user.ID,
		Name:             user.Name,
		Username:         user.Username,
		Email:            user.Email,
		Phone:            phone,
		HasAvatar:        user.HasAvatar,
		EmailVerified:    user.EmailVerified,
		PhoneVerified:    user.PhoneVerified,
		TwoFactorEnabled: user.TwoFactorEnabled,
		LastLogin:        lastLogin,
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
		Phone:       sql.NullString{},
		HasAvatar:   false,
		Password:    hashedPassword,
		RoleID:      1, // TODO: change this
		IsSuperuser: false,
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
	Password string `json:"password" binding:"required,max=128"`
}

type loginUserResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
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
			ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid email or password")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = utils.CheckPassword(req.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid email or password")))
		return
	}
	if user.DeletedAt.Valid || (user.LockedUntil.Valid && time.Now().Before(user.LockedUntil.Time)) {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("account is suspended")))
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

	refreshToken := db.CreateRefreshTokenParams{
		UserID:    user.ID,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(server.config.RefreshTokenDuration),
	}
	refreshTokenRes, err := server.store.CreateRefreshToken(ctx, refreshToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := loginUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenRes.ID.String(),
	}
	ctx.JSON(http.StatusOK, rsp)
}

type changePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=15,max=128"`
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

	// update password, delete all refresh tokens
	err = server.store.ExecTx(ctx, func(q *db.Queries) error {
		// save new pwd
		err := q.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
			ID:       user.ID,
			Password: hashedPassword,
		})
		if err != nil {
			return err
		}

		return q.HardDeleteRefreshTokensByUser(ctx, user.ID)
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "password changed successfully"})
}

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (server *Server) refreshToken(ctx *gin.Context) {
	var req refreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	tokenID, err := uuid.Parse(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var user db.User
	var newRefreshToken db.RefreshToken
	var accessToken string
	var errRefreshTokenExpired = errors.New("refresh token expired")

	err = server.store.ExecTx(ctx, func(q *db.Queries) error {
		oldToken, err := q.GetRefreshTokenByIDForUpdate(ctx, tokenID)
		if err != nil {
			return err
		}
		if !time.Now().Before(oldToken.ExpiresAt) {
			return errRefreshTokenExpired
		}

		user, err = q.GetUserByID(ctx, oldToken.UserID)
		if err != nil {
			return err
		}

		now := time.Now()
		if user.DeletedAt.Valid || (user.LockedUntil.Valid && now.Before(user.LockedUntil.Time)) {
			return errors.New("account is suspended")
		}

		accessToken, err = server.tokenMaker.CreateToken(
			user.Username,
			server.config.AccessTokenDuration,
		)
		if err != nil {
			return err
		}

		if err := q.HardDeleteRefreshToken(ctx, oldToken.ID); err != nil {
			return err
		}

		now = time.Now()
		newRefreshToken, err = q.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
			UserID:    oldToken.UserID,
			IssuedAt:  now,
			ExpiresAt: now.Add(server.config.RefreshTokenDuration),
		})
		return err
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, errRefreshTokenExpired) {
			ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid or expired refresh token")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := loginUserResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken.ID.String(),
	}
	ctx.JSON(http.StatusOK, rsp)
}

type logoutUserRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (server *Server) logoutUser(ctx *gin.Context) {
	var req logoutUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	tokenID, err := uuid.Parse(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.store.HardDeleteRefreshToken(ctx, tokenID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}
