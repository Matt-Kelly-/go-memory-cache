package store

//go:generate mockery --name=Store --inpackage --case underscore

type Store interface {
	Has(key string) bool
	Get(key string) (string, bool)
	Put(key, value string)
	Delete(key string)
}

type defaultStore struct {
	contents map[string]string
}

func newDefaultStore() *defaultStore {
	return &defaultStore{
		contents: make(map[string]string),
	}
}

func NewStore() Store {
	return newDefaultStore()
}

func NewStoreWithContents(contents map[string]string) Store {
	s := newDefaultStore()
	for k, v := range contents {
		s.contents[k] = v
	}
	return s
}

func (s *defaultStore) Has(key string) bool {
	_, ok := s.contents[key]
	return ok
}

func (s *defaultStore) Get(key string) (string, bool) {
	value, ok := s.contents[key]
	return value, ok
}

func (s *defaultStore) Put(key, value string) {
	s.contents[key] = value
}

func (s *defaultStore) Delete(key string) {
	delete(s.contents, key)
}
