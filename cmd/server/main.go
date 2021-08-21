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
	args := os.Args[1:]
	storeType := ""
	if len(args) > 0 {
		storeType = args[0]
	}

	cacheStore := store.NewStore()

	switch storeType {
	case "mutex":
		log.Print("Protecting store with mutex")
		cacheStore = store.WithMutex(cacheStore)
	case "rwmutex":
		log.Print("Protecting store with rw mutex")
		cacheStore = store.WithRWMutex(cacheStore)
	default:
		log.Print("Leaving store unprotected")
	}

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
