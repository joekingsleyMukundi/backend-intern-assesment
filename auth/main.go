package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/joekingsleyMukundi/backend-intern-assesment/auth/gapi"
	"github.com/joekingsleyMukundi/backend-intern-assesment/auth/pb"
	db "github.com/joekingsleyMukundi/backend-intern-assesment/common/db/sqlc"
	"github.com/joekingsleyMukundi/backend-intern-assesment/common/util"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := util.LoadConfig("../common/")
	if err != nil {
		log.Fatal("connot load config: ", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("connot connect to db: ", err)
	}
	store := db.NewStore(conn)
	runGrpcServer(config, store)
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot load server: ", err)
	}
	grpcserver := grpc.NewServer()
	pb.RegisterAuthServer(grpcserver, server)
	reflection.Register(grpcserver)
	listener, err := net.Listen("tcp", config.AuthSeviceGrpcServerAddress)
	if err != nil {
		log.Fatal("cannot create listener due  to: ", err)
	}
	log.Printf("Start grpc server at %s", listener.Addr().String())
	err = grpcserver.Serve(listener)
	if err != nil {
		log.Fatal("cannot start start server: ", err)
	}
}
