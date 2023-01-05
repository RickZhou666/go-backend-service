package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/RickZhou666/go-backend-service/util"
	_ "github.com/lib/pq"
)

// below are global variables
var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..") // go the parent of parent path
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// testDB, err := sql.Open(dbDriver, dbSource)
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
