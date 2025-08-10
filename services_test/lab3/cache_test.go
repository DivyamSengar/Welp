package services_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	cache "cse190-welp/applications"
	"cse190-welp/proto/mycache"
)

func constructCache(capacity int) (cache.Cache, error) {
	policy := os.Getenv("POLICY")

	switch policy {
	case "FIFO":
		return cache.NewFIFOCacheApp(capacity), nil
	case "Random":
		return cache.NewRandomCacheApp(capacity), nil
	case "LRU":
		return cache.NewLRUCacheApp(capacity), nil
	case "LFU":
		return cache.NewLFUCacheApp(capacity), nil
	default:
		return nil, fmt.Errorf("Unrecognized cache policy %s", policy)
	}
}

func TestCacheAppImplementation(t *testing.T) {
	c, err := constructCache(5)
	if err != nil {
		t.Fatal(err)
	}

	// Insert and retrieve items and check cache length
	for i := 0; i < 5; i++ {
		key_str := fmt.Sprintf("key%d", i)
		value_str := fmt.Sprintf("value%d", i)
		c.Set(&mycache.CacheItem{Key: key_str, Value: []byte(value_str)})

		val, err := c.Get(key_str)
		if err != nil || string(val.Value) != value_str {
			t.Errorf("Expected %s, got '%s'", value_str, val.Value)
		}

		if c.Len() != i+1 {
			t.Errorf("Expected cache length %d, got %d", i+1, c.Len())
		}
	}

	// Delete items in a different order and check cache length
	ids := []int{3, 0, 2, 4, 1}
	for index, i := range ids {
		key_str := fmt.Sprintf("key%d", i)

		err := c.Delete(key_str)
		if err != nil {
			t.Errorf("Error deleting '%s': %v", key_str, err)
		}

		if c.Len() != 4-index {
			t.Errorf("Expected cache length %d, got %d", 4-index, c.Len())
		}
	}

	// Clear cache
	c.Clear()
	if c.Len() != 0 {
		t.Error("Cache should be empty after clear")
	}
}

func TestCacheAppConcurrency(t *testing.T) {
	cacheImpl, err := constructCache(3)
	if err != nil {
		t.Fatal(err)
	}

	// Concurrent insertion
	go func() {
		cacheImpl.Set(&mycache.CacheItem{Key: "key1", Value: []byte("value1")})
	}()
	go func() {
		cacheImpl.Set(&mycache.CacheItem{Key: "key2", Value: []byte("value2")})
	}()
	go func() {
		cacheImpl.Set(&mycache.CacheItem{Key: "key3", Value: []byte("value3")})
	}()

	// Wait for insertions to complete
	time.Sleep(time.Millisecond * 100)

	// Concurrent access
	go func() {
		val, err := cacheImpl.Get("key1")
		if err != nil || string(val.Value) != "value1" {
			t.Errorf("Expected 'value1', got '%s'", val.Value)
		}
	}()
	go func() {
		val, err := cacheImpl.Get("key2")
		if err != nil || string(val.Value) != "value2" {
			t.Errorf("Expected 'value2', got '%s'", val.Value)
		}
	}()

	// Wait for access to complete
	time.Sleep(time.Millisecond * 100)

	// Concurrent deletion
	go func() {
		cacheImpl.Delete("key1")
	}()
	go func() {
		cacheImpl.Delete("key2")
	}()

	// Wait for deletions to complete
	time.Sleep(time.Millisecond * 100)

	// Concurrent clear
	go func() {
		cacheImpl.Clear()
	}()

	// Wait for clear to complete
	time.Sleep(time.Millisecond * 100)
	if cacheImpl.Len() != 0 {
		t.Error("Cache should be empty after concurrent clear")
	}
}
