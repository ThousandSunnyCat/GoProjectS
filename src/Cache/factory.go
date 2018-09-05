package Cache

import (
	"runtime"
	"time"
)

// 创建默认的新缓存对象
func New(defaultExpiration, cleanupInterval time.Duration) *Cache {
	items := make(map[string]Item)
	return newCacheWithJanitor(defaultExpiration, cleanupInterval, items)
}

// 创建新缓存对象
func NewFrom(defaultExpiration, cleanupInterval time.Duration, items map[string]Item) *Cache {
	return newCacheWithJanitor(defaultExpiration, cleanupInterval, items)
}

func newCacheWithJanitor(de time.Duration, ci time.Duration, m map[string]Item) *Cache {
	c := newCache(de, m)

	C := &Cache{c}
	if ci > 0 {
		runJanitor(c, ci)
		// C对象回收时 -> 执行stopJanitor方法
		runtime.SetFinalizer(C, stopJanitor)
	}
	return C
}

func newCache(de time.Duration, m map[string]Item) *cache {
	if de == 0 {
		de = -1
	}
	c := &cache {
		defaultExpiration:	de,
		items:				m,
	}

	return c
}