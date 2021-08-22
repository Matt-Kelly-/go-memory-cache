# go-memory-cache

## About The Project

A simple in-memory cache with GRPC API. This is a toy project to help me learn Go. The cache stores string values under string keys. It supports the following operations:
- `Has` Checks for the existence of a key. Returns a boolean indicating the existence.
- `Get` Reads the value for a key. Returns a boolean indicating the existence, and the value (or empty string if it doesn't exist).
- `Put` Sets the value for a key.
- `Delete` Deletes the value for a key.

I wanted to try different approaches to synchronisation, so the default cache store has no protection. I used the decorator pattern to create 2 wrappers to protect the cache with `sync.Mutex` and `sync.RWMutex` respectively.

### Built With
- [GRPC](https://grpc.io/)
- [Testify](https://github.com/stretchr/testify)
