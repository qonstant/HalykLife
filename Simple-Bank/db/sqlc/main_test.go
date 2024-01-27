package db

import (
	"database/sql"
	"log"
	"os"
	"testing"
	_ "github.com/lib/pq"
	"Simple-Bank/util"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	// conn, err := sql.Open(dbDriver, dbSource)
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Failed to open config", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Failed to open database", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
