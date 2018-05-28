package cache

import (
	"io"
	"sync"
	"time"

	cache "github.com/srikrsna/go-cache"
)

type inMemoryBackend struct {
	sync.Mutex

	data   map[string][]byte
	expiry map[string]*time.Timer
}

func (c *inMemoryBackend) Get(key string, w io.Writer, d time.Duration) error {
	c.Lock()
	defer c.Unlock()

	c.expire(key, d)

	if val, ok := c.data[key]; ok {
		_, err := w.Write(val)
		return err
	}

	return cache.ErrCacheMiss
}

func (c *inMemoryBackend) Set(key string, r []byte, d time.Duration) error {
	c.Lock()
	defer c.Unlock()

	c.data[key] = r
	c.expire(key, d)

	return nil
}

func (c *inMemoryBackend) Renew(key string, d time.Duration) error {
	c.Lock()
	defer c.Unlock()

	c.expire(key, d)

	return nil
}

func (c *inMemoryBackend) Delete(key string) error {
	c.Lock()
	defer c.Unlock()

	c.expiry[key].Stop()

	delete(c.data, key)
	delete(c.expiry, key)

	return nil
}

func (c *inMemoryBackend) expire(key string, d time.Duration) {
	if val, ok := c.expiry[key]; ok {
		val.Stop()
	}

	c.expiry[key] = time.AfterFunc(d, func() {
		c.Lock()
		defer c.Unlock()

		delete(c.data, key)
		delete(c.expiry, key)
	})
}

// NewInMemoryBackend ...
func NewInMemoryBackend() cache.Backend {
	return &inMemoryBackend{
		data:   map[string][]byte{},
		expiry: map[string]*time.Timer{},
	}
}
