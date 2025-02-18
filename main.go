package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/mrohadi/simplebank/cmd/api"
	db "github.com/mrohadi/simplebank/db/sqlc"
)

const (
	dbDriver   = "postgres"
	dbSource   = "postgresql://root:P@ssw0rd@localhost:5432/simple_bank?sslmode=disable"
	serverAddr = "0.0.0.0:8000"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to database", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddr)
	if err != nil {
		log.Fatal("Cannot started the server")
	}
}
