package applications

import (
	"container/list"
	"cse190-welp/proto/mycache"
	"errors"
	"log"
	"math/rand"
	"sync"
)

var (
	ErrItemNotFound = errors.New("mycache: cache miss")
)

// Cache is a simple Key-Value cache interface.
type Cache interface {
	// Len returns the number of elements in the cache.
	Len() int

	// Get retrieves the value for the specified key.
	Get(key string) (*mycache.CacheItem, error)

	// Set sets the value for the specified key. If the maximum capacity of the cache is exceeded,
	// an eviction policy is applied.
	Set(item *mycache.CacheItem) error

	// Delete deletes the value for the specified key.
	Delete(key string) error

	// Clear removes all items from the cache.
	Clear()
}

// FIFOCacheApp is a simple in-memory FIFO (First-In-First-Out) key-value cache.
type FIFOCacheApp struct {
	data     map[string]*mycache.CacheItem
	order    *list.List // Use a doubly-linked list to maintain FIFO order
	capacity int
	lock     sync.Mutex
}

// NewFIFOCacheApp returns a new FIFO Cache with the specified maximum capacity.
func NewFIFOCacheApp(capacity int) *FIFOCacheApp {
	log.Println("eviction policy: FIFO cache")
	return &FIFOCacheApp{
		data:     make(map[string]*mycache.CacheItem),
		order:    list.New(),
		capacity: capacity,
	}
}

// Len returns the number of elements in the cache.
func (c *FIFOCacheApp) Len() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return len(c.data)
}

// Get retrieves the value for the specified key.
func (c *FIFOCacheApp) Get(key string) (*mycache.CacheItem, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	value, ok := c.data[key]
	if !ok {
		return nil, ErrItemNotFound
	}
	return value, nil
}

// Set sets the value for the specified key. If the maximum capacity of the cache is exceeded,
// the oldest key-value pair will be evicted.
func (c *FIFOCacheApp) Set(item *mycache.CacheItem) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	// Don't do anything if cache has size of 0
	if c.capacity == 0 {
		return nil
	}

	if len(c.data) >= c.capacity {
		// If the cache is full, evict the oldest item (front of the list)
		oldestElement := c.order.Front()
		if oldestElement != nil {
			oldestKey := oldestElement.Value.(string)
			delete(c.data, oldestKey)
			c.order.Remove(oldestElement)
		}
	}

	key := item.Key
	c.data[key] = item
	c.order.PushBack(key) // Add the new key to the back of the list
	return nil
}

// Delete deletes the value for the specified key.
func (c *FIFOCacheApp) Delete(key string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	_, ok := c.data[key]
	if !ok {
		return ErrItemNotFound
	}

	// Remove the key from the data map and the list
	delete(c.data, key)
	for element := c.order.Front(); element != nil; element = element.Next() {
		if element.Value.(string) == key {
			c.order.Remove(element)
			break
		}
	}
	return nil
}

// Clear removes all items from the cache.
func (c *FIFOCacheApp) Clear() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.data = make(map[string]*mycache.CacheItem)
	c.order.Init()
}

// RandomCacheApp is a simple in-memory key-value cache.
type RandomCacheApp struct {
	lock     sync.Mutex // Mutex for protecting concurrent access
	data     map[string]*mycache.CacheItem
	capacity int
}

// NewRandomCacheApp returns a new Cache with the specified maximum capacity.
func NewRandomCacheApp(capacity int) *RandomCacheApp {
	log.Println("eviction policy: random cache")
	return &RandomCacheApp{
		lock:     sync.Mutex{},
		data:     make(map[string]*mycache.CacheItem),
		capacity: capacity,
	}
}

// Len returns the number of elements in the cache.
func (c *RandomCacheApp) Len() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return len(c.data)
}

// Get retrieves the value for the specified key.
func (c *RandomCacheApp) Get(key string) (*mycache.CacheItem, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	value, ok := c.data[key]
	if !ok {
		return nil, ErrItemNotFound
	}
	return value, nil
}

// Set sets the value for the specified key. If the maximum capacity of the cache is exceeded,
// the oldest key-value pair will be evicted.
func (c *RandomCacheApp) Set(item *mycache.CacheItem) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	// Don't do anything if cache has size of 0
	if c.capacity == 0 {
		return nil
	}

	if len(c.data) >= c.capacity {
		_ = c.evictRandomKey()
	}
	key := item.Key
	c.data[key] = item
	return nil
}

// Remove deletes the value for the specified key.
func (c *RandomCacheApp) Delete(key string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, ok := c.data[key]
	if !ok {
		return ErrItemNotFound
	}
	delete(c.data, key)
	return nil
}

// Clear removes all items from the cache.
func (c *RandomCacheApp) Clear() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.data = make(map[string]*mycache.CacheItem)
}

// evictRandomKey deletes a random key-value pair from the cache and returns the evicted key.
func (c *RandomCacheApp) evictRandomKey() string {
	randomKey := ""
	randomIndex := rand.Intn(len(c.data))

	i := 0
	for key := range c.data {
		if i == randomIndex {
			randomKey = key
			break
		}
		i++
	}
	delete(c.data, randomKey)
	return randomKey
}

// Find the element in a linked list given its key
func GetListElement(order *list.List, key string) *list.Element {
	// Traverse the linked list
	for element := order.Front(); element != nil; element = element.Next() {
		// If the key is found, return the element
		if element.Value == key {
			return element
		}
	}
	// If key not found, return nil
	return nil
}

// LRUCacheApp is an in-memory Least-Recently-Used key-value cache.
// To enforce LRU, a linkedlist is kept that relocates a node to the back each
// time it is accessed. Then the node to evict will always be at the front.
type LRUCacheApp struct {
	data     map[string]*mycache.CacheItem
	order    *list.List // Use a doubly-linked list to maintain LRU order
	capacity int
	lock     sync.Mutex
}

// NewLRUCacheApp returns a new LRU Cache with the specified maximum capacity.
func NewLRUCacheApp(capacity int) *LRUCacheApp {
	log.Println("eviction policy: LRU cache")
	return &LRUCacheApp{
		data:     make(map[string]*mycache.CacheItem),
		order:    list.New(),
		capacity: capacity,
	}
}

// Len returns the number of elements in the cache.
func (c *LRUCacheApp) Len() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return len(c.data)
}

// Get retrieves the value for the specified key.
func (c *LRUCacheApp) Get(key string) (*mycache.CacheItem, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	value, ok := c.data[key]
	if !ok {
		return nil, ErrItemNotFound
	}
	// Get the element we just accessed
	accessedElement := GetListElement(c.order, key)
	if accessedElement != nil {
		// Send the accessed key to the back of the linked list.
		c.order.MoveToBack(accessedElement)
	}
	return value, nil
}

// Set sets the value for the specified key. If the maximum capacity of the cache is exceeded,
// the oldest key-value pair will be evicted.
func (c *LRUCacheApp) Set(item *mycache.CacheItem) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	// Don't do anything if cache has size of 0
	if c.capacity == 0 {
		return ErrItemNotFound
	}
	// fmt.Println("Hello, World!")
	if len(c.data) >= c.capacity {
		// If the cache is full, evict the oldest item (front of the list)
		oldestElement := c.order.Front()
		if oldestElement != nil {
			oldestKey := oldestElement.Value.(string)
			delete(c.data, oldestKey)
			c.order.Remove(oldestElement)
		}
	}

	key := item.Key
	c.data[key] = item
	c.order.PushBack(key) // Add the new key to the back of the list
	return nil
}

// Delete deletes the value for the specified key.
func (c *LRUCacheApp) Delete(key string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	_, ok := c.data[key]
	if !ok {
		return ErrItemNotFound
	}

	// Remove the key from the data map and the list
	delete(c.data, key)
	for element := c.order.Front(); element != nil; element = element.Next() {
		if element.Value.(string) == key {
			c.order.Remove(element)
			break
		}
	}
	return nil
}

// Clear removes all items from the cache.
func (c *LRUCacheApp) Clear() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.data = make(map[string]*mycache.CacheItem)
	c.order.Init()
}

type LFUPair struct {
	key   string
	count int
}

// Find the element in a linked list of LFUPair given its key
func GetLFUListElement(order *list.List, key string) *list.Element {
	// Traverse the linked list
	for element := order.Front(); element != nil; element = element.Next() {
		// If the key is found, return the element
		if element.Value.(LFUPair).key == key {
			return element
		}
	}
	// If key not found, return nil
	return nil
}

// Find the element with the smallest amount of accesses
func GetLeastUsedElement(order *list.List) *list.Element {
	var minAccesses int = order.Front().Value.(LFUPair).count
	minElement := order.Front()
	// Traverse the linked list
	for element := order.Front(); element != nil; element = element.Next() {
		// If an element is found with less accesses, that is the new element
		if element.Value.(LFUPair).count < minAccesses {
			minAccesses = element.Value.(LFUPair).count
			minElement = element
		}
	}
	// Return the element with minimum accesses
	return minElement
}

// LFUCacheApp is an in-memory Least-Frequently-Used key-value cache.
// To enforce LFU, a linkedlist is kept that relocates a node to the back each
// time it is accessed. Then the node to evict will always be at the front.
type LFUCacheApp struct {
	data     map[string]*mycache.CacheItem
	order    *list.List // Doubly-linked list stores (key, count) pairs to maintain LFU order.
	capacity int
	lock     sync.Mutex
}

// NewLFUCacheApp returns a new LRU Cache with the specified maximum capacity.
func NewLFUCacheApp(capacity int) *LFUCacheApp {
	log.Println("eviction policy: LFU cache")
	return &LFUCacheApp{
		data:     make(map[string]*mycache.CacheItem),
		order:    list.New(),
		capacity: capacity,
	}
}

// Len returns the number of elements in the cache.
func (c *LFUCacheApp) Len() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return len(c.data)
}

// Get retrieves the value for the specified key.
func (c *LFUCacheApp) Get(key string) (*mycache.CacheItem, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	value, ok := c.data[key]
	if !ok {
		return nil, ErrItemNotFound
	}
	// Get the element we just accessed
	var accessedElement *list.Element = GetLFUListElement(c.order, key)
	// Increment the pair's access count
	newCount := accessedElement.Value.(LFUPair).count + 1
	if accessedElement != nil {
		// Update the list with the new pair
		c.order.Remove(accessedElement)
		c.order.PushBack(LFUPair{key: key, count: newCount})
	}
	return value, nil
}

// Set sets the value for the specified key. If the maximum capacity of the cache is exceeded,
// the least-accessed key-value pair will be evicted.
func (c *LFUCacheApp) Set(item *mycache.CacheItem) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	// Don't do anything if cache has size of 0
	if c.capacity == 0 {
		return nil
	}

	if len(c.data) >= c.capacity {
		// If the cache is full, evict the least-used item
		oldestElement := GetLeastUsedElement(c.order)
		oldestKey := oldestElement.Value.(LFUPair).key
		delete(c.data, oldestKey)
		c.order.Remove(oldestElement)
	}

	key := item.Key
	c.data[key] = item
	c.order.PushBack(LFUPair{key: key, count: 0}) // Add the new key to the back of the list
	return nil
}

// Delete deletes the value for the specified key.
func (c *LFUCacheApp) Delete(key string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	_, ok := c.data[key]
	if !ok {
		return ErrItemNotFound
	}

	// Remove the key from the data map and the list
	delete(c.data, key)
	for element := c.order.Front(); element != nil; element = element.Next() {
		if element.Value.(LFUPair).key == key {
			c.order.Remove(element)
			break
		}
	}
	return nil
}

// Clear removes all items from the cache.
func (c *LFUCacheApp) Clear() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.data = make(map[string]*mycache.CacheItem)
	c.order.Init()
}
