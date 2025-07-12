package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"project_expiry_backend/api"
	sqlc "project_expiry_backend/db/sqlc"
	"project_expiry_backend/utils"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}

	store := sqlc.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatalf("cannot create server: %v", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatalf("cannot start server: %v", err)
	}
	log.Print("server started on ", config.ServerAddress)
}
