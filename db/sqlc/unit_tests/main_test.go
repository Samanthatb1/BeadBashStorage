package tests

import (
	"database/sql" // provides a generic interface for SQL
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq" // provides the DB driver

	sqlc "github.com/samanthatb1/beadBashStorage/db/sqlc"
)

var testQueries *sqlc.Queries
var testDB *sql.DB

const (
	dbDriver = "postgres"
	dbSource = "postgres://root:secret@localhost:5432/BB-DB?sslmode=disable"
)

func TestMain(m *testing.M){
	var err error
	// Connect to postgres DB
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil{ log.Fatal("Cannot connect to db: " , err) }
	testQueries = sqlc.New(testDB)

	os.Exit(m.Run());
}