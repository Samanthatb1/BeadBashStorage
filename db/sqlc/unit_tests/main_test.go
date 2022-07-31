package tests

import (
	"database/sql" // provides a generic interface for SQL
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq" // provides the DB driver

	sqlc "github.com/samanthatb1/beadBashStorage/db/sqlc"
	"github.com/samanthatb1/beadBashStorage/util"
)

var testQueries *sqlc.Queries
var testDB *sql.DB

func TestMain(m *testing.M){
	// Load ENV variables
	config, err := util.LoadConfig("../../../")
	if err != nil {
		log.Fatal("Cannot load configurations (file / env): " , err)
	}

	// Connect to postgres DB
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil{ log.Fatal("Cannot connect to db: " , err) }
	testQueries = sqlc.New(testDB)

	os.Exit(m.Run());
}