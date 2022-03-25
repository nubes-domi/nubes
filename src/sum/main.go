package main

import (
	"log"
	"net"
	"net/http"
	"nubes/sum/db"
	"nubes/sum/router"
	"nubes/sum/rpc"
	"nubes/sum/utils"

	"github.com/gin-gonic/gin"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
)

type Handler struct {
	http    *gin.Engine
	grpc    *grpc.Server
	grpcWeb *grpcweb.WrappedGrpcServer
}

func NewHandler() *Handler {
	router := router.New()
	grpc, grpcWeb := rpc.NewServer()

	return &Handler{
		http:    router,
		grpc:    grpc,
		grpcWeb: grpcWeb,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if h.grpcWeb.IsGrpcWebRequest(req) {
		h.grpcWeb.ServeHTTP(w, req)
		return
	}
	h.http.ServeHTTP(w, req)
}

func main() {
	utils.PrepareKeys()
	db.InitDatabase()

	handler := NewHandler()

	lis, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	go handler.grpc.Serve(lis)

	http.ListenAndServe(":8080", handler)
}
