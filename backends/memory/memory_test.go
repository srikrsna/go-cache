package memory_test

import (
	"context"
	"testing"
	"time"

	"github.com/srikrsna/go-cache/backends/memory"
	"github.com/srikrsna/go-cache/testsuite"
)

func TestCache(t *testing.T) {
	testsuite.TestCache(t, memory.NewCache())
}

func TestPointerInput(t *testing.T) {
	c := memory.NewCache()
	c.Set(context.Background(), "random", 1234, time.Second)
	if memory.NewCache().Get(context.Background(), "random", 1234, time.Second) != memory.ErrNotPtr {
		t.Errorf("get should only accept pointers")
	}
}

func BenchmarkCache(b *testing.B) {
	testsuite.BenchmarkCache(b, memory.NewCache())
}
