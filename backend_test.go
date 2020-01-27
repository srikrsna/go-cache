package cache_test

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mediocregopher/radix/v3"

	cache "github.com/srikrsna/go-cache"
	mem "github.com/srikrsna/go-cache/memory"
	rad "github.com/srikrsna/go-cache/radix"
)

func TestMemoryBackend(t *testing.T) {
	RunTestsForBackend(mem.NewInMemoryBackend(), t)
}

func TestRedisRadixBackend(t *testing.T) {
	pool, err := radix.NewPool("tcp", "127.0.0.1:6379", 10)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	RunTestsForBackend(rad.NewRadixV3Backend(pool), t)
}

func RunTestsForBackend(c cache.Backend, t *testing.T) {
	getRandomObject := func() *bytes.Buffer {
		buf := new(bytes.Buffer)
		buf.WriteString("value")
		return buf
	}

	key := uuid.New().String()

	randObj := getRandomObject()

	if err := c.Set(key, randObj.Bytes(), time.Second); err != nil {
		t.Errorf("error while setting a random value: %v", err)
	}

	gotRandObj := new(bytes.Buffer)
	if err := c.Get(key, gotRandObj, time.Millisecond); err != nil {
		t.Errorf("error while getting value from cache %v", err)
	}

	if !reflect.DeepEqual(gotRandObj.Bytes(), randObj.Bytes()) {
		t.Errorf("Cache SET and GET Mismatch, expected: %v, got: %v", randObj, gotRandObj)
	}

	time.Sleep(2 * time.Millisecond)

	gotRandObj.Reset()
	if err := c.Get(key, gotRandObj, 2*time.Second); err != cache.ErrCacheMiss && gotRandObj.Len() > 0 {
		t.Errorf("key expired should have gotten: %v, got: %v", cache.ErrCacheMiss, err)
	}

	if err := c.Set(key, randObj.Bytes(), 2*time.Second); err != nil {
		t.Errorf("error while setting a random value: %v", err)
	}

	if err := c.Renew(key, time.Second); err != nil {
		t.Errorf("unexpected error while renewing key: %v", err)
	}

	time.Sleep(3 * time.Millisecond)

	gotRandObj.Reset()
	if err := c.Get(key, gotRandObj, time.Second); err != nil || gotRandObj.Len() == 0 {
		t.Errorf("error while getting value from cache %v", err)
	}

	if err := c.Delete(key); err != nil {
		t.Errorf("unexpected error while deleting key: %v", err)
	}

	gotRandObj.Reset()
	if err := c.Get(key, gotRandObj, time.Second); err != cache.ErrCacheMiss && gotRandObj.Len() > 0 {
		t.Errorf("key deleted should have gotten: %v, got: %v", cache.ErrCacheMiss, err)
	}
}
