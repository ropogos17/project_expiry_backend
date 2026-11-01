package api

import (
	"errors"
	"fmt"
	"net/http"
	db "project_expiry_backend/db/sqlc"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"project_expiry_backend/token"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

// AuthMiddleware creates a gin middleware for authorization
func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) != 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}

// Role-based access control middleware
func requireAdminOrEditor(store *db.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		//authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
		authPayload, exists := c.Get(authorizationPayloadKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("unauthorized")))
			return
		}
		payload, ok := authPayload.(*token.Payload)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("unauthorized")))
			return
		}

		// Get user to check role and superuser status
		user, err := store.GetUserByUsername(c, payload.Username)
		if err != nil {
			c.JSON(http.StatusUnauthorized, errorResponse(errors.New("user not found")))
			c.Abort()
			return
		}

		now := time.Now()
		if user.DeletedAt.Valid || (user.LockedUntil.Valid && now.Before(user.LockedUntil.Time)) {
			c.JSON(http.StatusUnauthorized, errorResponse(errors.New("account is suspended")))
			c.Abort()
			return
		}

		// Check if user is superuser or has admin/editor role (assuming role IDs 1=admin, 2=editor)
		if !user.IsSuperuser && user.RoleID != 2 && user.RoleID != 3 {
			c.JSON(http.StatusForbidden, errorResponse(errors.New("insufficient permissions: admin or editor role required")))
			c.Abort()
			return
		}

		c.Next()
	}
}
