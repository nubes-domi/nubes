package rpc

import (
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewServer() (*grpc.Server, *grpcweb.WrappedGrpcServer) {
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(ServerInterceptor))
	RegisterSessionsServer(grpcServer, &SessionsServerImpl{})
	RegisterUsersServer(grpcServer, &UsersServerImpl{})
	reflection.Register(grpcServer)

	web := grpcweb.WrapServer(grpcServer)

	return grpcServer, web
}
