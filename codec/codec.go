package cache

import (
	"bytes"
	"sync"
	"time"

	cache "github.com/srikrsna/go-cache"
	"github.com/ugorji/go/codec"
)

type codecCache struct {
	cache.Backend
	pool sync.Pool
	h    codec.Handle
}

func (c *codecCache) Get(key string, v interface{}, d time.Duration) error {
	buf := c.pool.Get().(*bytes.Buffer)
	defer c.pool.Put(buf)

	buf.Reset()

	if err := c.Backend.Get(key, buf, d); err != nil {
		return err
	}

	if buf.Len() <= 0 {
		return cache.ErrCacheMiss
	}

	return codec.NewDecoder(buf, c.h).Decode(v)
}

func (c *codecCache) Set(key string, v interface{}, d time.Duration) error {
	buf := c.pool.Get().(*bytes.Buffer)
	defer c.pool.Put(buf)

	buf.Reset()
	if err := codec.NewEncoder(buf, c.h).Encode(v); err != nil {
		return err
	}

	return c.Backend.Set(key, buf.Bytes(), d)
}

// NewJSONCache ...
func NewJSONCache(backend cache.Backend) cache.Cache {
	return NewCodecCache(backend, new(codec.JsonHandle))
}

// NewMsgPackCache ...
func NewMsgPackCache(backend cache.Backend) cache.Cache {
	return NewCodecCache(backend, new(codec.MsgpackHandle))
}

// NewCborCache ...
func NewCborCache(backend cache.Backend) cache.Cache {
	return NewCodecCache(backend, new(codec.CborHandle))
}

// NewBincCache ...
func NewBincCache(backend cache.Backend) cache.Cache {
	return NewCodecCache(backend, new(codec.BincHandle))
}

// NewCodecCache ...
func NewCodecCache(b cache.Backend, h codec.Handle) cache.Cache {
	return &codecCache{
		Backend: b,
		pool: sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
		h: h,
	}
}
