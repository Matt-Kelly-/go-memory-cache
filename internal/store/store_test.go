package store_test

import (
	"github.com/Matt-Kelly-/go-memory-cache/internal/store"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"sync"
	"testing"
)

type storeTestSuite struct {
	suite.Suite

	createStore             func() store.Store
	createStoreWithContents func(map[string]string) store.Store
}

func (suite *storeTestSuite) TestHas() {

	suite.T().Run("empty store", func(t *testing.T) {
		s := suite.createStore()
		result := s.Has("test key")
		require.False(t, result)
	})

	suite.T().Run("wrong key", func(t *testing.T) {
		s := suite.createStoreWithContents(map[string]string{
			"test key": "test value",
		})
		result := s.Has("other key")
		require.False(t, result)
	})

	suite.T().Run("right key", func(t *testing.T) {
		s := suite.createStoreWithContents(map[string]string{
			"test key": "test value",
		})
		result := s.Has("test key")
		require.True(t, result)
	})

}

func (suite *storeTestSuite) TestGet() {

	suite.T().Run("empty store", func(t *testing.T) {
		s := suite.createStore()
		value, ok := s.Get("test key")
		require.Empty(t, value)
		require.False(t, ok)
	})

	suite.T().Run("wrong key", func(t *testing.T) {
		s := suite.createStoreWithContents(map[string]string{
			"test key": "test value",
		})
		value, ok := s.Get("other key")
		require.Empty(t, value)
		require.False(t, ok)
	})

	suite.T().Run("right key", func(t *testing.T) {
		s := suite.createStoreWithContents(map[string]string{
			"test key": "test value",
		})
		value, ok := s.Get("test key")
		require.Equal(t, "test value", value)
		require.True(t, ok)
	})

}

func (suite *storeTestSuite) TestPut() {

	suite.T().Run("empty store", func(t *testing.T) {
		s := suite.createStore()
		s.Put("test key", "test value")
		value, ok := s.Get("test key")
		require.Equal(t, "test value", value)
		require.True(t, ok)
	})

	suite.T().Run("replace", func(t *testing.T) {
		s := suite.createStoreWithContents(map[string]string{
			"test key": "test value",
		})
		s.Put("test key", "new test value")
		value, ok := s.Get("test key")
		require.Equal(t, "new test value", value)
		require.True(t, ok)
	})

}

func (suite *storeTestSuite) TestDelete() {

	suite.T().Run("empty store", func(t *testing.T) {
		s := suite.createStore()
		s.Delete("test key")
		require.False(t, s.Has("test key"))
	})

	suite.T().Run("wrong key", func(t *testing.T) {
		s := suite.createStoreWithContents(map[string]string{
			"test key": "test value",
		})
		s.Delete("other key")
		require.True(t, s.Has("test key"))
		require.False(t, s.Has("other key"))
	})

	suite.T().Run("right key", func(t *testing.T) {
		s := suite.createStoreWithContents(map[string]string{
			"test key": "test value",
		})
		s.Delete("test key")
		require.False(t, s.Has("test key"))
	})

}

func TestDefaultStore(t *testing.T) {
	suite.Run(t, &storeTestSuite{
		createStore: func() store.Store {
			return store.NewStore()
		},
		createStoreWithContents: func(contents map[string]string) store.Store {
			return store.NewStoreWithContents(contents)
		},
	})
}

func TestMutexDecorator(t *testing.T) {
	suite.Run(t, &storeTestSuite{
		createStore: func() store.Store {
			return store.WithMutex(store.NewStore())
		},
		createStoreWithContents: func(contents map[string]string) store.Store {
			return store.WithMutex(store.NewStoreWithContents(contents))
		},
	})
}

func TestRWMutexDecorator(t *testing.T) {
	suite.Run(t, &storeTestSuite{
		createStore: func() store.Store {
			return store.WithRWMutex(store.NewStore())
		},
		createStoreWithContents: func(contents map[string]string) store.Store {
			return store.WithRWMutex(store.NewStoreWithContents(contents))
		},
	})
}

type storeLockingTestSuite struct {
	suite.Suite

	createStore func() store.Store
}

func (suite *storeLockingTestSuite) TestLocking() {

	testStore := suite.createStore()

	// Run operations in parallel to allow for race detection
	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		testStore.Has("test key")
		wg.Done()
	}()
	go func() {
		testStore.Get("test key")
		wg.Done()
	}()
	go func() {
		testStore.Put("test key", "test value")
		wg.Done()
	}()
	go func() {
		testStore.Delete("test key")
		wg.Done()
	}()
	wg.Wait()
}

func TestMutexDecoratorLocking(t *testing.T) {
	suite.Run(t, &storeLockingTestSuite{
		createStore: func() store.Store {
			return store.WithMutex(store.NewStore())
		},
	})
}

func TestRWMutexDecoratorLocking(t *testing.T) {
	suite.Run(t, &storeLockingTestSuite{
		createStore: func() store.Store {
			return store.WithRWMutex(store.NewStore())
		},
	})
}
