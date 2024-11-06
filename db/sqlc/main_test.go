package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/vansh123456/simplebank/util"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("couldnt load config of main_test", err)
	}
	testDB, err = sql.Open(config.DBdriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db", err)
	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}
