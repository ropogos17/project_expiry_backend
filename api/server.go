package api

import (
	"fmt"
	db "project_expiry_backend/db/sqlc"
	"project_expiry_backend/token"
	"project_expiry_backend/utils"

	"github.com/gin-gonic/gin"
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
		return nil, fmt.Errorf("cannot create token maker: %w", err)
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
	if err := router.SetTrustedProxies([]string{"172.16.0.0/12"}); err != nil {
		panic(err)
	}
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

		// Announcements
		protectedRoutes.GET("/announcements", server.listAnnouncements)
		protectedRoutes.GET("/announcements/:id", server.getAnnouncement)
	}

	// Admin/Editor only routes
	adminOrEditorRoutes := router.Group("/").Use(authMiddleware, requireAdminOrEditor(server.store))
	{
		// User
		adminOrEditorRoutes.GET("/users", server.listUsers)

		// Announcements
		adminOrEditorRoutes.POST("/announcements", server.createAnnouncement)
		adminOrEditorRoutes.PUT("/announcements/:id", server.updateAnnouncement)
		adminOrEditorRoutes.DELETE("/announcements/:id", server.deleteAnnouncement)
		adminOrEditorRoutes.POST("/announcements/:id/restore", server.restoreAnnouncement)
	}

	// Admin only routes
	adminRoutes := router.Group("/").Use(authMiddleware, requireAdmin(server.store))
	{
		// Roles
		adminRoutes.POST("/roles", server.createRole)
		adminRoutes.GET("/roles/:id", server.getRole)
		adminRoutes.GET("/roles", server.listRoles)
		adminRoutes.DELETE("/roles/:id", server.deleteRole)
		adminRoutes.PUT("/roles/:id", server.updateRole)

		// Announcements
		adminRoutes.DELETE("/announcements/:id/hard", server.hardDeleteAnnouncement)

		// Action types
		adminRoutes.POST("/action-types", server.createActionType)
		adminRoutes.GET("/action-types", server.listActionTypes)
		adminRoutes.GET("/action-types/:id", server.getActionType)
		adminRoutes.PUT("/action-types/:id", server.updateActionType)
		adminRoutes.DELETE("/action-types/:id/hard", server.hardDeleteActionType)
		adminRoutes.GET("/action-types/exists/:name", server.actionTypeExistsByName)

		// Categories
		adminRoutes.POST("/categories", server.createCategory)
		adminRoutes.GET("/categories", server.listCategories)
		adminRoutes.GET("/categories/:id", server.getCategory)
		adminRoutes.PUT("/categories/:id", server.updateCategory)
		adminRoutes.DELETE("/categories/:id/hard", server.hardDeleteCategory)
		adminRoutes.GET("/categories/exists/:name", server.categoryExistsByName)

		// Logs
		adminRoutes.POST("/logs", server.createLog)
		adminRoutes.GET("/logs", server.listLogs)
		adminRoutes.GET("/logs/:id", server.getLog)
		adminRoutes.PUT("/logs/:id", server.updateLog)
		adminRoutes.DELETE("/logs/:id/hard", server.hardDeleteLog)
	}

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
