package locker

import "sync"

type Locker[K comparable] struct {
	m  map[K]struct{}
	mu sync.Mutex
}

func New[K comparable]() *Locker[K] {
	return &Locker[K]{
		m: make(map[K]struct{}),
	}
}

func (l *Locker[K]) TryLock(key K) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if _, ok := l.m[key]; ok {
		return false
	}
	l.m[key] = struct{}{}
	return true
}

func (l *Locker[K]) Unlock(key K) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.m, key)
}
