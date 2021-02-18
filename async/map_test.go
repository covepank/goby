package async_test

import (
	"strconv"
	"sync"
	"testing"

	"github.com/sanbsy/goby/async"
)

func BenchmarkConMap_Set(b *testing.B) {
	cmap := async.NewConMap()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cmap.Set("key", "value")
		}
	})
}

func BenchmarkSyncMap_Set(b *testing.B) {
	smap := sync.Map{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			smap.Store("key", "value")
		}
	})
}

func BenchmarkConMap_SetValues(b *testing.B) {
	data := make(map[string]interface{}, 1024)
	for i := 0; i < 1024; i++ {
		data[strconv.Itoa(i)] = i
	}

	cmap := async.NewConMap()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cmap.SetValues(data)
		}
	})
}

func BenchmarkConMap_SetValues2(b *testing.B) {
	data := make(map[string]interface{}, 1024)
	for i := 0; i < 1024; i++ {
		data[strconv.Itoa(i)] = i
	}

	cmap := async.NewConMap()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for key, value := range data {
				cmap.Set(key, value)
			}
		}
	})
}
