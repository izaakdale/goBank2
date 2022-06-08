package main

import (
	"database/sql"
	"log"

	"github.com/izaakdale/goBank2/api"
	db "github.com/izaakdale/goBank2/db/sqlc"
	_ "github.com/lib/pq"
)

var driver = "postgres"
var source = "postgresql://root:secret@localhost:5432/goBank?sslmode=disable"
var serverAddress = "localhost:8080"

func main() {
	dbConn, err := sql.Open(driver, source)
	if err != nil {
		log.Fatal("Failed to connect to db")
	}

	store := db.NewStore(dbConn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("Unable to start the server: ", err.Error())
	}
}
