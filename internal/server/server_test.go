package server_test

import (
	"context"
	"github.com/Matt-Kelly-/go-memory-cache/api"
	"github.com/Matt-Kelly-/go-memory-cache/internal/server"
	"github.com/Matt-Kelly-/go-memory-cache/internal/store"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"testing"
)

func newLogger() *log.Logger {
	return log.New(ioutil.Discard, "", 0)
}

func TestHas(t *testing.T) {
	t.Run("exists", func(t *testing.T) {
		mockStore := new(store.MockStore)
		mockStore.On("Has", "test key").Return(true)

		testServer := server.NewServer(mockStore, newLogger())
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

		testServer := server.NewServer(mockStore, newLogger())
		response, err := testServer.Has(context.Background(), &api.HasRequest{
			Key: "test key",
		})

		require.NotNil(t, response)
		require.False(t, response.Exists)
		require.Nil(t, err)
	})
}

func TestGet(t *testing.T) {
	t.Run("exists", func(t *testing.T) {
		mockStore := new(store.MockStore)
		mockStore.On("Get", "test key").Return("test value", true)

		testServer := server.NewServer(mockStore, newLogger())
		response, err := testServer.Get(context.Background(), &api.GetRequest{
			Key: "test key",
		})

		require.NotNil(t, response)
		require.Equal(t, "test value", response.Value)
		require.True(t, response.Exists)
		require.Nil(t, err)
	})

	t.Run("does not exist", func(t *testing.T) {
		mockStore := new(store.MockStore)
		mockStore.On("Get", "test key").Return("", false)

		testServer := server.NewServer(mockStore, newLogger())
		response, err := testServer.Get(context.Background(), &api.GetRequest{
			Key: "test key",
		})

		require.NotNil(t, response)
		require.Equal(t, "", response.Value)
		require.False(t, response.Exists)
		require.Nil(t, err)
	})
}

func TestPut(t *testing.T) {
	mockStore := new(store.MockStore)
	mockStore.On("Put", "test key", "test value")

	testServer := server.NewServer(mockStore, newLogger())
	response, err := testServer.Put(context.Background(), &api.PutRequest{
		Key:   "test key",
		Value: "test value",
	})

	require.NotNil(t, response)
	require.Nil(t, err)

	mockStore.AssertCalled(t, "Put", "test key", "test value")
}

func TestDelete(t *testing.T) {
	mockStore := new(store.MockStore)
	mockStore.On("Delete", "test key")

	testServer := server.NewServer(mockStore, newLogger())
	response, err := testServer.Delete(context.Background(), &api.DeleteRequest{
		Key: "test key",
	})

	require.NotNil(t, response)
	require.Nil(t, err)

	mockStore.AssertCalled(t, "Delete", "test key")
}
