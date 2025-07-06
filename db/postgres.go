package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"project_expiry_backend/config"
)

func NewPostgresDB(cfg config.Config) (*sql.DB, error) {
	dsn := cfg.GetDSN()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return db, nil
}
