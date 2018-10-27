package cache

import (
	"context"
	"errors"
	"io"
	"time"
)

// ErrCacheMiss is the returned by the caching backend when a cache for the given key is not found.
// It should be returned by Get and Renew methods.
var ErrCacheMiss = errors.New("cache miss")

// Cache provides methods to set, get arbitary data structures.
// The exact serialization and backend depends on the implementation.
type Cache interface {
	// Get retrives an data structure from cache.
	Get(ctx context.Context, key string, v interface{}, d time.Duration) error

	// Set marshals and saves an arbitary go data structure.
	Set(ctx context.Context, key string, v interface{}, d time.Duration) error

	// Renew is used to renew a cache
	Renew(ctx context.Context, key string, d time.Duration) error

	// Evict is used to evict a key from cache.
	Evict(ctx context.Context, key string) error
}

// MarhsalUnmarshaler can be used to marhsal un marshal an object
type MarhsalUnmarshaler interface {
	// Marshal to marshals a data structure into the given io.Writer
	Marshal(w io.Writer, v interface{}) error

	// Unmarshal unmarshals the data structure present in r in its encoded form into v. v should be a pointer type.
	Unmarshal(r io.Reader, v interface{}) error
}
