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
