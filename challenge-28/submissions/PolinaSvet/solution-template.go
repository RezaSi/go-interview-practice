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

//
// LRU Cache Implementation
//

type NodeLRU struct {
	key   string
	value interface{}
	prev  *NodeLRU
	next  *NodeLRU
}

type LRUCache struct {
	// TODO: Add necessary fields for LRU implementation
	// Hint: Use a doubly-linked list + hash map
	capacity int
	cache    map[string]*NodeLRU
	head     *NodeLRU // Most recently used
	tail     *NodeLRU // Least recently used
	mu       sync.RWMutex
	hits     int
	misses   int
}

// NewLRUCache creates a new LRU cache with the specified capacity
func NewLRUCache(capacity int) *LRUCache {
	// TODO: Implement LRU cache constructor
	if capacity <= 0 {
		return nil
	}

	return &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*NodeLRU),
		head:     nil,
		tail:     nil,
		hits:     0,
		misses:   0,
	}
}

func (c *LRUCache) removeNode(node *NodeLRU) {
	if node.prev != nil {
		node.prev.next = node.next
	} else {
		c.head = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	} else {
		c.tail = node.prev
	}

	node.prev = nil
	node.next = nil
}

func (c *LRUCache) addToFront(node *NodeLRU) {
	node.next = c.head
	node.prev = nil

	if c.head != nil {
		c.head.prev = node
	}
	c.head = node

	if c.tail == nil {
		c.tail = node
	}
}

func (c *LRUCache) moveToFront(node *NodeLRU) {
	if node == c.head {
		return
	}

	c.removeNode(node)
	c.addToFront(node)
}

func (c *LRUCache) Get(key string) (interface{}, bool) {
	// TODO: Implement LRU get operation
	// Should move accessed item to front (most recently used position)
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, exists := c.cache[key]; exists {
		c.moveToFront(node)
		c.hits++
		return node.value, true
	}

	c.misses++
	return nil, false
}

func (c *LRUCache) Put(key string, value interface{}) {
	// TODO: Implement LRU put operation
	// Should add new item to front and evict least recently used if at capacity
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, exists := c.cache[key]; exists {
		node.value = value
		c.moveToFront(node)
		return
	}

	newNode := &NodeLRU{
		key:   key,
		value: value,
	}

	if len(c.cache) >= c.capacity {
		if c.tail != nil {
			delete(c.cache, c.tail.key)
			c.removeNode(c.tail)
		}
	}

	c.addToFront(newNode)
	c.cache[key] = newNode

}

func (c *LRUCache) Delete(key string) bool {
	// TODO: Implement delete operation
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, exists := c.cache[key]; exists {
		c.removeNode(node)
		delete(c.cache, key)
		return true
	}
	return false
}

func (c *LRUCache) Clear() {
	// TODO: Implement clear operation
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache = make(map[string]*NodeLRU)
	c.head = nil
	c.tail = nil
	c.hits = 0
	c.misses = 0
}

func (c *LRUCache) Size() int {
	// TODO: Return current cache size
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.cache)
}

func (c *LRUCache) Capacity() int {
	// TODO: Return cache capacity
	return c.capacity
}

func (c *LRUCache) HitRate() float64 {
	// TODO: Calculate and return hit rate
	c.mu.RLock()
	defer c.mu.RUnlock()

	total := c.hits + c.misses
	if total == 0 {
		return 0.0
	}
	return float64(c.hits) / float64(total)
}

//
// LFU Cache Implementation
//

type NodeLFU struct {
	key       string
	value     interface{}
	freq      int
	prev      *NodeLFU
	next      *NodeLFU
	freqGroup *FreqGroup
}

type FreqGroup struct {
	freq  int
	head  *NodeLFU
	tail  *NodeLFU
	count int
}

type LFUCache struct {
	// TODO: Add necessary fields for LFU implementation
	// Hint: Use frequency tracking with efficient eviction
	capacity   int
	minFreq    int
	cache      map[string]*NodeLFU
	freqGroups map[int]*FreqGroup
	mu         sync.RWMutex
	hits       int
	misses     int
}

// NewLFUCache creates a new LFU cache with the specified capacity
func NewLFUCache(capacity int) *LFUCache {
	// TODO: Implement LFU cache constructor
	if capacity <= 0 {
		return nil
	}
	return &LFUCache{
		capacity:   capacity,
		minFreq:    1,
		cache:      make(map[string]*NodeLFU),
		freqGroups: make(map[int]*FreqGroup),
		hits:       0,
		misses:     0,
	}
}

func (c *LFUCache) createFreqGroup(freq int) *FreqGroup {
	group := &FreqGroup{
		freq:  freq,
		head:  nil,
		tail:  nil,
		count: 0,
	}
	c.freqGroups[freq] = group
	return group
}

func (c *LFUCache) removeNodeFromGroup(node *NodeLFU) {
	if node.freqGroup == nil {
		return
	}

	group := node.freqGroup

	if node.prev != nil {
		node.prev.next = node.next
	} else {
		group.head = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	} else {
		group.tail = node.prev
	}

	group.count--
	node.prev = nil
	node.next = nil

	if group.count == 0 {
		delete(c.freqGroups, group.freq)
		if group.freq == c.minFreq {
			c.minFreq++
		}
	}
}

func (c *LFUCache) addNodeToGroup(node *NodeLFU, group *FreqGroup) {
	node.freqGroup = group
	node.freq = group.freq

	node.next = group.head
	node.prev = nil

	if group.head != nil {
		group.head.prev = node
	}
	group.head = node

	if group.tail == nil {
		group.tail = node
	}

	group.count++
}

func (c *LFUCache) moveNodeToNewFreq(node *NodeLFU) {
	oldFreq := node.freq
	newFreq := oldFreq + 1

	c.removeNodeFromGroup(node)

	newGroup, exists := c.freqGroups[newFreq]
	if !exists {
		newGroup = c.createFreqGroup(newFreq)
	}

	c.addNodeToGroup(node, newGroup)

	if oldFreq == c.minFreq {
		if _, exists := c.freqGroups[oldFreq]; !exists {
			c.minFreq = newFreq
		}
	}
}

func (c *LFUCache) evictLFU() {
	minGroup, exists := c.freqGroups[c.minFreq]
	if !exists || minGroup.tail == nil {
		return
	}

	nodeToEvict := minGroup.tail
	delete(c.cache, nodeToEvict.key)
	c.removeNodeFromGroup(nodeToEvict)
}

func (c *LFUCache) Get(key string) (interface{}, bool) {
	// TODO: Implement LFU get operation
	// Should increment frequency count of accessed item
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, exists := c.cache[key]; exists {
		c.moveNodeToNewFreq(node)
		c.hits++
		return node.value, true
	}

	c.misses++
	return nil, false
}

func (c *LFUCache) Put(key string, value interface{}) {
	// TODO: Implement LFU put operation
	// Should evict least frequently used item if at capacity
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, exists := c.cache[key]; exists {
		node.value = value
		c.moveNodeToNewFreq(node)
		return
	}

	if len(c.cache) >= c.capacity {
		c.evictLFU()
	}

	newNode := &NodeLFU{
		key:   key,
		value: value,
		freq:  1,
	}

	group, exists := c.freqGroups[1]
	if !exists {
		group = c.createFreqGroup(1)
	}
	c.addNodeToGroup(newNode, group)

	c.cache[key] = newNode
	c.minFreq = 1
}

func (c *LFUCache) Delete(key string) bool {
	// TODO: Implement delete operation
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, exists := c.cache[key]; exists {
		c.removeNodeFromGroup(node)
		delete(c.cache, key)
		return true
	}
	return false
}

func (c *LFUCache) Clear() {
	// TODO: Implement clear operation
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache = make(map[string]*NodeLFU)
	c.freqGroups = make(map[int]*FreqGroup)
	c.minFreq = 1
	c.hits = 0
	c.misses = 0
}

func (c *LFUCache) Size() int {
	// TODO: Return current cache size
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.cache)
}

func (c *LFUCache) Capacity() int {
	// TODO: Return cache capacity
	return c.capacity
}

func (c *LFUCache) HitRate() float64 {
	// TODO: Calculate and return hit rate
	c.mu.RLock()
	defer c.mu.RUnlock()

	total := c.hits + c.misses
	if total == 0 {
		return 0.0
	}
	return float64(c.hits) / float64(total)
}

// FIFO Cache Implementation
type NodeFIFO struct {
	key   string
	value interface{}
	prev  *NodeFIFO
	next  *NodeFIFO
}

type FIFOCache struct {
	// TODO: Add necessary fields for FIFO implementation
	// Hint: Use a queue or circular buffer
	capacity int
	cache    map[string]*NodeFIFO
	head     *NodeFIFO // Newest item
	tail     *NodeFIFO // Oldest item
	mu       sync.RWMutex
	hits     int
	misses   int
}

// NewFIFOCache creates a new FIFO cache with the specified capacity
func NewFIFOCache(capacity int) *FIFOCache {
	// TODO: Implement FIFO cache constructor
	if capacity <= 0 {
		return nil
	}

	return &FIFOCache{
		capacity: capacity,
		cache:    make(map[string]*NodeFIFO),
		head:     nil,
		tail:     nil,
		hits:     0,
		misses:   0,
	}
}

func (c *FIFOCache) removeNode(node *NodeFIFO) {
	if node.prev != nil {
		node.prev.next = node.next
	} else {
		c.head = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	} else {
		c.tail = node.prev
	}

	node.prev = nil
	node.next = nil
}

func (c *FIFOCache) addToFront(node *NodeFIFO) {
	node.next = c.head
	node.prev = nil

	if c.head != nil {
		c.head.prev = node
	}
	c.head = node

	if c.tail == nil {
		c.tail = node
	}
}

func (c *FIFOCache) moveToFront(node *NodeFIFO) {
	if node == c.head {
		return
	}

	c.removeNode(node)
	c.addToFront(node)
}

func (c *FIFOCache) Get(key string) (interface{}, bool) {
	// TODO: Implement FIFO get operation
	// Note: Get operations don't affect eviction order in FIFO
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, exists := c.cache[key]; exists {
		c.hits++
		return node.value, true
	}

	c.misses++
	return nil, false
}

func (c *FIFOCache) Put(key string, value interface{}) {
	// TODO: Implement FIFO put operation
	// Should evict first-in item if at capacity
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, exists := c.cache[key]; exists {
		node.value = value
		return
	}

	newNode := &NodeFIFO{
		key:   key,
		value: value,
	}

	if len(c.cache) >= c.capacity {
		if c.tail != nil {
			delete(c.cache, c.tail.key)
			c.removeNode(c.tail)
		}
	}

	c.addToFront(newNode)
	c.cache[key] = newNode
}

func (c *FIFOCache) Delete(key string) bool {
	// TODO: Implement delete operation
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, exists := c.cache[key]; exists {
		c.removeNode(node)
		delete(c.cache, key)
		return true
	}
	return false
}

func (c *FIFOCache) Clear() {
	// TODO: Implement clear operation
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache = make(map[string]*NodeFIFO)
	c.head = nil
	c.tail = nil
	c.hits = 0
	c.misses = 0
}

func (c *FIFOCache) Size() int {
	// TODO: Return current cache size
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.cache)
}

func (c *FIFOCache) Capacity() int {
	// TODO: Return cache capacity
	return c.capacity
}

func (c *FIFOCache) HitRate() float64 {
	// TODO: Calculate and return hit rate
	c.mu.RLock()
	defer c.mu.RUnlock()

	total := c.hits + c.misses
	if total == 0 {
		return 0.0
	}
	return float64(c.hits) / float64(total)
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
	return &ThreadSafeCache{cache: cache}
}

func (c *ThreadSafeCache) Get(key string) (interface{}, bool) {
	// TODO: Implement thread-safe get operation
	// Hint: Use read lock for better performance
	c.mu.RLock()
	defer c.mu.RUnlock()

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
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cache.Size()
}

func (c *ThreadSafeCache) Capacity() int {
	// TODO: Implement thread-safe capacity operation
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cache.Capacity()
}

func (c *ThreadSafeCache) HitRate() float64 {
	// TODO: Implement thread-safe hit rate operation
	c.mu.RLock()
	defer c.mu.RUnlock()
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
		return NewLRUCache(capacity)
	case LFU:
		// TODO: Return LFU cache
		return NewLFUCache(capacity)
	case FIFO:
		// TODO: Return FIFO cache
		return NewFIFOCache(capacity)
	default:
		// TODO: Return default cache or handle error
		return NewFIFOCache(capacity)
	}

}

// NewThreadSafeCacheWithPolicy creates a thread-safe cache with the specified policy
func NewThreadSafeCacheWithPolicy(policy CachePolicy, capacity int) Cache {
	// TODO: Implement thread-safe cache factory
	// Should create cache with policy and wrap it with thread safety
	cache := NewCache(policy, capacity)
	if cache == nil {
		return nil
	}
	return NewThreadSafeCache(cache)
}
