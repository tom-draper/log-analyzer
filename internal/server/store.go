package server

import (
	"encoding/json"
	"sync"
)

// Store holds the current dashboard data behind a lock and notifies subscribers
// (SSE clients) whenever it changes, so the dashboard can update live.
type Store struct {
	mu   sync.RWMutex
	data *Data

	subMu sync.Mutex
	subs  map[chan struct{}]struct{}
}

func NewStore(data *Data) *Store {
	return &Store{
		data: data,
		subs: make(map[chan struct{}]struct{}),
	}
}

// snapshotJSON marshals the current data under a read lock.
func (s *Store) snapshotJSON() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return json.Marshal(s.data)
}

// Update mutates the data under a write lock and then notifies subscribers.
func (s *Store) Update(fn func(*Data)) {
	s.mu.Lock()
	fn(s.data)
	s.mu.Unlock()
	s.notify()
}

func (s *Store) subscribe() chan struct{} {
	ch := make(chan struct{}, 1)
	s.subMu.Lock()
	s.subs[ch] = struct{}{}
	s.subMu.Unlock()
	return ch
}

func (s *Store) unsubscribe(ch chan struct{}) {
	s.subMu.Lock()
	delete(s.subs, ch)
	s.subMu.Unlock()
}

// notify signals every subscriber. Sends are non-blocking; a subscriber that
// already has a pending notification simply coalesces the update.
func (s *Store) notify() {
	s.subMu.Lock()
	defer s.subMu.Unlock()
	for ch := range s.subs {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}
