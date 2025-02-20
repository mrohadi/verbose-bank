package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/mrohadi/simplebank/cmd/api"
	db "github.com/mrohadi/simplebank/db/sqlc"
	"github.com/mrohadi/simplebank/utils"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot read config")
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to database", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Cannot started the server")
	}
}
