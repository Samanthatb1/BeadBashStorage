package main

import (
	"database/sql"
	"log"

	"github.com/samanthatb1/beadBashStorage/api"
	db "github.com/samanthatb1/beadBashStorage/db/sqlc"
	"github.com/samanthatb1/beadBashStorage/util"
)

// Starts server and DB connection
func main() {
	// Load variables from env file
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load configurations (file / env): " , err)
	}

	// Connect to postgres DB
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil { log.Fatal("Cannot connect to db: ", err) }

	store := db.NewStore(conn) // Create a store instance with original and additional DB operations
	server := api.NewServer(store) // Create server instance based on the store

	// Start server by passing in an address to run on
	err = server.Start(config.ServerAddress)
	if err != nil { log.Fatal("Cannot start server: ", err) }
}