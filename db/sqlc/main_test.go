package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/izaakdale/goBank2/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDb *sql.DB
var driver = "postgres"
var source = "postgresql://root:secret@localhost:5432/goBank?sslmode=disable"

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Unable to load config file during testing")
	}

	testDb, err = sql.Open(config.DbDriver, config.DbSoruce)
	if err != nil {
		log.Fatal("Failed to connect to db")
	}

	testQueries = New(testDb)
	os.Exit(m.Run())
}
