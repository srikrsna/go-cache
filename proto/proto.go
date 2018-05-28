package cache

import (
	"bytes"
	"errors"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	cache "github.com/srikrsna/go-cache"
)

type protoCache struct {
	cache.Backend
	pool      sync.Pool
	protoPool sync.Pool
}

func (c *protoCache) Get(key string, v interface{}, d time.Duration) error {
	pb, ok := v.(proto.Message)
	if !ok {
		return errors.New("not a proto message")
	}

	buf := c.pool.Get().(*bytes.Buffer)
	defer c.pool.Put(buf)

	buf.Reset()

	if err := c.Backend.Get(key, buf, d); err != nil {
		return err
	}

	if buf.Len() <= 0 {
		return cache.ErrCacheMiss
	}

	return proto.Unmarshal(buf.Bytes(), pb)
}

func (c *protoCache) Set(key string, v interface{}, d time.Duration) error {
	pb, ok := v.(proto.Message)
	if !ok {
		return errors.New("not a proto message")
	}

	buf := c.protoPool.Get().(*proto.Buffer)
	defer c.protoPool.Put(buf)

	buf.Reset()
	if err := buf.Marshal(pb); err != nil {
		return err
	}

	return c.Backend.Set(key, buf.Bytes(), d)
}

// NewProtoBufferCache ...
func NewProtoBufferCache(b cache.Backend) cache.Cache {
	return &protoCache{
		Backend: b,
		pool: sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
		protoPool: sync.Pool{
			New: func() interface{} {
				return proto.NewBuffer([]byte{})
			},
		},
	}
}
