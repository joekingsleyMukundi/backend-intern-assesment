package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbUri    = "postgresql://root:mypassword@localhost:5432/payd?sslmode=disable"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	testDb, err := sql.Open(dbDriver, dbUri)
	if err != nil {
		log.Fatal("ConnectioN error: ", err)
	}
	testQueries = New(testDb)
	os.Exit(m.Run())
}
