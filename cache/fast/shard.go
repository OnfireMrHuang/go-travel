package fast

import (
	"container/list"
	"go-travel/cache"
	"sync"
)

type entry struct {
	key    string
	value  interface{}
	weight int
	index  int
}

func (e *entry) Len() int {
	return cache.CalcLen(e.value) + 4 + 4
}

type cacheShard struct {
	locker sync.RWMutex

	// 最大存放entry个数
	maxEntries int
	// 当一个entry从缓存中移除是调用该回调函数，默认为nil
	// groupcache中的key是任意的可比较类型，value是interface{}
	onEvicted func(key string, value interface{})

	ll    *list.List
	cache map[string]*list.Element
}

func newCacheShard(maxEntries int, onEvicted func(key string, value interface{})) *cacheShard {
	return &cacheShard{
		maxEntries: maxEntries,
		onEvicted:  onEvicted,
		ll:         list.New(),
		cache:      make(map[string]*list.Element),
	}
}

// get从cache中获取key对应的值，nil表示key不存在
func (c *cacheShard) get(key string) interface{} {
	c.locker.RLocker()
	defer c.locker.RUnlock()

	if e, ok := c.cache[key]; ok {
		c.ll.MoveToBack(e)
		return e.Value.(*entry).value
	}
	return nil
}
