package test

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
	db "project_expiry_backend/db/sqlc"
	"project_expiry_backend/utils"
	"testing"
)

var testQueries *db.Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := utils.LoadConfig("../../..")
	if err != nil {
		log.Fatalf("cannot load config file: %v", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}
	testQueries = db.New(testDB)

	os.Exit(m.Run())
}
