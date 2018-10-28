package redis

import (
	"bytes"
	"context"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"

	"github.com/srikrsna/go-cache"
)

// Cache implements an in memory object cache
type Cache struct {
	Pool *redis.Pool

	MarshalUnmarshaler cache.MarhsalUnmarshaler
}

var _ cache.Cache = (*Cache)(nil)

// NewCache returns an in memory cache backend.
func NewCache(p *redis.Pool, m cache.MarhsalUnmarshaler) *Cache {
	return &Cache{p, m}
}

// Get retrives an data structure from cache.
func (c *Cache) Get(ctx context.Context, key string, v interface{}, d time.Duration) error {
	conn, err := c.Pool.GetContext(ctx)
	if err != nil {
		return err
	}

	conn.Send("GET", key)
	conn.Send("EXPIRE", key, int(d.Seconds()))
	conn.Flush()
	data, err := conn.Receive()
	conn.Close()

	b, err := redis.Bytes(data, err)
	if err != nil {
		if err == redis.ErrNil {
			return cache.ErrCacheMiss
		}

		return err
	}

	return c.MarshalUnmarshaler.Unmarshal(bytes.NewReader(b), v)
}

// Set marshals and saves an arbitary go data structure.
func (c *Cache) Set(ctx context.Context, key string, v interface{}, d time.Duration) error {
	buf := rpool.Get().(*bytes.Buffer)
	buf.Reset()
	if err := c.MarshalUnmarshaler.Marshal(buf, v); err != nil {
		rpool.Put(buf)
		return err
	}

	conn, err := c.Pool.GetContext(ctx)
	if err != nil {
		rpool.Put(buf)
		return err
	}
	_, err = conn.Do("SETEX", key, int(d.Seconds()), buf.Bytes())
	conn.Close()
	rpool.Put(buf)
	return err
}

// Renew is used to renew a cache
func (c *Cache) Renew(ctx context.Context, key string, d time.Duration) error {
	conn, err := c.Pool.GetContext(ctx)
	if err != nil {
		return err
	}

	r, err := redis.Int(conn.Do("EXPIRE", key, int(d.Seconds())))
	conn.Close()
	if err != nil {
		return err
	}

	if r == 0 {
		return cache.ErrCacheMiss
	}

	return nil
}

// Evict is used to evict a key from cache.
func (c *Cache) Evict(ctx context.Context, key string) error {
	conn, err := c.Pool.GetContext(ctx)
	if err != nil {
		return err
	}

	count, err := redis.Int(conn.Do("DEL", key))
	conn.Close()
	if err != nil {
		return err
	}

	if count < 1 {
		return cache.ErrCacheMiss
	}

	return err
}

var rpool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}
