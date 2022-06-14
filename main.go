package main

import (
	"database/sql"
	"log"

	"github.com/izaakdale/goBank2/api"
	db "github.com/izaakdale/goBank2/db/sqlc"
	"github.com/izaakdale/goBank2/util"
	_ "github.com/lib/pq"
)

// var driver = "postgres"
// var source = "postgresql://root:secret@localhost:5432/goBank?sslmode=disable"
// var serverAddress = "localhost:8080"

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Unable to load config file")
	}

	dbConn, err := sql.Open(config.DbDriver, config.DbSoruce)
	if err != nil {
		log.Fatal("Failed to connect to db")
	}

	store := db.NewStore(dbConn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Failed to create server: " + err.Error())
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Unable to start the server: ", err.Error())
	}
}
