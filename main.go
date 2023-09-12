package main

import (
	"database/sql"
	"log"

	"github.com/tonisco/simple-bank-go/api"
	db "github.com/tonisco/simple-bank-go/db/sqlc"

	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://postgres:user@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("Failed to connect to db:", err)
	}

	store := db.NewStore(conn)

	server := api.NewServer(&store)

	err = server.Start(serverAddress)

	if err != nil {
		log.Fatal("Cannot start server", err)
	}
}
