package cache

import (
	"sync"
)

// Cache interface defines the contract for all cache implementations
type Cache interface {
	Get(key string) (value interface{}, found bool)
	Put(key string, value interface{})
	Delete(key string) bool
	Clear()
	Size() int
	Capacity() int
	HitRate() float64
}

// CachePolicy represents the eviction policy type
type CachePolicy int

const (
	LRU CachePolicy = iota
	LFU
	FIFO
)

// LRU Cache Implementation
type Metrics struct {
	hits      int
	misses    int
	evictions int
}

type LRUCache struct {
	// TODO: Add necessary fields for LRU implementation
	// Hint: Use a doubly-linked list + hash map
	capacity int
	size     int
	cache    map[string]*LRUNode
	metrics  Metrics
	head     *LRUNode
	tail     *LRUNode
}

// NewLRUCache creates a new LRU cache with the specified capacity
func NewLRUCache(capacity int) *LRUCache {
	if capacity < 1 {
		return nil
	}
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*LRUNode, capacity),
		metrics:  Metrics{},
	}
}

type LRUNode struct {
	key   string
	value interface{}
	next  *LRUNode
	prev  *LRUNode
}

func newLRUNode(key string, value interface{}) *LRUNode {
	if key == "" {
		return nil
	}
	return &LRUNode{
		key:   key,
		value: value,
	}
}

func (c *LRUCache) evictLRUNode() bool {
	if c.size == 0 {
		return false
	}
	node := c.tail
	return c.removeLRUNode(node)
}
func (c *LRUCache) removeLRUNode(node *LRUNode) bool {
	if c.size == 0 || node == nil {
		return false
	}

	if c.size == 1 {
		c.head = nil
		c.tail = nil
		delete(c.cache, node.key)
		c.size--
		return true
	}
	if c.tail == node {
		c.tail = node.prev
		c.tail.next = nil
		node.prev = nil
		delete(c.cache, node.key)
		c.size--
		return true
	}
	if c.head == node {
		c.head = node.next
		c.head.prev = nil
		node.next = nil
		delete(c.cache, node.key)
		c.size--
		return true
	}
	node.prev.next = node.next
	node.next.prev = node.prev
	node.next = nil
	node.prev = nil
	delete(c.cache, node.key)
	c.size--
	return true
}

func (c *LRUCache) moveToFront(node *LRUNode) bool {
	if node == nil {
		return false
	}
	if c.head == node {
		return true
	}
	if _, exists := c.cache[node.key]; exists {
		c.removeLRUNode(node)
	}
	if c.head == nil {
		c.head = node
		c.tail = node
	} else {
		node.next = c.head
		c.head.prev = node
		c.head = node
		node.prev = nil
	}
	c.cache[node.key] = node
	c.size++
	return true
}

func (c *LRUCache) Get(key string) (interface{}, bool) {
	if key == "" {
		return nil, false
	}
	if node, exists := c.cache[key]; exists {
		c.moveToFront(node)
		c.metrics.hits++
		return node.value, true
	}
	c.metrics.misses++
	return nil, false
}

func (c *LRUCache) Put(key string, value interface{}) {
	if key == "" {
		return
	}
	if node, exists := c.cache[key]; exists {
		node.value = value
		c.moveToFront(node)
		return
	}
	node := newLRUNode(key, value)
	if c.size == c.capacity {
		c.evictLRUNode()
	}
	c.moveToFront(node)
}

func (c *LRUCache) Delete(key string) bool {
	if key == "" {
		return false
	}
	if node, exists := c.cache[key]; exists {
		c.removeLRUNode(node)
		return true
	}
	return false
}

func (c *LRUCache) Clear() {
	c.size = 0
	c.cache = make(map[string]*LRUNode)
	c.metrics = Metrics{}
}

func (c *LRUCache) Size() int {
	return c.size
}

func (c *LRUCache) Capacity() int {
	return c.capacity
}

func (c *LRUCache) HitRate() float64 {
	requests := c.metrics.misses + c.metrics.hits
	if requests == 0 {
		return 0.0
	}
	return float64(c.metrics.hits) / float64(requests)
}

//
// LFU Cache Implementation
//

type LFUCache struct {
	// TODO: Add necessary fields for LFU implementation
	// Hint: Use frequency tracking with efficient eviction
	capacity   int
	size       int
	cache      map[string]*LFUNode
	freqGroups map[int]*FreqGroup
	freqList   []int
	metrics    Metrics
}

type LFUNode struct {
	key   string
	value interface{}
	freq  int
	next  *LFUNode
	prev  *LFUNode
}

type FreqGroup struct {
	freq int
	size int
	head *LFUNode
	tail *LFUNode
}

// NewLFUCache creates a new LFU cache with the specified capacity
func NewLFUCache(capacity int) *LFUCache {
	// TODO: Implement LFU cache constructor
	if capacity <= 0 {
		return nil
	}
	return &LFUCache{
		capacity: capacity,
		cache:    make(map[string]*LFUNode, capacity),
		freqList: make([]int, 0),
	}
}

func newLFUNode(key string, value interface{}) *LFUNode {
	if key == "" {
		return nil
	}
	return &LFUNode{
		key:   key,
		value: value,
	}
}

func newFreqGroup(freq int) *FreqGroup {
	if freq <= 0 {
		return nil
	}
	return &FreqGroup{
		freq: freq,
	}
}

func (c *LFUCache) evictLFUNode() bool {
	if bucket, exists := c.freqGroups[c.freqList[0]]; exists {
		c.removeLFUNode(bucket.tail)
		return true
	}
	return false
}

func (c *LFUCache) removeLFUNode(node *LFUNode) bool {
	if node == nil {
		return false
	}
	if bucket, exists := c.freqGroups[node.freq]; exists {
		if bucket.size == 0 {
			delete(c.freqGroups, bucket.freq)
			return false
		}
		if bucket.size == 1 {
			bucket.head = nil
			bucket.tail = nil
			delete(c.freqGroups, bucket.freq)
			c.size--
			return true
		}
		if bucket.head == node {
			bucket.head = node.next
			bucket.head.prev = nil
			node.next = nil
			bucket.size--
			c.size--
			return true
		}
		if bucket.tail == node {
			bucket.tail = node.prev
			bucket.tail.next = nil
			node.prev = nil
			bucket.size--
			c.size--
			return true
		}
	}
	return false
}

func (c *LFUCache) Get(key string) (interface{}, bool) {
	// TODO: Implement LFU get operation
	// Should increment frequency count of accessed item
	return nil, false
}

func (c *LFUCache) Put(key string, value interface{}) {
	// TODO: Implement LFU put operation
	// Should evict least frequently used item if at capacity
}

func (c *LFUCache) Delete(key string) bool {
	// TODO: Implement delete operation
	return false
}

func (c *LFUCache) Clear() {
	// TODO: Implement clear operation
}

func (c *LFUCache) Size() int {
	// TODO: Return current cache size
	return 0
}

func (c *LFUCache) Capacity() int {
	// TODO: Return cache capacity
	return 0
}

func (c *LFUCache) HitRate() float64 {
	// TODO: Calculate and return hit rate
	return 0.0
}

//
// FIFO Cache Implementation
//

type FIFOCache struct {
	// TODO: Add necessary fields for FIFO implementation
	// Hint: Use a queue or circular buffer
	capacity int
	size     int
	cache    map[string]*FIFONode
	metrics  Metrics
	head     *FIFONode
	tail     *FIFONode
}

type FIFONode struct {
	key   string
	value interface{}
	next  *FIFONode
	prev  *FIFONode
}

// NewFIFOCache creates a new FIFO cache with the specified capacity
func NewFIFOCache(capacity int) *FIFOCache {
	// TODO: Implement FIFO cache constructor
	if capacity <= 0 {
		return nil
	}
	return &FIFOCache{
		capacity: capacity,
		cache:    make(map[string]*FIFONode),
		metrics:  Metrics{},
	}
}

func newFIFONode(key string, value interface{}) *FIFONode {
	if key == "" {
		return nil
	}
	return &FIFONode{
		key:   key,
		value: value,
	}
}

func (c *FIFOCache) evictFIFONode() bool {
	if c.size == 0 {
		return false
	}
	node := c.tail
	return c.removeFIFONode(node)
}

func (c *FIFOCache) removeFIFONode(node *FIFONode) bool {
	if node == nil || c.size == 0 {
		return false
	}
	if c.size == 1 {
		c.head = nil
		c.tail = nil
		delete(c.cache, node.key)
		c.size = 0
		return true
	}
	if c.head == node {
		c.head = node.next
		c.head.prev = nil
		node.next = nil
		delete(c.cache, node.key)
		c.size--
		return true
	}
	if c.tail == node {
		c.tail = node.prev
		c.tail.next = nil
		node.prev = nil
		delete(c.cache, node.key)
		c.size--
		return true
	}
	node.prev.next = node.next
	node.next.prev = node.prev
	node.next = nil
	node.prev = nil
	delete(c.cache, node.key)
	c.size--
	return true
}

func (c *FIFOCache) Get(key string) (interface{}, bool) {
	// TODO: Implement FIFO get operation
	// Note: Get operations don't affect eviction order in FIFO
	if key == "" {
		return nil, false
	}
	if node, exists := c.cache[key]; exists {
		c.metrics.hits++
		return node.value, true
	}
	c.metrics.misses++
	return nil, false
}

func (c *FIFOCache) Put(key string, value interface{}) {
	// TODO: Implement FIFO put operation
	// Should evict first-in item if at capacity
	if key == "" {
		return
	}
	if node, exists := c.cache[key]; exists {
		node.value = value
		return
	}
	node := newFIFONode(key, value)

	if c.size == c.capacity {
		c.evictFIFONode()
		c.metrics.evictions++
	}
	if c.size == 0 {
		c.head = node
		c.tail = node
		c.cache[key] = node
		c.size++
		return
	}
	node.next = c.head
	c.head = node
	node.next.prev = node
	c.cache[key] = node
	c.size++
}

func (c *FIFOCache) Delete(key string) bool {
	// TODO: Implement delete operation
	if key == "" {
		return false
	}
	if node, exists := c.cache[key]; exists {
		return c.removeFIFONode(node)
	}
	return false
}

func (c *FIFOCache) Clear() {
	// TODO: Implement clear operation
	c.cache = make(map[string]*FIFONode)
	c.head = nil
	c.tail = nil
	c.size = 0
	c.metrics = Metrics{}
}

func (c *FIFOCache) Size() int {
	// TODO: Return current cache size
	return c.size
}

func (c *FIFOCache) Capacity() int {
	// TODO: Return cache capacity
	return c.capacity
}

func (c *FIFOCache) HitRate() float64 {
	// TODO: Calculate and return hit rate
	total := c.metrics.hits + c.metrics.misses
	if total > 0 {
		return float64(c.metrics.hits) / float64(total)
	}
	return 0.0
}

//
// Thread-Safe Cache Wrapper
//

type ThreadSafeCache struct {
	cache Cache
	mu    sync.RWMutex
	// TODO: Add any additional fields if needed
}

// NewThreadSafeCache wraps any cache implementation to make it thread-safe
func NewThreadSafeCache(cache Cache) *ThreadSafeCache {
	// TODO: Implement thread-safe wrapper constructor
	if cache == nil {
		return nil
	}
	return &ThreadSafeCache{
		cache: cache,
	}
}

func (c *ThreadSafeCache) Get(key string) (interface{}, bool) {
	// TODO: Implement thread-safe get operation
	// Hint: Use read lock for better performance
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.cache.Get(key)
}

func (c *ThreadSafeCache) Put(key string, value interface{}) {
	// TODO: Implement thread-safe put operation
	// Hint: Use write lock
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache.Put(key, value)
}

func (c *ThreadSafeCache) Delete(key string) bool {
	// TODO: Implement thread-safe delete operation
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.cache.Delete(key)
}

func (c *ThreadSafeCache) Clear() {
	// TODO: Implement thread-safe clear operation
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache.Clear()
}

func (c *ThreadSafeCache) Size() int {
	// TODO: Implement thread-safe size operation
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.cache.Size()
}

func (c *ThreadSafeCache) Capacity() int {
	// TODO: Implement thread-safe capacity operation
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.cache.Capacity()
}

func (c *ThreadSafeCache) HitRate() float64 {
	// TODO: Implement thread-safe hit rate operation
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.cache.HitRate()
}

//
// Cache Factory Functions
//

// NewCache creates a cache with the specified policy and capacity
func NewCache(policy CachePolicy, capacity int) Cache {
	// TODO: Implement cache factory
	// Should create appropriate cache type based on policy
	switch policy {
	case LRU:
		// TODO: Return LRU cache
		return NewLRUCache(capacity)
	case LFU:
		// TODO: Return LFU cache
		return NewLFUCache(capacity)
	case FIFO:
		// TODO: Return FIFO cache
		return NewFIFOCache(capacity)
	default:
		// TODO: Return default cache or handle error
		return NewLRUCache(capacity)
	}
}

// NewThreadSafeCacheWithPolicy creates a thread-safe cache with the specified policy
func NewThreadSafeCacheWithPolicy(policy CachePolicy, capacity int) Cache {
	// TODO: Implement thread-safe cache factory
	// Should create cache with policy and wrap it with thread safety
	switch policy {
	case LRU:
		// TODO: Return LRU cache
		return NewThreadSafeCache(NewLRUCache(capacity))
	case LFU:
		// TODO: Return LFU cache
		return NewThreadSafeCache(NewLFUCache(capacity))
	case FIFO:
		// TODO: Return FIFO cache
		return NewThreadSafeCache(NewFIFOCache(capacity))
	default:
		// TODO: Return default cache or handle error
		return NewThreadSafeCache(NewLRUCache(capacity))
	}
}
