package test

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
	db "project_expiry_backend/db/sqlc"
	"testing"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://ropogos:S3cr3t@localhost:5432/project_expiry_db?sslmode=disable"
)

var testQueries *db.Queries

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}
	testQueries = db.New(conn)

	os.Exit(m.Run())
}
