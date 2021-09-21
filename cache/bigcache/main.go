package bigcache

import (
	"github.com/allegro/bigcache/v2"
	"log"
	"time"
)

func main() {
	cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	if err != nil {
		log.Println(err)
		return
	}

	entry, err := cache.Get("my-unique-key")
	if err != nil {
		log.Println(err)
		return
	}
	if entry == nil {
		// 从缓存中没有获取到，则从数据库获取，然后设置缓存
		entry = []byte("value")
		cache.Set("my-unique-key", entry)
	}
	log.Println(string(entry))
}
