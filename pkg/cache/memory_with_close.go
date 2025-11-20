package cache

import "time"

type MemoryCacheWithClose[K comparable, V any] struct {
	MemoryCache[K, V]
}

func NewMemoryCacheWithClose[K comparable, V any](ttl time.Duration, preventClean bool) Cache[K, V] {
	cache := &MemoryCacheWithClose[K, V]{
		MemoryCache: MemoryCache[K, V]{
			ttl:          ttl,
			preventClean: preventClean,
		},
	}

	if !preventClean {
		cache.startGC(5 * time.Minute)
	}

	return cache
}

func (c *MemoryCacheWithClose[K, V]) deleteInternal(key K) {
	if v, found := c.data.LoadAndDelete(key); found {
		c.count--
		item := v.(cacheItem[V])
		if closer, ok := any(item.value).(Closable); ok {
			closer.Close()
		}
	}
}

func (c *MemoryCacheWithClose[K, V]) Delete(key K) {
	c.deleteInternal(key)
}

func (c *MemoryCacheWithClose[K, V]) Clear() {
	c.data.Range(func(key, value any) bool {
		c.deleteInternal(key.(K))
		return true
	})
}

func (c *MemoryCacheWithClose[K, V]) startGC(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			now := time.Now()

			c.data.Range(func(key, value any) bool {
				item := value.(cacheItem[V])
				if now.After(item.expiresAt) {
					c.Delete(key.(K))
				}
				return true
			})
		}
	}()
}
