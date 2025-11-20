package cache

import (
	"iter"
	"sync"
	"time"
)

type MemoryCache[K comparable, V any] struct {
	data         sync.Map
	ttl          time.Duration
	count        int
	preventClean bool
}

type cacheItem[V any] struct {
	value     V
	expiresAt time.Time
}

func NewMemoryCache[K comparable, V any](ttl time.Duration, preventClean bool) Cache[K, V] {
	cache := &MemoryCache[K, V]{
		ttl:          ttl,
		preventClean: preventClean,
	}

	if !preventClean {
		cache.startGC(5 * time.Minute)
	}

	return cache
}

func (c *MemoryCache[K, V]) Set(key K, value V) {
	if _, found := c.data.Load(key); !found {
		c.count++
	}

	c.data.Store(key, cacheItem[V]{
		value:     value,
		expiresAt: time.Now().Add(c.ttl),
	})
}

func (c *MemoryCache[K, V]) Get(key K) (V, bool) {
	item, found := c.data.Load(key)
	if !found {
		var zero V
		return zero, false
	}

	cachedItem := item.(cacheItem[V])
	if !c.preventClean && time.Now().After(cachedItem.expiresAt) {
		c.data.Delete(key)
		var zero V
		return zero, false
	}

	// Refresh expiration on successful access
	cachedItem.expiresAt = time.Now().Add(c.ttl)
	c.data.Store(key, cachedItem)

	return cachedItem.value, true
}

func (c *MemoryCache[K, V]) Delete(key K) {
	if _, found := c.data.LoadAndDelete(key); found {
		c.count--
	}
}

func (c *MemoryCache[K, V]) GetAll() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		c.data.Range(func(key, value any) bool {
			return yield(key.(K), value.(cacheItem[V]).value)
		})
	}
}

func (c *MemoryCache[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		c.data.Range(func(key, _ any) bool {
			return yield(key.(K))
		})
	}
}

func (c *MemoryCache[K, V]) Count() int {
	return c.count
}

func (c *MemoryCache[K, V]) Clear() {
	c.data.Clear()
	c.count = 0
}

func (c *MemoryCache[K, V]) startGC(interval time.Duration) {
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
