package api

import (
	"errors"
	"fmt"
	"net/http"
	db "project_expiry_backend/db/sqlc"
	"slices"
	"strings"
	"time"

	"project_expiry_backend/token"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"

	adminRoleID  = 2
	editorRoleID = 3
)

// AuthMiddleware creates a gin middleware for authorization
func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fields := strings.Fields(ctx.GetHeader(authorizationHeaderKey))
		if len(fields) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		if len(fields) != 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		if !strings.EqualFold(fields[0], authorizationTypeBearer) {
			err := fmt.Errorf("unsupported authorization type")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		payload, err := tokenMaker.VerifyToken(fields[1])
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}

// requireAdminOrEditor allows superusers, admins and editors.
func requireAdminOrEditor(store *db.Store) gin.HandlerFunc {
	return requireRole(store, []int32{2, 3})
}

// requireAdmin allows only superusers and admins.
func requireAdmin(store *db.Store) gin.HandlerFunc {
	return requireRole(store, []int32{2})
}

// Role-based access control middleware
func requireRole(store *db.Store, roles []int32) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, exists := c.Get(authorizationPayloadKey)
		payload, ok := val.(*token.Payload)
		if !exists || !ok || payload == nil {
			abortWithError(c, http.StatusUnauthorized, errors.New("unauthorized"))
			return
		}

		// Get user to check role and superuser status
		user, err := store.GetUserByUsername(c, payload.Username)
		if err != nil {
			abortWithError(c, http.StatusUnauthorized, errors.New("user not found"))
			return
		}

		now := time.Now()
		if user.DeletedAt.Valid || (user.LockedUntil.Valid && now.Before(user.LockedUntil.Time)) {
			abortWithError(c, http.StatusUnauthorized, errors.New("account is suspended"))
			return
		}

		hasRequiredRole := slices.Contains(roles, user.RoleID)
		if !user.IsSuperuser && !hasRequiredRole {
			abortWithError(c, http.StatusUnauthorized, errors.New("insufficient permissions: role does not match"))
			return
		}

		c.Next()
	}
}

// abortWithError aborts the request and returns a JSON error response.
func abortWithError(c *gin.Context, status int, err error) {
	c.AbortWithStatusJSON(status, errorResponse(err))
}
