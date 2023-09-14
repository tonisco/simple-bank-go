package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/tonisco/simple-bank-go/util"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	var err error

	config, err := util.LoadConfig("../..")

	if err != nil {
		log.Fatal("cannot load environment variables")
	}

	testDb, err = sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("Failed to connect to db:", err)
	}

	testQueries = New(testDb)

	os.Exit(m.Run())
}
