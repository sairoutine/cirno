package main

import (
	"testing"

	"github.com/bradfitz/gomemcache/memcache"
)

var mc = memcache.New("localhost:11212")

func BenchmarkSingle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DoAnyThing()
	}
}

func BenchmarkParallel(b *testing.B) {
	b.SetParallelism(5)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			DoAnyThing()
		}
	})
}

func DoAnyThing() {
	mc.Get("id")
}
