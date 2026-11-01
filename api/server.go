package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	db "project_expiry_backend/db/sqlc"
	"project_expiry_backend/token"
	"project_expiry_backend/utils"
)

type Server struct {
	config     utils.Config
	store      *db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config utils.Config, store *db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("Cannot create token maker: %w", err)
	}
	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}

	server.SetupRouter()
	return server, nil
}

func (server *Server) SetupRouter() {
	router := gin.Default()
	authMiddleware := authMiddleware(server.tokenMaker)

	// Auth routes (public)
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", server.registerUser)
		authRoutes.POST("/login", server.loginUser)
		authRoutes.POST("/refresh", server.refreshToken)
		authRoutes.POST("/logout", server.logoutUser)
		pw := router.Group("/auth").Use(authMiddleware)
		{
			pw.PUT("/change-password", server.changePassword)
		}
	}

	// Protected routes
	protectedRoutes := router.Group("/").Use(authMiddleware)
	{
		// User routes
		protectedRoutes.GET("/users/:id", server.getUser)
		protectedRoutes.PUT("/users/:id", server.updateUser)
		protectedRoutes.DELETE("/users/:id", server.deleteUser)

		// Role routes
		/*protectedRoutes.POST("/roles", server.createRole)
		protectedRoutes.GET("/roles/:id", server.getRole)
		protectedRoutes.GET("/roles", server.listRoles)
		protectedRoutes.DELETE("/roles/:id", server.deleteRole)
		protectedRoutes.PUT("/roles/:id", server.updateRole)

		protectedRoutes.GET("/announcements", server.listAnnouncements)
		protectedRoutes.GET("/announcements/:id", server.getAnnouncement)*/
	}

	// Admin/Editor only routes for announcements
	adminRoutes := router.Group("/").Use(authMiddleware, requireAdminOrEditor(server.store))
	{
		// User routes
		adminRoutes.GET("/users", server.listUsers)

		/*adminRoutes.POST("/announcements", server.createAnnouncement)
		adminRoutes.PUT("/announcements/:id", server.updateAnnouncement)
		adminRoutes.DELETE("/announcements/:id", server.deleteAnnouncement)
		adminRoutes.POST("/announcements/:id/restore", server.restoreAnnouncement)
		adminRoutes.DELETE("/announcements/:id/hard", server.hardDeleteAnnouncement)*/
	}

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
