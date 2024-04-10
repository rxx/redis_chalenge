package main

import "sync"

type Store struct {
	data map[string]string
	mu   sync.RWMutex
}

func (r *Store) Get(key string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	val, ok := r.data[key]
	return val, ok
}

func (r *Store) Set(key, value string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[key] = value
}

func NewStore() *Store {
	return &Store{data: make(map[string]string)}
}
