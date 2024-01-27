package common

import (
	"context"
	"fmt"
	"sync"
)

var ErrLockResourceLocked = fmt.Errorf("resource locked")

type LockUnlockFunc = func()

func NewLockStore[T comparable]() *LockStore[T] {
	return &LockStore[T]{
		mu:    sync.Mutex{},
		locks: make(map[T]chan struct{}),
	}
}

// LockStore handles concurrent resource locking.
type LockStore[T comparable] struct {
	mu    sync.Mutex
	locks map[T]chan struct{}
}

func (s *LockStore[T]) GetLock(slug T) bool {
	s.mu.Lock()
	_, found := s.locks[slug]
	s.mu.Unlock()

	return found
}

func (s *LockStore[T]) ListLock(slugs ...T) []bool {
	s.mu.Lock()
	res := make([]bool, 0, len(slugs))
	for _, slug := range slugs {
		_, found := s.locks[slug]
		res = append(res, found)
	}
	s.mu.Unlock()

	return res
}

func (s *LockStore[T]) ListSlug() []T {
	s.mu.Lock()
	res := make([]T, 0, len(s.locks))
	for slug := range s.locks {
		res = append(res, slug)
	}
	s.mu.Unlock()

	return res
}

// Lock blocks until it is able to lock the resource.
func (s *LockStore[T]) Lock(ctx context.Context, slug T) (LockUnlockFunc, error) {
	for {
		s.mu.Lock()
		r, found := s.locks[slug]
		if found {
			// Already locked
			s.mu.Unlock()

			// Wait for unlock
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case _, ok := <-r:
				if !ok {
					// Previous lock removed
					continue
				}
			}

			// Use previous lock
			return s.unlockFn(slug), nil
		}

		// Create new lock
		s.locks[slug] = make(chan struct{})
		s.mu.Unlock()
		return s.unlockFn(slug), nil
	}
}

// TryLock tries to lock the resource if it is available.
func (s *LockStore[T]) TryLock(slug T) (LockUnlockFunc, error) {
	s.mu.Lock()
	_, found := s.locks[slug]
	if found {
		s.mu.Unlock()
		return nil, fmt.Errorf("%w: %v", ErrLockResourceLocked, slug)
	}
	s.locks[slug] = make(chan struct{})
	s.mu.Unlock()

	return s.unlockFn(slug), nil
}

func (s *LockStore[T]) unlockFn(slug T) LockUnlockFunc {
	return sync.OnceFunc(func() {
		s.mu.Lock()
		// Assume lock exists
		r := s.locks[slug]
		select {
		case r <- struct{}{}:
		// Gave lock
		default:
			// Remove lock
			close(r)
			delete(s.locks, slug)
		}
		s.mu.Unlock()
	})
}
