package main

import (
	"net/http"
	"nubes/sum/db"
	"nubes/sum/router"
	"nubes/sum/rpc"
	"nubes/sum/utils"

	"github.com/gin-gonic/gin"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
)

type Handler struct {
	http *gin.Engine
	grpc *grpcweb.WrappedGrpcServer
}

func NewHandler() *Handler {
	router := router.New()
	_, server := rpc.NewServer()

	return &Handler{
		http: router,
		grpc: server,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if h.grpc.IsGrpcWebRequest(req) {
		h.grpc.ServeHTTP(w, req)
		return
	}
	h.http.ServeHTTP(w, req)
}

func main() {
	utils.PrepareKeys()
	db.InitDatabase()

	http.ListenAndServe(":8080", NewHandler())
}
