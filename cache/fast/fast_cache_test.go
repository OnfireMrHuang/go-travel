package fast

import (
	"testing"
)

func BenchmarkTourFastCacheSetParallel(b *testing.B) {
	//cache := NewFastCache(b.N,maxEntrySize,nil)
	//rand.Seed(time.Now().Unix())
	//
	//b.RunParallel(func(pb *testing.PB) {
	//	id := rand.Intn(1000)
	//	counter := 0
	//	for pb.Next() {
	//		cache.Set(parallelKey(id,counter),value())
	//		counter = counter + 1
	//	}
	//
	//})
	//
}
