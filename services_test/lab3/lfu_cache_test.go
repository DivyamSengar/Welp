package services_test

import (
	"testing"
	"time"

	cache "cse190-welp/applications"
	"cse190-welp/proto/mycache"
)

func TestLFUCacheEviction(t *testing.T) {
	cache := cache.NewLFUCacheApp(3) // Create a cache with size 3

	cache.Set(&mycache.CacheItem{Key: "key1", Value: []byte("value1")})
	cache.Set(&mycache.CacheItem{Key: "key2", Value: []byte("value2")})
	cache.Set(&mycache.CacheItem{Key: "key3", Value: []byte("value3")})

	// Access several keys to change their frequencies
	_, _ = cache.Get("key1")
	_, _ = cache.Get("key2")
	_, _ = cache.Get("key1")
	_, _ = cache.Get("key2")
	_, _ = cache.Get("key3")

	// Insert a new item that should cause eviction
	cache.Set(&mycache.CacheItem{Key: "key4", Value: []byte("value4")})

	// Check that key3 got evicted due to its low frequency
	_, err := cache.Get("key3")
	if err == nil {
		t.Errorf("Expected cache miss for 'key3'")
	}

	// Check that key1 did not get evicted
	_, err = cache.Get("key1")
	if err != nil {
		t.Errorf("Expected cache hit for 'key1'")
	}

	// Check that key2 did not get evicted
	_, err = cache.Get("key2")
	if err != nil {
		t.Errorf("Expected cache hit for 'key2'")
	}

	// Access keys to increase their frequency
	cache.Set(&mycache.CacheItem{Key: "key3", Value: []byte("value3")})
	for i := 0; i < 6; i++ {
		_, _ = cache.Get("key3")
	}
	_, _ = cache.Get("key1")

	// Insert a new item that should cause eviction
	cache.Set(&mycache.CacheItem{Key: "key5", Value: []byte("value5")})

	// Check that key1 did not get evicted
	_, err = cache.Get("key1")
	if err != nil {
		t.Errorf("Expected cache hit for 'key1'")
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

func TestLFUCacheConcurrencyEviction(t *testing.T) {
	cache := cache.NewLFUCacheApp(3) // Create a cache with size 3

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

	// Concurrent access that increases frequency of key1
	go func() {
		_, _ = cache.Get("key1")
	}()
	// Concurrent access that increases frequency of key2
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

	// Check that key3 got evicted due to its low frequency
	_, err := cache.Get("key3")
	if err == nil {
		t.Errorf("Expected cache miss for 'key3'")
	}
}
