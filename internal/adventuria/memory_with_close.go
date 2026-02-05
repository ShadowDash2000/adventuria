package adventuria

import (
	"adventuria/pkg/cache"
	"iter"
	"sync"
	"time"
)

type Closable interface {
	Close(ctx AppContext)
}

type WithClose[K comparable, V any] interface {
	cache.Cache[K, V]
	Close()
}

type cacheItem[V any] struct {
	value     V
	expiresAt time.Time
}

type MemoryCacheWithClose[K comparable, V any] struct {
	data         sync.Map
	ttl          time.Duration
	count        int
	preventClean bool
}

func NewMemoryCacheWithClose[K comparable, V any](ctx AppContext, ttl time.Duration, preventClean bool) *MemoryCacheWithClose[K, V] {
	cache := &MemoryCacheWithClose[K, V]{
		ttl:          ttl,
		preventClean: preventClean,
	}

	if !preventClean {
		cache.startGC(ctx, 5*time.Minute)
	}

	return cache
}

func (c *MemoryCacheWithClose[K, V]) Set(key K, value V) {
	if _, found := c.data.Load(key); !found {
		c.count++
	}

	c.data.Store(key, cacheItem[V]{
		value:     value,
		expiresAt: time.Now().Add(c.ttl),
	})
}

func (c *MemoryCacheWithClose[K, V]) Get(key K) (V, bool) {
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

func (c *MemoryCacheWithClose[K, V]) deleteInternal(ctx AppContext, key K) {
	if v, found := c.data.LoadAndDelete(key); found {
		c.count--
		item := v.(cacheItem[V])
		if closer, ok := any(item.value).(Closable); ok {
			closer.Close(ctx)
		}
	}
}

func (c *MemoryCacheWithClose[K, V]) Delete(ctx AppContext, key K) {
	c.deleteInternal(ctx, key)
}

func (c *MemoryCacheWithClose[K, V]) GetAll() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		c.data.Range(func(key, value any) bool {
			return yield(key.(K), value.(cacheItem[V]).value)
		})
	}
}

func (c *MemoryCacheWithClose[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		c.data.Range(func(key, _ any) bool {
			return yield(key.(K))
		})
	}
}

func (c *MemoryCacheWithClose[K, V]) Count() int {
	return c.count
}

func (c *MemoryCacheWithClose[K, V]) Clear(ctx AppContext) {
	c.data.Range(func(key, value any) bool {
		c.deleteInternal(ctx, key.(K))
		return true
	})
}

func (c *MemoryCacheWithClose[K, V]) startGC(ctx AppContext, interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			now := time.Now()

			c.data.Range(func(key, value any) bool {
				item := value.(cacheItem[V])
				if now.After(item.expiresAt) {
					c.Delete(ctx, key.(K))
				}
				return true
			})
		}
	}()
}
