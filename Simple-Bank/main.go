package main

import (
	"Simple-Bank/api"
	db "Simple-Bank/db/sqlc"
	"Simple-Bank/util"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Failed to open database", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Failed to create server", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Failed to start server", err)
	}

}
