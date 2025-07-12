package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	db "project_expiry_backend/db/sqlc"
	"project_expiry_backend/token"
)

type Server struct {
	store      *db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(store *db.Store) (*Server, error) {
	// TODO: replace the symmetricKey with a real one
	tokenMaker, err := token.NewPasetoMaker("")
	if err != nil {
		return nil, fmt.Errorf("Cannot create token maker: %w", err)
	}
	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
	}
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
