package server

import (
	"context"
	"github.com/Matt-Kelly-/go-memory-cache/api"
	"github.com/Matt-Kelly-/go-memory-cache/internal/store"
	"log"
)

type defaultServer struct {
	api.UnimplementedCacheServer

	store  store.Store
	logger *log.Logger
}

func NewServer(store store.Store, logger *log.Logger) api.CacheServer {
	return defaultServer{
		store:  store,
		logger: logger,
	}
}

func (s defaultServer) Has(ctx context.Context, request *api.HasRequest) (*api.HasResponse, error) {
	s.logger.Printf("Request: Has %v", request)
	result := s.store.Has(request.Key)
	return &api.HasResponse{
		Exists: result,
	}, nil
}

func (s defaultServer) Get(ctx context.Context, request *api.GetRequest) (*api.GetResponse, error) {
	s.logger.Printf("Request: Get %v", request)
	value, exists := s.store.Get(request.Key)
	return &api.GetResponse{
		Exists: exists,
		Value:  value,
	}, nil
}

func (s defaultServer) Put(ctx context.Context, request *api.PutRequest) (*api.PutResponse, error) {
	s.logger.Printf("Request: Put %v", request)
	s.store.Put(request.Key, request.Value)
	return &api.PutResponse{}, nil
}

func (s defaultServer) Delete(ctx context.Context, request *api.DeleteRequest) (*api.DeleteResponse, error) {
	s.logger.Printf("Request: Delete %v", request)
	s.store.Delete(request.Key)
	return &api.DeleteResponse{}, nil
}
