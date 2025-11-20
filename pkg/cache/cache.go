package cache

import (
	"iter"
)

type Cache[K comparable, V any] interface {
	Set(K, V)
	Get(K) (V, bool)
	Delete(K)
	GetAll() iter.Seq2[K, V]
	Keys() iter.Seq[K]
	Count() int
	Clear()
}

type Closable interface {
	Close()
}

type WithClose[K comparable, V any] interface {
	Cache[K, V]
	Close()
}
