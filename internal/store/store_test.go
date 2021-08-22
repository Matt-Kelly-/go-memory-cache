package store_test

import (
	"github.com/Matt-Kelly-/go-memory-cache/internal/store"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"sync"
	"testing"
)

func createDefaultStore() store.Store {
	return store.NewStore()
}

func createStoreWithMutexDecorator() store.Store {
	return store.WithMutex(store.NewStore())
}

func createStoreWithRWMutexDecorator() store.Store {
	return store.WithRWMutex(store.NewStore())
}

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
		createStore: createStoreWithMutexDecorator,
	})
}

func TestRWMutexDecoratorLocking(t *testing.T) {
	suite.Run(t, &storeLockingTestSuite{
		createStore: createStoreWithRWMutexDecorator,
	})
}

func benchmarkHas(b *testing.B, createStore func() store.Store) {
	b.Run("serial miss", func(b *testing.B) {
		testStore := createStore()
		testKey := "test key"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			testStore.Has(testKey)
		}
	})

	b.Run("serial hit", func(b *testing.B) {
		testStore := createStore()
		testKey := "test key"
		testStore.Put(testKey, "test value")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			testStore.Has(testKey)
		}
	})

	b.Run("parallel miss", func(b *testing.B) {
		testStore := createStore()
		testKey := "test key"
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				testStore.Has(testKey)
			}
		})
	})

	b.Run("parallel hit", func(b *testing.B) {
		testStore := createStore()
		testKey := "test key"
		testStore.Put(testKey, "test value")
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				testStore.Has(testKey)
			}
		})
	})

	b.Run("parallel mixed", func(b *testing.B) {
		testStore := createStore()
		testKeys := []string{"test key 1", "test key 2"}
		testStore.Put(testKeys[0], "test value")
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			index := 0
			for pb.Next() {
				index = (index + 1) % 2
				testStore.Has(testKeys[index])
			}
		})
	})
}

func benchmarkGet(b *testing.B, createStore func() store.Store) {
	b.Run("serial miss", func(b *testing.B) {
		testStore := createStore()
		testKey := "test key"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			testStore.Get(testKey)
		}
	})

	b.Run("serial hit", func(b *testing.B) {
		testStore := createStore()
		testKey := "test key"
		testStore.Put(testKey, "test value")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			testStore.Get(testKey)
		}
	})

	b.Run("parallel miss", func(b *testing.B) {
		testStore := createStore()
		testKey := "test key"
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				testStore.Get(testKey)
			}
		})
	})

	b.Run("parallel hit", func(b *testing.B) {
		testStore := createStore()
		testKey := "test key"
		testStore.Put(testKey, "test value")
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				testStore.Get(testKey)
			}
		})
	})

	b.Run("parallel mixed", func(b *testing.B) {
		testStore := createStore()
		testKeys := []string{"test key 1", "test key 2"}
		testStore.Put(testKeys[0], "test value")
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			index := 0
			for pb.Next() {
				index = (index + 1) % 2
				testStore.Get(testKeys[index])
			}
		})
	})
}

func benchmarkPut(b *testing.B, createStore func() store.Store, parallel bool) {
	b.Run("serial miss", func(b *testing.B) {
		testStore := createStore()
		testKey := "test key"
		testValue := "test value"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			testStore.Put(testKey, testValue)
		}
	})

	b.Run("serial hit", func(b *testing.B) {
		testStore := createStore()
		testKey := "test key"
		testValue := "test value"
		testStore.Put(testKey, testValue)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			testStore.Put(testKey, testValue)
		}
	})

	if !parallel {
		return
	}

	b.Run("parallel miss", func(b *testing.B) {
		testStore := createStore()
		testKey := "test key"
		testValue := "test value"
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				testStore.Put(testKey, testValue)
			}
		})
	})

	b.Run("parallel hit", func(b *testing.B) {
		testStore := createStore()
		testKey := "test key"
		testValue := "test value"
		testStore.Put(testKey, testValue)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				testStore.Put(testKey, testValue)
			}
		})
	})

	b.Run("parallel mixed", func(b *testing.B) {
		testStore := createStore()
		testKeys := []string{"test key 1", "test key 2"}
		testValues := []string{"test value 1", "test value 2"}
		testStore.Put(testKeys[0], testValues[0])
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			index := 0
			for pb.Next() {
				index = (index + 1) % 2
				testStore.Put(testKeys[index], testValues[index])
			}
		})
	})
}

func benchmarkDelete(b *testing.B, createStore func() store.Store, parallel bool) {
	b.Run("serial miss", func(b *testing.B) {
		testStore := createStore()
		testKey := "test key"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			testStore.Delete(testKey)
		}
	})

	if !parallel {
		return
	}

	b.Run("parallel miss", func(b *testing.B) {
		testStore := createStore()
		testKey := "test key"
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				testStore.Delete(testKey)
			}
		})
	})
}

func BenchmarkDefaultStoreHas(b *testing.B) {
	benchmarkHas(b, createDefaultStore)
}

func BenchmarkDefaultStoreGet(b *testing.B) {
	benchmarkGet(b, createDefaultStore)
}

func BenchmarkDefaultStorePut(b *testing.B) {
	benchmarkPut(b, createDefaultStore, false)
}

func BenchmarkDefaultStoreDelete(b *testing.B) {
	benchmarkDelete(b, createDefaultStore, false)
}

func BenchmarkMutexDecoratorHas(b *testing.B) {
	benchmarkHas(b, createStoreWithMutexDecorator)
}

func BenchmarkMutexDecoratorGet(b *testing.B) {
	benchmarkGet(b, createStoreWithMutexDecorator)
}

func BenchmarkMutexDecoratorPut(b *testing.B) {
	benchmarkPut(b, createStoreWithMutexDecorator, true)
}

func BenchmarkMutexDecoratorDelete(b *testing.B) {
	benchmarkDelete(b, createStoreWithMutexDecorator, true)
}

func BenchmarkRWMutexDecoratorHas(b *testing.B) {
	benchmarkHas(b, createStoreWithRWMutexDecorator)
}

func BenchmarkRWMutexDecoratorGet(b *testing.B) {
	benchmarkGet(b, createStoreWithRWMutexDecorator)
}

func BenchmarkRWMutexDecoratorPut(b *testing.B) {
	benchmarkPut(b, createStoreWithRWMutexDecorator, true)
}

func BenchmarkRWMutexDecoratorDelete(b *testing.B) {
	benchmarkDelete(b, createStoreWithRWMutexDecorator, true)
}
