package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"project_expiry_backend/config"
	"project_expiry_backend/db"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config.LoadConfig()

	postgresDB, err := db.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to postgres database: %v", err)
	}
	defer postgresDB.Close()

	//userHandler := handlers.NewUserHandler(postgresDB)

	router := gin.Default()
	//routes.SetupRoutes(router, userHandler)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	fmt.Printf("Server started on port 8080")
}
