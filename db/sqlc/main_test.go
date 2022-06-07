package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDb *sql.DB
var driver = "postgres"
var source = "postgresql://root:secret@localhost:5432/goBank?sslmode=disable"

func TestMain(m *testing.M) {

	var err error
	testDb, err = sql.Open(driver, source)
	if err != nil {
		log.Fatal("Failed to connect to db")
	}

	testQueries = New(testDb)
	os.Exit(m.Run())
}
