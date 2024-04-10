package main

import (
	"sync"
	"time"
)

type kvalue struct {
	value string
	time  time.Time
}

type Store struct {
	data map[string]kvalue
	mu   sync.RWMutex
}

func (v *Store) Get(key string) (string, bool) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	val, ok := v.data[key]

	if ok && val.time.Before(time.Now()) && !val.time.IsZero() {
		return "", false
	}

	return val.value, ok
}

func (v *Store) Set(key, value string, expireDuration time.Duration) {
	v.mu.Lock()
	defer v.mu.Unlock()

	var expiredTime time.Time
	if expireDuration > 0 {
		expiredTime = time.Now().Add(expireDuration)
	}
	v.data[key] = kvalue{value: value, time: expiredTime}
}

func NewStore() *Store {
	return &Store{data: make(map[string]kvalue)}
}
