package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDBConn *sql.DB

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:P@ssw0rd@localhost:5432/simple_bank?sslmode=disable"
)

func TestMain(m *testing.M) {
	var err error
	testDBConn, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to database", err)
	}

	testQueries = New(testDBConn)

	os.Exit(m.Run())
}
