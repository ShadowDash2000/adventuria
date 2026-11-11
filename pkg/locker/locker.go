package locker

import (
	"context"
	"sync"
	"time"
)

type Locker[K comparable] struct {
	m  map[K]struct{}
	mu sync.Mutex
}

func New[K comparable]() *Locker[K] {
	return &Locker[K]{
		m: make(map[K]struct{}),
	}
}

func (l *Locker[K]) Lock(ctx context.Context, key K) error {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		if err := ctx.Err(); err != nil {
			return err
		}

		if l.TryLock(key) {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
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

func (l *Locker[K]) HasAny() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.m) > 0
}
