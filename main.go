package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/vansh123456/simplebank/api"
	db "github.com/vansh123456/simplebank/db/sqlc"
	"github.com/vansh123456/simplebank/util"
)

func main() {

	config, err := util.LoadConfig(".") // dot defines load the configs here
	if err != nil {
		log.Fatal("err loading configs:", err)
	}
	conn, err := sql.Open(config.DBdriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
