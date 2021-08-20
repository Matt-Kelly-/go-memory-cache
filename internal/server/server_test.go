package server_test

import (
	"context"
	"github.com/Matt-Kelly-/go-memory-cache/api"
	"github.com/Matt-Kelly-/go-memory-cache/internal/server"
	"github.com/Matt-Kelly-/go-memory-cache/internal/store"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHas(t *testing.T) {
	t.Run("exists", func(t *testing.T) {
		mockStore := new(store.MockStore)
		mockStore.On("Has", "test key").Return(true)

		testServer := server.NewServer(mockStore)
		response, err := testServer.Has(context.Background(), &api.HasRequest{
			Key: "test key",
		})

		require.NotNil(t, response)
		require.True(t, response.Exists)
		require.Nil(t, err)
	})

	t.Run("does not exist", func(t *testing.T) {
		mockStore := new(store.MockStore)
		mockStore.On("Has", "test key").Return(false)

		testServer := server.NewServer(mockStore)
		response, err := testServer.Has(context.Background(), &api.HasRequest{
			Key: "test key",
		})

		require.NotNil(t, response)
		require.False(t, response.Exists)
		require.Nil(t, err)
	})
}
