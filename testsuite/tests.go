package testsuite

import (
	"context"
	"encoding/gob"
	"testing"
	"time"

	"github.com/google/gofuzz"
	. "github.com/smartystreets/goconvey/convey"
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

	fuzzer := fuzz.New()

	var (
		key  string
		data Data
		ctx  context.Context
	)

	ctx = context.Background()

	fuzzer.Fuzz(&data)
	fuzzer.Fuzz(&key)

	Convey("Once a data has been added to the cache with a unique key and specific expiration", t, func() {
		So(c.Set(ctx, key, &data, time.Second), ShouldBeNil)

		Convey("It should be retrieved as is with the given unique key under the expiration limit", func() {
			var sv Data
			So(c.Get(ctx, key, &sv, time.Second), ShouldBeNil)
			So(sv, ShouldResemble, data)
		})

		Convey("Once the time has expired it should no longer be accessible", func() {
			time.Sleep(40*time.Millisecond + time.Second)
			var sv Data
			So(c.Get(ctx, key, &sv, time.Second), ShouldEqual, cache.ErrCacheMiss)
		})

	})

	Convey("Once a data has been added to cache and evicted", t, func() {
		So(c.Set(ctx, key, &data, time.Second), ShouldBeNil)
		So(c.Evict(ctx, key), ShouldBeNil)

		Convey("It should no longer be accessible returning a cache miss error", func() {
			var sv Data
			og := sv
			So(c.Get(ctx, key, &sv, time.Second), ShouldEqual, cache.ErrCacheMiss)
			So(sv, ShouldResemble, og)
		})

		Convey("It should no longer be renewable returning a cache miss error", func() {
			So(c.Renew(ctx, key, time.Second), ShouldEqual, cache.ErrCacheMiss)
		})

		Convey("If tried to evict it it should return a cache miss", func() {
			So(c.Evict(ctx, key), ShouldEqual, cache.ErrCacheMiss)
		})
	})

	Convey("Once a data has been added to cache", t, func() {
		So(c.Set(ctx, key, &data, time.Second), ShouldBeNil)
		Convey("If it is reset with new data and new expiry", func() {
			var nd Data
			fuzzer.Fuzz(&nd)
			So(c.Set(ctx, key, &nd, 3*time.Second), ShouldBeNil)

			Convey("The new value should be returned upon a successful get", func() {
				var sv Data
				So(c.Get(ctx, key, &sv, time.Millisecond), ShouldBeNil)
				So(&sv, ShouldResemble, &nd)
			})

			Convey("The new expiry should be honoured", func() {
				time.Sleep(2 * time.Second)
				var sv Data
				So(c.Get(ctx, key, &sv, time.Millisecond), ShouldBeNil)
				So(&sv, ShouldResemble, &nd)
				time.Sleep(40*time.Millisecond + time.Second)
				So(c.Get(ctx, key, &sv, time.Second), ShouldEqual, cache.ErrCacheMiss)
			})

		})

	})

	Convey("Once a data has been added", t, func() {
		So(c.Set(ctx, key, &data, time.Second), ShouldBeNil)
		Convey("If it has been renewed", func() {
			So(c.Renew(ctx, key, 2*time.Second), ShouldBeNil)
			Convey("It's expiry should have been updated", func() {
				time.Sleep(time.Second)
				var sv Data
				So(c.Get(ctx, key, &sv, time.Millisecond), ShouldBeNil)
				time.Sleep(40*time.Millisecond + time.Second)
				So(c.Get(ctx, key, &sv, time.Second), ShouldEqual, cache.ErrCacheMiss)
			})
		})
	})

	Convey("Trying to retrieve an unset value should return cache miss", t, func() {
		var sv Data
		So(c.Get(ctx, key, &sv, time.Second), ShouldEqual, cache.ErrCacheMiss)
	})

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
