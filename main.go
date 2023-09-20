package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/tonisco/simple-bank-go/api"
	db "github.com/tonisco/simple-bank-go/db/sqlc"
	"github.com/tonisco/simple-bank-go/gapi"
	"github.com/tonisco/simple-bank-go/pb"
	"github.com/tonisco/simple-bank-go/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("Cannot not load environment variables", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("Failed to connect to db:", err)
	}

	store := db.NewStore(conn)
	runGRPCServer(config, store)
}

func runGRPCServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("Cannot start server", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot create listener")
	}
	log.Printf("started grpc server on %s", listener.Addr().String())

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start grpc server")
	}
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Cannot start server", err)
	}

	err = server.Start(config.HTTPServerAddress)

	if err != nil {
		log.Fatal("Cannot start server", err)
	}
}
