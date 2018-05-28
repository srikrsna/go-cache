package cache // import "github.com/srikrsna/go-cache"

import (
	"errors"
	"io"
	"time"
)

var (
	// ErrCacheMiss is the error returned by the Cache interface upon a cache miss
	ErrCacheMiss = errors.New("cache miss")
)

//go:generate mockgen -destination=mocks/cache_mock.go -package=cache_mock github.com/srikrsna/go-cache Backend

// Backend is the low level cache backend typically, redis, and memcached
type Backend interface {
	Get(key string, w io.Writer, d time.Duration) error
	Set(key string, r []byte, d time.Duration) error
	Renew(key string, d time.Duration) error
	Delete(key string) error
}

//go:generate mockgen -destination=mocks/backend_mock.go -package=cache_mock github.com/srikrsna/go-cache Cache

// Cache is the high level interface that caches go's objects using a defined serialization format
// Examples: JSON, Protocol Buffers, Message Pack, Bson, Gob
type Cache interface {
	Get(key string, v interface{}, d time.Duration) error
	Set(key string, v interface{}, d time.Duration) error
	Renew(key string, d time.Duration) error
	Delete(key string) error
}
