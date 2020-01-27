package cache

import (
	"io"
	"strconv"
	"time"

	"github.com/mediocregopher/radix/v3"
	cache "github.com/srikrsna/go-cache"
)

type radixV3 struct {
	radix.Client
}

// NewRadixV3Backend returns a redis backend powered by radix.v3 client library
func NewRadixV3Backend(c radix.Client) cache.Backend {
	return &radixV3{c}
}

func (c *radixV3) Get(key string, w io.Writer, d time.Duration) error {
	p := radix.Pipeline(
		radix.Cmd(w, "GET", key),
		radix.Cmd(nil, "PEXPIRE", key, strconv.Itoa(int(d/time.Millisecond))),
	)
	return c.Do(p)
}

func (c *radixV3) Set(key string, b []byte, d time.Duration) error {
	return c.Do(radix.Cmd(nil, "PSETEX", key, strconv.Itoa(int(d/time.Millisecond)), string(b)))
}

func (c *radixV3) Renew(key string, d time.Duration) error {
	return c.Do(radix.Cmd(nil, "PEXPIRE", key, strconv.Itoa(int(d/time.Millisecond))))
}

func (c *radixV3) Delete(key string) error {
	return c.Do(radix.Cmd(nil, "DEL", key))
}
