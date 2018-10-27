package memory

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"time"

	"github.com/srikrsna/go-cache"
)

// ErrNotPtr is returned by Get
var ErrNotPtr = errors.New("can only use pointers")

// Cache implements an in memory object cache
type Cache struct {
	sync.Mutex
	data map[string]interface{}

	timers map[string]*time.Timer
}

var _ cache.Cache = (*Cache)(nil)

// NewCache returns an in memory cache backend.
func NewCache() *Cache {
	return &Cache{
		data:   make(map[string]interface{}),
		timers: make(map[string]*time.Timer),
	}
}

// Get retrives an data structure from cache.
func (c *Cache) Get(ctx context.Context, key string, v interface{}, d time.Duration) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		return ErrNotPtr
	}

	c.Lock()
	defer c.Unlock()
	stored, ok := c.data[key]
	if !ok {
		return cache.ErrCacheMiss
	}

	val.Elem().Set(reflect.ValueOf(stored))

	t := c.timers[key]

	t.Stop()
	t.Reset(d)

	return nil
}

// Set marshals and saves an arbitary go data structure.
func (c *Cache) Set(ctx context.Context, key string, v interface{}, d time.Duration) error {
	c.Lock()
	defer c.Unlock()

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	c.data[key] = val.Interface()
	t := c.timers[key]
	if t != nil {
		t.Stop()
		t.Reset(d)
	} else {
		c.timers[key] = time.AfterFunc(d, func() {
			c.Lock()
			defer c.Unlock()

			delete(c.data, key)
			delete(c.timers, key)
		})
	}

	return nil
}

// Renew is used to renew a cache
func (c *Cache) Renew(ctx context.Context, key string, d time.Duration) error {
	c.Lock()
	defer c.Unlock()

	t := c.timers[key]
	if t == nil {
		return cache.ErrCacheMiss
	}

	t.Stop()
	t.Reset(d)

	return nil
}

// Evict is used to evict a key from cache.
func (c *Cache) Evict(ctx context.Context, key string) error {
	c.Lock()
	defer c.Unlock()

	delete(c.data, key)
	t := c.timers[key]
	if t == nil {
		return cache.ErrCacheMiss
	}

	t.Stop()
	delete(c.timers, key)

	return nil
}
