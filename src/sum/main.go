package main

import (
	"log"
	"net"
	"nubes/sum/db"
	"nubes/sum/router"
	"nubes/sum/rpc"
	"nubes/sum/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	utils.PrepareKeys()
	db.InitDatabase()

	lis, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(rpc.ServerInterceptor),
	)
	rpc.RegisterSessionsServer(grpcServer, &rpc.SessionsServerImpl{})
	reflection.Register(grpcServer)
	go grpcServer.Serve(lis)

	router := router.New()
	router.Run("localhost:8080")
}
