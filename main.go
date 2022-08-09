package main

import (
	"database/sql"
	"log"

	"github.com/samanthatb1/beadBashStorage/api"
	db "github.com/samanthatb1/beadBashStorage/db/sqlc"
	"github.com/samanthatb1/beadBashStorage/util"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Starts server and DB connection
func main() {
	// Load variables from env file
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load configurations (file / env): " , err)
	}

	// Run Migrations
	runDBMigration(config.MigrationURL, config.DBSource)

	// Connect to postgres DB
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil { log.Fatal("Cannot connect to db: ", err) }

	store := db.NewStore(conn) // Create a store instance with original and additional DB operations
	server := api.NewServer(store) // Create server instance based on the store

	// Start server by passing in an address to run on
	err = server.Start(config.ServerAddress)
	if err != nil { log.Fatal("Cannot start server: ", err) }
	
	log.Println("server is up and running")
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal("cannot create new migrate instance:", err, "migrationURL: " , migrationURL, " db source: " , dbSource)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("failed to run migrate up:", err)
	}

	log.Println("db migrated successfully")
}