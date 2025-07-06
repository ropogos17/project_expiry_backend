package routes

import (
	"github.com/gin-gonic/gin"
	"project_expiry_backend/handlers"
)

func SetupRoutes(r *gin.Engine, handler *handlers.UserHandler) {
	r.GET("/users", handler.GetUsers)
	r.GET("/user/:id", handler.GetUser)
	r.POST("/user", handler.CreateUser)
}
