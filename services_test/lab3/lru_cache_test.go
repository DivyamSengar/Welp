package services_test

import (
	"testing"
	"time"

	cache "cse190-welp/applications"
	"cse190-welp/proto/mycache"
)

func TestLRUCacheEviction(t *testing.T) {
	cache := cache.NewLRUCacheApp(3) // Create a cache with size 3

	cache.Set(&mycache.CacheItem{Key: "key1", Value: []byte("value1")})
	cache.Set(&mycache.CacheItem{Key: "key2", Value: []byte("value2")})
	cache.Set(&mycache.CacheItem{Key: "key3", Value: []byte("value3")})

	// Access several keys to update their positions in the cache
	_, _ = cache.Get("key1")
	_, _ = cache.Get("key2")
	_, _ = cache.Get("key1")
	_, _ = cache.Get("key2")
	_, _ = cache.Get("key3")

	// Insert a new item that should cause eviction
	cache.Set(&mycache.CacheItem{Key: "key4", Value: []byte("value4")})

	// Check that key1 got evicted due to its least recently used position
	_, err := cache.Get("key1")
	if err == nil {
		t.Errorf("Expected cache miss for 'key1'")
	}

	// Check that key3 did not get evicted
	_, err = cache.Get("key3")
	if err != nil {
		t.Errorf("Expected cache hit for 'key3'")
	}

	// Insert a new item that should cause eviction
	cache.Set(&mycache.CacheItem{Key: "key5", Value: []byte("value5")})

	// Check that key4 did not get evicted
	_, err = cache.Get("key4")
	if err != nil {
		t.Errorf("Expected cache hit for 'key4'")
	}

	// Check that key3 did not get evicted
	_, err = cache.Get("key3")
	if err != nil {
		t.Errorf("Expected cache hit for 'key3'")
	}

	// Check that key2 did get evicted
	_, err = cache.Get("key2")
	if err == nil {
		t.Errorf("Expected cache miss for 'key2'")
	}
}

func TestLRUCacheConcurrencyEviction(t *testing.T) {
	cache := cache.NewLRUCacheApp(3) // Create a cache with size 3

	// Concurrent insertion
	go func() {
		cache.Set(&mycache.CacheItem{Key: "key1", Value: []byte("value1")})
	}()
	go func() {
		cache.Set(&mycache.CacheItem{Key: "key2", Value: []byte("value2")})
	}()
	go func() {
		cache.Set(&mycache.CacheItem{Key: "key3", Value: []byte("value3")})
	}()

	// Wait for insertions to complete
	time.Sleep(time.Millisecond * 100)

	// Concurrent access that updates key1's position
	go func() {
		_, _ = cache.Get("key1")
	}()
	// Concurrent access that updates key2's position
	go func() {
		_, _ = cache.Get("key2")
	}()

	// Wait for access to complete
	time.Sleep(time.Millisecond * 100)

	// Insert a new item that should cause eviction
	go func() {
		cache.Set(&mycache.CacheItem{Key: "key4", Value: []byte("value4")})
	}()

	// Wait for insertion to complete
	time.Sleep(time.Millisecond * 100)

	// Check that key3 got evicted due to its least recently used position
	_, err := cache.Get("key3")
	if err == nil {
		t.Errorf("Expected cache miss for 'key3'")
	}
}
