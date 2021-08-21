package main

import (
	"github.com/Matt-Kelly-/go-memory-cache/api"
	"github.com/Matt-Kelly-/go-memory-cache/internal/server"
	"github.com/Matt-Kelly-/go-memory-cache/internal/store"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

const (
	port = ":50051"
)

func main() {
	cacheStore := store.NewStore()

	logger := log.New(os.Stdout, "Server ", log.LstdFlags)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	api.RegisterCacheServer(grpcServer, server.NewServer(cacheStore, logger))

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
