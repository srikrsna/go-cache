package cache_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/duration"

	"github.com/google/uuid"
	cache "github.com/srikrsna/go-cache"
	codec "github.com/srikrsna/go-cache/codec"
	mem "github.com/srikrsna/go-cache/memory"
	p "github.com/srikrsna/go-cache/proto"
)

func TestMsgPackCache(t *testing.T) {
	c := codec.NewMsgPackCache(mem.NewInMemoryBackend())

	RunTestsForCache(c, t)
}

func TestJSONCache(t *testing.T) {
	c := codec.NewJSONCache(mem.NewInMemoryBackend())

	RunTestsForCache(c, t)
}

func TestCBORCache(t *testing.T) {
	c := codec.NewCborCache(mem.NewInMemoryBackend())

	RunTestsForCache(c, t)
}

func TestBincCache(t *testing.T) {
	c := codec.NewBincCache(mem.NewInMemoryBackend())

	RunTestsForCache(c, t)
}

func TestProtoCache(t *testing.T) {
	c := p.NewProtoBufferCache(mem.NewInMemoryBackend())

	RunTestsForCache(c, t)
}

func RunTestsForCache(c cache.Cache, t *testing.T) {
	getRandomObject := func() interface{} {
		return &duration.Duration{Seconds: 20, Nanos: 20}
	}

	key := uuid.New().String()

	randObj := getRandomObject()

	if err := c.Set(key, randObj, time.Second); err != nil {
		t.Errorf("error while setting a random value: %v", err)
	}

	gotRandObj := duration.Duration{}
	if err := c.Get(key, &gotRandObj, time.Millisecond); err != nil {
		t.Errorf("error while getting value from cache %v", err)
	}

	if !reflect.DeepEqual(&gotRandObj, randObj) {
		t.Errorf("Cache SET and GET Mismatch, expected: %v, got: %v", randObj, gotRandObj)
	}

	time.Sleep(2 * time.Millisecond)

	if err := c.Get(key, &gotRandObj, time.Second); err != cache.ErrCacheMiss {
		t.Errorf("key expired should have gotten: %v, got: %v", cache.ErrCacheMiss, err)
	}

	if err := c.Set(key, randObj, 2*time.Millisecond); err != nil {
		t.Errorf("error while setting a random value: %v", err)
	}

	if err := c.Renew(key, time.Second); err != nil {
		t.Errorf("unexpected error while renewing key: %v", err)
	}

	time.Sleep(3 * time.Millisecond)

	if err := c.Get(key, &gotRandObj, time.Second); err != nil {
		t.Errorf("error while getting value from cache %v", err)
	}

	if err := c.Delete(key); err != nil {
		t.Errorf("unexpected error while deleting key: %v", err)
	}

	if err := c.Get(key, &gotRandObj, time.Second); err != cache.ErrCacheMiss {
		t.Errorf("key deleted should have gotten: %v, got: %v", cache.ErrCacheMiss, err)
	}
}

func BenchmarkCache_Proto(b *testing.B) {
	c := p.NewProtoBufferCache(mem.NewInMemoryBackend())
	RunBenchmarkForCache(c, b)
}

func BenchmarkCache_Json(b *testing.B) {
	c := codec.NewJSONCache(mem.NewInMemoryBackend())
	RunBenchmarkForCache(c, b)
}

func BenchmarkCache_Msgpack(b *testing.B) {
	c := codec.NewMsgPackCache(mem.NewInMemoryBackend())
	RunBenchmarkForCache(c, b)
}

func BenchmarkCache_Cbor(b *testing.B) {
	c := codec.NewCborCache(mem.NewInMemoryBackend())
	RunBenchmarkForCache(c, b)
}

func BenchmarkCache_Binc(b *testing.B) {
	c := codec.NewBincCache(mem.NewInMemoryBackend())
	RunBenchmarkForCache(c, b)
}

func RunBenchmarkForCache(c cache.Cache, b *testing.B) {
	key := "key"
	obj := &duration.Duration{Seconds: 20, Nanos: 20}

	b.Run("Set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			c.Set(key, obj, time.Second)
		}
	})

	b.Run("Set Parallel", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				c.Set(key, obj, time.Second)
			}
		})
	})

	b.Run("Get", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			c.Get(key, obj, time.Second)
		}
	})

	b.Run("Get Parallel", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				c.Get(key, obj, time.Second)
			}
		})
	})
}
