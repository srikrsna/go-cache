package testsuite

import (
	"context"
	"encoding/gob"
	"reflect"
	"testing"
	"time"

	"github.com/fortytw2/leaktest"

	"github.com/google/gofuzz"
	"github.com/srikrsna/go-cache"
)

func init() {
	gob.Register(&Data{})
}

// Data is the type used during the test suite.
type Data struct {
	Name     string
	BirthDay time.Time
	Phone    string
	Siblings int
	Spouse   bool
	Money    float64
}

// TestCache runs all test cases that Cache should to satify
func TestCache(t *testing.T, c cache.Cache) {
	defer leaktest.Check(t)()

	ctx := context.Background()

	fuzzer := fuzz.New()

	var (
		key   string
		v, sv Data
	)

	fuzzer.Fuzz(&v)
	key = "key"

	v.BirthDay = time.Now().UTC()

	if err := c.Set(ctx, key, &v, time.Second); err != nil {
		t.Fatalf("error while setting object: %v", err)
	}

	if err := c.Get(ctx, key, &sv, time.Second); err != nil {
		t.Fatalf("error getting data from cache: %v", err)
	}

	if !reflect.DeepEqual(v, sv) {
		t.Errorf("data structure that was set should be the same as the one that is retreived")
	}

	if err := c.Evict(ctx, key); err != nil {
		t.Fatalf("error evicting cache: %v", err)
	}

	if c.Get(ctx, key, &sv, time.Second) != cache.ErrCacheMiss {
		t.Errorf("get after evict should return ErrCacheMiss")
	}

	if !reflect.DeepEqual(v, sv) {
		t.Errorf("get returning ErrCacheMiss should not modify passed in v")
	}

	if err := c.Set(ctx, key, &v, time.Second); err != nil {
		t.Errorf("error while adding data to cache: %v", err)
	}

	time.Sleep(50*time.Millisecond + time.Second) // 50 millisecond buffer It should be acceptable

	if c.Get(ctx, key, &sv, time.Second) != cache.ErrCacheMiss {
		t.Errorf("cache should be expired according to the time provided in set")
	}

	if err := c.Set(ctx, key, &v, time.Second); err != nil {
		t.Errorf("error while adding data to cache: %v", err)
	}

	if err := c.Set(ctx, key, &v, time.Second); err != nil {
		t.Errorf("error while setting cache for the same key: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	if err := c.Renew(ctx, key, time.Second); err != nil {
		t.Errorf("unable to renew cache")
	}

	if err := c.Get(ctx, key, &sv, time.Second); err != nil {
		t.Fatalf("error getting data from cache: %v", err)
	}

	time.Sleep(50*time.Millisecond + time.Second) // 50 millisecond buffer It should be acceptable

	if c.Get(ctx, key, &sv, time.Second) != cache.ErrCacheMiss {
		t.Errorf("cache should be expired according to the time provided in set")
	}

	if c.Renew(ctx, key, time.Second) != cache.ErrCacheMiss {
		t.Errorf("renew should return ErrCacheMiss if key is missing")
	}

	if c.Evict(ctx, key) != cache.ErrCacheMiss {
		t.Errorf("evict should return ErrCacheMiss if key is missing")
	}
}

// BenchmarkCache benchmarks Caches
func BenchmarkCache(b *testing.B, c cache.Cache) {
	var (
		key string
		v   Data
	)

	fuzzer := fuzz.New()
	ctx := context.Background()

	fuzzer.Fuzz(&key)
	fuzzer.Fuzz(&v)
	if c.Set(ctx, key, v, time.Second) != nil {
		b.Errorf("unable to start benchmark: cannot run set. Please run the tests first")
	}
	b.Run("Parallel Get", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			var sv Data
			for pb.Next() {
				c.Get(ctx, key, &sv, time.Second)
			}
		})
	})
	b.Run("Parallel Set", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				c.Set(ctx, key, &v, time.Second)
			}
		})
	})
}
