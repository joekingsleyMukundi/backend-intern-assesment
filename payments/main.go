package main

import (
	"database/sql"
	"log"
	"net"

	db "github.com/joekingsleyMukundi/backend-intern-assesment/common/db/sqlc"
	"github.com/joekingsleyMukundi/backend-intern-assesment/common/util"
	"github.com/joekingsleyMukundi/backend-intern-assesment/payments/gapi"
	"github.com/joekingsleyMukundi/backend-intern-assesment/payments/pb"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := util.LoadConfig("../common/")
	if err != nil {
		log.Fatal("cannot load config file: ", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot conect to db: ", err)
	}
	store := db.NewStore(conn)
	runGrpcServer(config, store)
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot start gpai server: ", err)
	}
	grpcSever := grpc.NewServer()
	pb.RegisterPaymentsServer(grpcSever, server)
	reflection.Register(grpcSever)
	listener, err := net.Listen("tcp", config.PaymentSeviceGrpcServerAddress)
	if err != nil {
		log.Fatal("cannot create listener due  to: ", err)
	}
	log.Printf("Start grpc server at %s", listener.Addr().String())
	err = grpcSever.Serve(listener)
	if err != nil {
		log.Fatal("cannot start start server: ", err)
	}
}
