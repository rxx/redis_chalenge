package main

import (
	"testing"
	"time"
)

func TestStore(t *testing.T) {
	testStore := NewStore()

	if _, ok := testStore.Get("key"); ok {
		t.Errorf("Error expected because key is not found: %v", testStore)
	}

	testStore.Set("key", "val", 0)

	if val, ok := testStore.Get("key"); !ok || val != "val" {
		t.Errorf("Error on getting value %v", testStore)
	}

	if !testStore.data["key"].time.IsZero() {
		t.Error("Expected 'key' time is zero")
	}

	const expireDuration = 10 * time.Millisecond

	testStore.Set("key", "val", expireDuration)

	// it must exists just after setting. Could you increase duration please? ^_^
	if val, ok := testStore.Get("key"); !ok || val != "val" {
		t.Errorf("Error on getting value with duration %v", testStore)
	}

	if testStore.data["key"].time.IsZero() {
		t.Error("Expected 'key' time set")
	}

	time.Sleep(expireDuration)

	if _, ok := testStore.Get("key"); ok {
		t.Errorf("It should be expired after waiting %v. Isn't it? Check: %v", expireDuration, testStore)
	}
}
