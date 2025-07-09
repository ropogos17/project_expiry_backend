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

	router.POST("/roles", server.createRole)
	router.GET("/roles/:id", server.getRole)
	router.GET("/roles", server.listRoles)
	router.DELETE("/roles", server.deleteRole)
	router.PUT("/roles/:id", server.updateRole)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
