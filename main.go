package main

import (
	"database/sql"
	"log"

	"github.com/tonisco/simple-bank-go/api"
	db "github.com/tonisco/simple-bank-go/db/sqlc"
	"github.com/tonisco/simple-bank-go/util"

	_ "github.com/lib/pq"
)

func main() {
	config,err := util.LoadConfig(".")

	if err != nil{
		log.Fatal("Cannot not load environment variables", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("Failed to connect to db:", err)
	}

	store := db.NewStore(conn)

	server := api.NewServer(&store)

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("Cannot start server", err)
	}
}
