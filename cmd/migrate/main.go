package main

import (
	"errors"
	"fmt"
	"os"

	"project_expiry_backend/utils"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if err := run(); err != nil {
		_, err := fmt.Fprintln(os.Stderr, err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) != 2 || os.Args[1] != "up" {
		return errors.New("usage: go run ./cmd/migrate up")
	}

	config, err := utils.LoadConfig(".")
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if config.DBSource == "" {
		return errors.New("DB_SOURCE is not configured")
	}

	m, err := migrate.New("file://db/migration", config.DBSource)
	if err != nil {
		return errors.New("failed to initialize migrator")
	}
	defer func() {
		_, _ = m.Close()
	}()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return errors.New("database migration failed")
	}

	return nil
}
