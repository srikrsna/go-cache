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
	if memory.NewCache().Get(context.Background(), "random", 1234, time.Second) == nil {
		t.Errorf("get should only accept pointers")
	}
}

func BenchmarkCache(b *testing.B) {
	testsuite.BenchmarkCache(b, memory.NewCache())
}
