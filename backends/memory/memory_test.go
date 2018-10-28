package memory_test

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/srikrsna/go-cache/backends/memory"
	"github.com/srikrsna/go-cache/testsuite"
)

func TestCache(t *testing.T) {
	testsuite.TestCache(t, memory.NewCache())
}

func TestPointerInput(t *testing.T) {
	c := memory.NewCache()
	Convey("If a value is set", t, func() {
		So(c.Set(context.Background(), "random", 1234, time.Second), ShouldBeNil)
		Convey("It should return an error if a pointer is not provided", func() {
			So(c.Get(context.Background(), "random", 1234, time.Second), ShouldEqual, memory.ErrNotPtr)
		})
	})
}

func BenchmarkCache(b *testing.B) {
	testsuite.BenchmarkCache(b, memory.NewCache())
}
