package api

import (
	"github.com/gin-gonic/gin"
	db "project_expiry_backend/db/sqlc"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// Role routes
	router.POST("/roles", server.createRole)
	router.GET("/roles/:id", server.getRole)
	router.GET("/roles", server.listRoles)
	router.DELETE("/roles/:id", server.deleteRole)
	router.PUT("/roles/:id", server.updateRole)

	// ActionType routes
	router.POST("/action-types", server.createActionType)
	router.GET("/action-types/:id", server.getActionType)
	router.GET("/action-types", server.listActionTypes)
	router.DELETE("/action-types/:id", server.deleteActionType)
	router.PUT("/action-types/:id", server.updateActionType)

	// Announcement routes
	router.POST("/announcements", server.createAnnouncement)
	router.GET("/announcements/:id", server.getAnnouncement)
	router.GET("/announcements", server.listAnnouncements)
	router.DELETE("/announcements/:id", server.deleteAnnouncement)
	router.PUT("/announcements/:id", server.updateAnnouncement)

	// Category routes
	router.POST("/categories", server.createCategory)
	router.GET("/categories/:id", server.getCategory)
	router.GET("/categories", server.listCategories)
	router.DELETE("/categories/:id", server.deleteCategory)
	router.PUT("/categories/:id", server.updateCategory)

	// Item routes
	router.POST("/items", server.createItem)
	router.GET("/items/:id", server.getItem)
	router.GET("/items", server.listItems)
	router.DELETE("/items/:id", server.deleteItem)
	router.PUT("/items/:id", server.updateItem)

	// Language routes
	router.POST("/languages", server.createLanguage)
	router.GET("/languages/:id", server.getLanguage)
	router.GET("/languages", server.listLanguages)
	router.DELETE("/languages/:id", server.deleteLanguage)
	router.PUT("/languages/:id", server.updateLanguage)

	// Log routes
	router.POST("/logs", server.createLog)
	router.GET("/logs/:id", server.getLog)
	router.GET("/logs", server.listLogs)
	router.DELETE("/logs/:id", server.deleteLog)
	router.PUT("/logs/:id", server.updateLog)

	// Share routes
	router.POST("/shares", server.createShare)
	router.GET("/shares/:id", server.getShare)
	router.GET("/shares", server.listShares)
	router.DELETE("/shares/:id", server.deleteShare)
	router.PUT("/shares/:id", server.updateShare)

	// ShoppingListItem routes
	router.POST("/shopping-list-items", server.createShoppingListItem)
	router.GET("/shopping-list-items/:id", server.getShoppingListItem)
	router.GET("/shopping-list-items", server.listShoppingListItems)
	router.DELETE("/shopping-list-items/:id", server.deleteShoppingListItem)
	router.PUT("/shopping-list-items/:id", server.updateShoppingListItem)

	// Storage routes
	router.POST("/storages", server.createStorage)
	router.GET("/storages/:id", server.getStorage)
	router.GET("/storages", server.listStorages)
	router.DELETE("/storages/:id", server.deleteStorage)
	router.PUT("/storages/:id", server.updateStorage)

	// User routes
	router.POST("/users", server.createUser)
	router.GET("/users/:id", server.getUser)
	router.GET("/users", server.listUsers)
	router.DELETE("/users/:id", server.deleteUser)
	router.PUT("/users/:id", server.updateUser)

	// UserStatus routes
	router.POST("/user-statuses", server.createUserStatus)
	router.GET("/user-statuses/:id", server.getUserStatus)
	router.GET("/user-statuses", server.listUserStatuses)
	router.DELETE("/user-statuses/:id", server.deleteUserStatus)
	router.PUT("/user-statuses/:id", server.updateUserStatus)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
