package cache

import (
	"sync"
	"time"
)

type MemoryCache[T any] struct {
	data sync.Map
	ttl  time.Duration
}

type cacheItem[T any] struct {
	value     T
	expiresAt time.Time
}

func NewMemoryCache[T any](ttl time.Duration, preventClear bool) *MemoryCache[T] {
	cache := &MemoryCache[T]{ttl: ttl}
	if !preventClear {
		cache.startGC(5 * time.Minute)
	}
	return cache
}

func (c *MemoryCache[T]) Set(key string, value T) {
	c.data.Store(key, cacheItem[T]{
		value:     value,
		expiresAt: time.Now().Add(c.ttl),
	})
}

func (c *MemoryCache[T]) Get(key string) (T, bool) {
	item, found := c.data.Load(key)
	if !found {
		var zero T
		return zero, false
	}

	cachedItem := item.(cacheItem[T])
	if time.Now().After(cachedItem.expiresAt) {
		c.data.Delete(key)
		var zero T
		return zero, false
	}

	return cachedItem.value, true
}

func (c *MemoryCache[T]) Delete(key string) {
	c.data.Delete(key)
}

func (c *MemoryCache[T]) startGC(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			now := time.Now()

			c.data.Range(func(key, value interface{}) bool {
				cacheItem := value.(cacheItem[T])
				if now.After(cacheItem.expiresAt) {
					c.data.Delete(key)
				}
				return true
			})
		}
	}()
}
