package redis_test

import (
	"testing"
	"time"

	"github.com/srikrsna/go-cache/formats/gob"

	"github.com/srikrsna/go-cache/backends/redis"

	"github.com/srikrsna/go-cache/testsuite"

	rp "github.com/gomodule/redigo/redis"
)

func TestCache(t *testing.T) {
	testsuite.TestCache(t, redis.NewCache(getPool(), gob.Gob{}))
}

func BenchmarkCache(b *testing.B) {
	testsuite.BenchmarkCache(b, redis.NewCache(getPool(), gob.Gob{}))
}

func getPool() *rp.Pool {
	return &rp.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (rp.Conn, error) { return rp.Dial("tcp", ":6379") },
	}
}
