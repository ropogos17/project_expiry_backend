package api

import (
	"database/sql"
	"errors"
	"net/http"
	db "project_expiry_backend/db/sqlc"
	"project_expiry_backend/token"
	"strings"

	"github.com/gin-gonic/gin"
)

// -----------------------------------------------------------------------------
// Requests
// -----------------------------------------------------------------------------

type getUserRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

type listUsersRequest struct {
	PageID   int32 `form:"page_id" binding:"omitempty,min=1,max=1000000"`
	PageSize int32 `form:"page_size" binding:"omitempty,min=1,max=100"`
}

type updateUserRequest struct {
	Name      *string `json:"name" binding:"omitempty,min=1,max=100"`
	Username  *string `json:"username" binding:"omitempty,alphanum,max=50"`
	Email     *string `json:"email" binding:"omitempty,email,max=255"`
	Phone     *string `json:"phone" binding:"omitempty,max=30"`
	HasAvatar *bool   `json:"has_avatar"`

	// only admin
	EmailVerified    *bool  `json:"email_verified"`
	PhoneVerified    *bool  `json:"phone_verified"`
	TwoFactorEnabled *bool  `json:"two_factor_enabled"`
	RoleID           *int32 `json:"role_id" binding:"omitempty,min=1"`
	IsSuperuser      *bool  `json:"is_superuser"`
}

type deleteUserRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

// -----------------------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------------------

func getAuthorizationPayload(c *gin.Context) (*token.Payload, bool) {
	value, exists := c.Get(authorizationPayloadKey)
	if !exists {
		return nil, false
	}

	payload, ok := value.(*token.Payload)
	if !ok || payload == nil {
		return nil, false
	}

	return payload, true
}

func (server *Server) getAuthenticatedUser(c *gin.Context) (db.User, int, error) {
	payload, ok := getAuthorizationPayload(c)
	if !ok {
		return db.User{}, http.StatusUnauthorized, errors.New("unauthorized")
	}

	user, err := server.store.GetUserByUsername(c, payload.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.User{}, http.StatusUnauthorized, errors.New("user not found")
		}

		return db.User{}, http.StatusInternalServerError, errors.New("could not retrieve user")
	}

	// deleted user
	if user.DeletedAt.Valid {
		return db.User{}, http.StatusUnauthorized, errors.New("account unavailable")
	}

	return user, http.StatusOK, nil
}

func isAdmin(user db.User) bool {
	return user.IsSuperuser || user.RoleID == 2
}

func hasAdminFields(req updateUserRequest) bool {
	return req.EmailVerified != nil ||
		req.PhoneVerified != nil ||
		req.TwoFactorEnabled != nil ||
		req.RoleID != nil ||
		req.IsSuperuser != nil
}

func hasUpdateFields(req updateUserRequest) bool {
	return req.Name != nil ||
		req.Username != nil ||
		req.Email != nil ||
		req.Phone != nil ||
		req.HasAvatar != nil ||
		hasAdminFields(req)
}

// -----------------------------------------------------------------------------
// GET /users/:id
// -----------------------------------------------------------------------------

func (server *Server) getUser(c *gin.Context) {
	var req getUserRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authUser, status, err := server.getAuthenticatedUser(c)
	if err != nil {
		c.JSON(status, errorResponse(err))
		return
	}

	admin := isAdmin(authUser)

	// only admin
	if authUser.ID != req.ID && !admin {
		c.JSON(http.StatusForbidden, errorResponse(errors.New("access denied")))
		return
	}

	user, err := server.store.GetUserByID(c, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, errorResponse(errors.New("user not found")))
			return
		}

		c.JSON(http.StatusInternalServerError, errorResponse(errors.New("could not retrieve user")))
		return
	}

	// deleted users
	if user.DeletedAt.Valid && !admin {
		c.JSON(http.StatusNotFound, errorResponse(errors.New("user not found")))
		return
	}

	c.JSON(http.StatusOK, newUserResponse(user))
}

// -----------------------------------------------------------------------------
// GET /users
// -----------------------------------------------------------------------------

func (server *Server) listUsers(c *gin.Context) {
	authUser, status, err := server.getAuthenticatedUser(c)
	if err != nil {
		c.JSON(status, errorResponse(err))
		return
	}

	if !isAdmin(authUser) {
		c.JSON(http.StatusForbidden, errorResponse(errors.New("admin permission required")))
		return
	}

	var req listUsersRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// default paging
	if req.PageID == 0 {
		req.PageID = 1
	}

	if req.PageSize == 0 {
		req.PageSize = 20
	}

	offset := (req.PageID - 1) * req.PageSize

	users, err := server.store.ListUsers(c, db.ListUsersParams{
		Limit:  req.PageSize,
		Offset: offset,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(errors.New("could not retrieve users")))
		return
	}

	responses := make([]userResponse, 0, len(users))

	for _, user := range users {
		responses = append(responses, newUserResponse(user))
	}

	c.JSON(http.StatusOK, gin.H{
		"users":     responses,
		"page_id":   req.PageID,
		"page_size": req.PageSize,
	})
}

// -----------------------------------------------------------------------------
// PUT /users/:id
// -----------------------------------------------------------------------------

func (server *Server) updateUser(c *gin.Context) {
	var uriReq getUserRequest

	if err := c.ShouldBindUri(&uriReq); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req updateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !hasUpdateFields(req) {
		c.JSON(http.StatusBadRequest, errorResponse(errors.New("no fields provided for update")))
		return
	}

	authUser, status, err := server.getAuthenticatedUser(c)
	if err != nil {
		c.JSON(status, errorResponse(err))
		return
	}

	admin := isAdmin(authUser)

	// only admin
	if authUser.ID != uriReq.ID && !admin {
		c.JSON(http.StatusForbidden, errorResponse(errors.New("access denied")))
		return
	}

	if !admin && hasAdminFields(req) {
		c.JSON(http.StatusForbidden, errorResponse(errors.New("admin permission required")))
		return
	}

	user, err := server.store.GetUserByID(c, uriReq.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, errorResponse(errors.New("user not found")))
			return
		}

		c.JSON(http.StatusInternalServerError, errorResponse(errors.New("could not retrieve user")))
		return
	}

	if user.DeletedAt.Valid {
		c.JSON(http.StatusNotFound, errorResponse(errors.New("user not found")))
		return
	}

	params := db.UpdateUserParams{
		ID:                user.ID,
		Name:              user.Name,
		Username:          user.Username,
		Email:             user.Email,
		Phone:             user.Phone,
		HasAvatar:         user.HasAvatar,
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
	}

	if req.Name != nil {
		params.Name = *req.Name
	}

	if req.Username != nil {
		params.Username = *req.Username
	}

	if req.Email != nil {
		if params.Email != *req.Email {
			params.EmailVerified = false
		}

		params.Email = *req.Email
	}

	if req.Phone != nil {
		phone := strings.TrimSpace(*req.Phone)

		if phone == "" {
			params.Phone = sql.NullString{}
			params.PhoneVerified = false
		} else {
			if !params.Phone.Valid || params.Phone.String != phone {
				params.PhoneVerified = false
			}

			params.Phone = sql.NullString{
				String: phone,
				Valid:  true,
			}
		}
	}

	if req.HasAvatar != nil {
		params.HasAvatar = *req.HasAvatar
	}

	if admin {
		if req.EmailVerified != nil {
			params.EmailVerified = *req.EmailVerified
		}

		if req.PhoneVerified != nil {
			if *req.PhoneVerified && !params.Phone.Valid {
				c.JSON(http.StatusBadRequest, errorResponse(errors.New("cannot verify missing phone number")))
				return
			}

			params.PhoneVerified = *req.PhoneVerified
		}

		if req.TwoFactorEnabled != nil {
			if *req.TwoFactorEnabled && !params.TwoFactorSecret.Valid {
				c.JSON(http.StatusBadRequest, errorResponse(errors.New("two-factor secret is not configured")))
				return
			}

			params.TwoFactorEnabled = *req.TwoFactorEnabled
		}

		if req.RoleID != nil {
			params.RoleID = *req.RoleID
		}

		if req.IsSuperuser != nil {
			params.IsSuperuser = *req.IsSuperuser
		}
	}

	updatedUser, err := server.store.UpdateUser(c, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, errorResponse(errors.New("user not found")))
			return
		}

		c.JSON(http.StatusInternalServerError, errorResponse(errors.New("could not update user")))
		return
	}

	c.JSON(http.StatusOK, newUserResponse(updatedUser))
}

// -----------------------------------------------------------------------------
// DELETE /users/:id
// -----------------------------------------------------------------------------

func (server *Server) deleteUser(c *gin.Context) {
	var req deleteUserRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authUser, status, err := server.getAuthenticatedUser(c)
	if err != nil {
		c.JSON(status, errorResponse(err))
		return
	}

	admin := isAdmin(authUser)
	if authUser.ID != req.ID && !admin {
		c.JSON(http.StatusForbidden, errorResponse(errors.New("access denied")))
		return
	}

	user, err := server.store.GetUserByID(c, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, errorResponse(errors.New("user not found")))
			return
		}

		c.JSON(http.StatusInternalServerError, errorResponse(errors.New("could not retrieve user")))
		return
	}

	if user.DeletedAt.Valid {
		c.Status(http.StatusNoContent)
		return
	}

	// delete refresh tokens
	err = server.store.ExecTx(c, func(q *db.Queries) error {
		if err := q.SoftDeleteUser(c, req.ID); err != nil {
			return err
		}

		_, err := q.HardDeleteRefreshTokensByUser(c, req.ID)
		return err
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(errors.New("could not delete user")))
		return
	}

	c.Status(http.StatusNoContent)
}
