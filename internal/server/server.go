package server

import (
	"context"
	"github.com/Matt-Kelly-/go-memory-cache/api"
	"github.com/Matt-Kelly-/go-memory-cache/internal/store"
)

type defaultServer struct {
	api.UnimplementedCacheServer

	store store.Store
}

func NewServer(store store.Store) api.CacheServer {
	return defaultServer{
		store: store,
	}
}

func (s defaultServer) Has(ctx context.Context, request *api.HasRequest) (*api.HasResponse, error) {
	result := s.store.Has(request.Key)
	return &api.HasResponse{
		Exists: result,
	}, nil
}

func (s defaultServer) Get(ctx context.Context, request *api.GetRequest) (*api.GetResponse, error) {
	value, exists := s.store.Get(request.Key)
	return &api.GetResponse{
		Exists: exists,
		Value:  value,
	}, nil
}

func (s defaultServer) Put(ctx context.Context, request *api.PutRequest) (*api.PutResponse, error) {
	s.store.Put(request.Key, request.Value)
	return &api.PutResponse{}, nil
}

func (s defaultServer) Delete(ctx context.Context, request *api.DeleteRequest) (*api.DeleteResponse, error) {
	s.store.Delete(request.Key)
	return &api.DeleteResponse{}, nil
}
