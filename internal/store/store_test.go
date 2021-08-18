package store_test

import (
	"github.com/Matt-Kelly-/go-memory-cache/internal/store"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHas(t *testing.T) {

	t.Run("empty store", func(t *testing.T) {
		s := store.NewStore()
		result := s.Has("test key")
		require.False(t, result)
	})

	t.Run("wrong key", func(t *testing.T) {
		s := store.NewStoreWithContents(map[string]string{
			"test key": "test value",
		})
		result := s.Has("other key")
		require.False(t, result)
	})

	t.Run("right key", func(t *testing.T) {
		s := store.NewStoreWithContents(map[string]string{
			"test key": "test value",
		})
		result := s.Has("test key")
		require.True(t, result)
	})

}

func TestGet(t *testing.T) {

	t.Run("empty store", func(t *testing.T) {
		s := store.NewStore()
		value, ok := s.Get("test key")
		require.Empty(t, value)
		require.False(t, ok)
	})

	t.Run("wrong key", func(t *testing.T) {
		s := store.NewStoreWithContents(map[string]string{
			"test key": "test value",
		})
		value, ok := s.Get("other key")
		require.Empty(t, value)
		require.False(t, ok)
	})

	t.Run("right key", func(t *testing.T) {
		s := store.NewStoreWithContents(map[string]string{
			"test key": "test value",
		})
		value, ok := s.Get("test key")
		require.Equal(t, "test value", value)
		require.True(t, ok)
	})

}

func TestPut(t *testing.T) {

	t.Run("empty store", func(t *testing.T) {
		s := store.NewStore()
		s.Put("test key", "test value")
		value, ok := s.Get("test key")
		require.Equal(t, "test value", value)
		require.True(t, ok)
	})

	t.Run("replace", func(t *testing.T) {
		s := store.NewStoreWithContents(map[string]string{
			"test key": "test value",
		})
		s.Put("test key", "new test value")
		value, ok := s.Get("test key")
		require.Equal(t, "new test value", value)
		require.True(t, ok)
	})

}

func TestDelete(t *testing.T) {

	t.Run("empty store", func(t *testing.T) {
		s := store.NewStore()
		s.Delete("test key")
		require.False(t, s.Has("test key"))
	})

	t.Run("wrong key", func(t *testing.T) {
		s := store.NewStoreWithContents(map[string]string{
			"test key": "test value",
		})
		s.Delete("other key")
		require.True(t, s.Has("test key"))
		require.False(t, s.Has("other key"))
	})

	t.Run("right key", func(t *testing.T) {
		s := store.NewStoreWithContents(map[string]string{
			"test key": "test value",
		})
		s.Delete("test key")
		require.False(t, s.Has("test key"))
	})

}
