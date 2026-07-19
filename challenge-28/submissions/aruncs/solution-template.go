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

const head = "head"
const tail = "tail"

type Cachevalue struct {
	value interface{}
}

type ListNode struct {
	key   string
	value interface{}
	next  *ListNode
	prev  *ListNode
}

type LFUListNode struct {
	key       string
	value     interface{}
	next      *LFUListNode
	prev      *LFUListNode
	frequency int
}

type FrequncyListDetails struct {
	head *LFUListNode
	tail *LFUListNode
}

type LRUCache struct {
	// TODO: Add necessary fields for LRU implementation
	// Hint: Use a doubly-linked list + hash map
	mu           sync.RWMutex
	head         *ListNode
	tail         *ListNode
	hashMap      map[string]*ListNode
	cap          int
	requestCount int
	hitCount     int
}

// NewLRUCache creates a new LRU cache with the specified capacity
func NewLRUCache(capacity int) *LRUCache {
	if capacity < 0 {
		capacity = 0
	}
	// TODO: Implement LRU cache constructor
	cache := &LRUCache{
		head: &ListNode{
			key: head,
		},
		tail: &ListNode{
			key: tail,
		},
		hashMap: make(map[string]*ListNode, capacity),
		cap:     capacity,
	}

	cache.head.next = cache.tail
	cache.tail.prev = cache.head

	return cache
}
func (c *LRUCache) removeNode(node *ListNode) {

	node.prev.next = node.next
	node.next.prev = node.prev
}

func (c *LRUCache) addToFront(node *ListNode) {

	node.prev = c.head
	node.next = c.head.next

	node.next.prev = node
	c.head.next = node
}

func (c *LRUCache) Get(key string) (interface{}, bool) {
	// TODO: Implement LRU get operation
	// Should move accessed item to front (most recently used position)
	c.incrementrequestCount()

	val, ok := c.hashMap[key]
	if !ok {
		return nil, false
	}

	c.incrementhitCount()

	c.removeNode(val)
	c.addToFront(val)

	return val.value, true
}

func (c *LRUCache) Put(key string, value interface{}) {
	// TODO: Implement LRU put operation
	// Should add new item to front and evict least recently used if at capacity
	if c.cap == 0 {
		return
	}
	existing, ok := c.hashMap[key]
	if ok {
		existing.value = value
		c.removeNode(existing)
		c.addToFront(existing)
		return
	}

	if len(c.hashMap) == c.cap {
		toBeDeleted := c.tail.prev
		c.removeNode(toBeDeleted)
		delete(c.hashMap, toBeDeleted.key)
	}

	node := &ListNode{
		key:   key,
		value: value,
	}
	c.addToFront(node)

	c.hashMap[key] = node
}

func (c *LRUCache) Delete(key string) bool {
	// TODO: Implement delete operation
	val, ok := c.hashMap[key]
	if !ok {
		return false
	}

	delete(c.hashMap, key)

	val.prev.next = val.next
	val.next.prev = val.prev

	return true
}

func (c *LRUCache) Clear() {
	// TODO: Implement clear operation

	c.head.next = c.tail
	c.tail.prev = c.head

	clear(c.hashMap)
}

func (c *LRUCache) Size() int {
	// TODO: Return current cache size

	return len(c.hashMap)
}

func (c *LRUCache) Capacity() int {
	// TODO: Return cache capacity
	return c.cap
}

func (c *LRUCache) HitRate() float64 {
	// TODO: Calculate and return hit rate
	if c.requestCount == 0 {
		return 0
	}
	return float64(c.hitCount) / float64(c.requestCount)
}

func (c *LRUCache) incrementrequestCount() {
	c.requestCount = c.requestCount + 1
}

func (c *LRUCache) incrementhitCount() {
	c.hitCount = c.hitCount + 1
}

//
// LFU Cache Implementation
//

type LFUCache struct {
	// TODO: Add necessary fields for LFU implementation
	// Hint: Use frequency tracking with efficient eviction
	mu             sync.RWMutex
	hashMap        map[string]*LFUListNode
	Leastfrequency int
	FrequncyList   map[int]*FrequncyListDetails
	cap            int
	requestCount   int
	hitCount       int
}

// NewLFUCache creates a new LFU cache with the specified capacity
func NewLFUCache(capacity int) *LFUCache {
	// TODO: Implement LFU cache constructor
	if capacity < 0 {
		capacity = 0
	}

	return &LFUCache{
		hashMap:        make(map[string]*LFUListNode, capacity),
		FrequncyList:   make(map[int]*FrequncyListDetails, capacity),
		Leastfrequency: 0,
		cap:            capacity,
	}
}

func (c *LFUCache) Get(key string) (interface{}, bool) {
	// TODO: Implement LFU get operation
	// Should increment frequency count of accessed item
	c.incrementrequestCount()

	node, exists := c.hashMap[key]
	if !exists {
		return nil, false
	}

	c.incrementhitCount()

	c.removeNode(node)
	node.frequency++
	c.addToFront(node)
	return node.value, true
}

func (c *LFUCache) Put(key string, value interface{}) {
	// TODO: Implement LFU put operation
	// Should evict least frequently used item if at capacity
	if c.cap == 0 {
		return
	}

	node, exists := c.hashMap[key]

	// exists
	if exists {
		c.removeNode(node)
		node.value = value
		node.frequency += 1
		c.addToFront(node)
		return
	}

	//not exists
	if len(c.hashMap) == c.cap {
		list := c.FrequncyList[c.Leastfrequency]
		nodeToRemove := list.tail.prev
		delete(c.hashMap, nodeToRemove.key)
		c.removeNode(nodeToRemove)
	}

	_, exists = c.FrequncyList[1]
	if !exists {
		c.initialiseNewfrequency(1)
	}
	node = &LFUListNode{
		key:       key,
		value:     value,
		frequency: 1,
	}
	c.hashMap[key] = node
	c.Leastfrequency = 1
	c.addToFront(node)
}

func (c *LFUCache) Delete(key string) bool {
	// TODO: Implement delete operation

	node, exists := c.hashMap[key]

	if !exists {
		return false
	}

	delete(c.hashMap, key)

	c.removeNode(node)
	return false
}

func (c *LFUCache) Clear() {
	// TODO: Implement clear operation

	clear(c.hashMap)
	clear(c.FrequncyList)
}

func (c *LFUCache) Size() int {
	// TODO: Return current cache size
	return len(c.hashMap)
}

func (c *LFUCache) Capacity() int {
	// TODO: Return cache capacity
	return c.cap
}

func (c *LFUCache) HitRate() float64 {
	// TODO: Calculate and return hit rate
	if c.requestCount == 0 {
		return 0
	}
	return float64(c.hitCount) / float64(c.requestCount)
}

func (c *LFUCache) initialiseNewfrequency(frequency int) *FrequncyListDetails {

	freqDetails := &FrequncyListDetails{
		head: &LFUListNode{
			key: head,
		},
		tail: &LFUListNode{
			key: tail,
		},
	}

	freqDetails.head.next = freqDetails.tail
	freqDetails.tail.prev = freqDetails.head

	c.FrequncyList[frequency] = freqDetails
	return freqDetails
}

func (c *LFUCache) removeNode(node *LFUListNode) {

	node.next.prev = node.prev
	node.prev.next = node.next
	c.removefrequencyIfEmpty()

}

func (c *LFUCache) addToFront(node *LFUListNode) {

	frequency := node.frequency
	list := c.FrequncyList[frequency]
	if list == nil {
		c.initialiseNewfrequency(frequency)
		list = c.FrequncyList[frequency]
	}
	node.prev = list.head
	node.next = list.head.next

	node.next.prev = node
	list.head.next = node
}

func (c *LFUCache) removefrequencyIfEmpty() {

	frequency := c.Leastfrequency
	list := c.FrequncyList[frequency]
	if list.head.next.key == tail {
		delete(c.FrequncyList, frequency)
		frequency++
		for {
			if _, ok := c.FrequncyList[frequency]; ok {
				c.Leastfrequency = frequency
				break
			}

			if len(c.FrequncyList) == 0 {
				c.Leastfrequency = 0
				break
			}
		}
	}
}
func (c *LFUCache) incrementrequestCount() {
	c.requestCount = c.requestCount + 1
}

func (c *LFUCache) incrementhitCount() {
	c.hitCount = c.hitCount + 1
}

//
// FIFO Cache Implementation
//

type FIFOCache struct {
	// TODO: Add necessary fields for FIFO implementation
	// Hint: Use a queue or circular buffer
	mu           sync.RWMutex
	head         *ListNode
	tail         *ListNode
	hashMap      map[string]*ListNode
	cap          int
	requestCount int
	hitCount     int
}

// NewFIFOCache creates a new FIFO cache with the specified capacity
func NewFIFOCache(capacity int) *FIFOCache {
	// TODO: Implement FIFO cache constructor
	if capacity < 0 {
		capacity = 0
	}

	cache := &FIFOCache{
		head: &ListNode{
			key: head,
		},
		tail: &ListNode{
			key: tail,
		},
		hashMap: make(map[string]*ListNode, capacity),
		cap:     capacity,
	}

	cache.head.next = cache.tail
	cache.tail.prev = cache.head

	return cache
}

func (c *FIFOCache) Get(key string) (interface{}, bool) {
	// TODO: Implement FIFO get operation
	// Note: Get operations don't affect eviction order in FIFO
	c.incrementrequestCount()

	val, ok := c.hashMap[key]
	if !ok {
		return nil, false
	}

	c.incrementhitCount()

	return val.value, true
}

func (c *FIFOCache) Put(key string, value interface{}) {
	// TODO: Implement FIFO put operation
	// Should evict first-in item if at capacity

	if c.cap == 0 {
		return
	}

	existing, ok := c.hashMap[key]
	if ok {
		existing.value = value
		return
	}

	if len(c.hashMap) == c.cap {
		toBeDeleted := c.head.next
		c.removeNode(toBeDeleted)
		delete(c.hashMap, toBeDeleted.key)
	}

	node := &ListNode{
		key:   key,
		value: value,
	}
	c.addToLast(node)

	c.hashMap[key] = node
}

func (c *FIFOCache) Delete(key string) bool {
	// TODO: Implement delete operation
	node, exists := c.hashMap[key]
	if !exists {
		return false
	}

	delete(c.hashMap, node.key)

	c.removeNode(node)
	return true
}

func (c *FIFOCache) Clear() {
	// TODO: Implement clear operation

	c.head.next = c.tail
	c.tail.prev = c.head

	clear(c.hashMap)
}

func (c *FIFOCache) Size() int {
	// TODO: Return current cache size

	return len(c.hashMap)
}

func (c *FIFOCache) Capacity() int {
	// TODO: Return cache capacity
	return c.cap
}

func (c *FIFOCache) HitRate() float64 {
	// TODO: Calculate and return hit rate
	if c.requestCount == 0 {
		return 0
	}
	return float64(c.hitCount) / float64(c.requestCount)
}

func (c *FIFOCache) addToLast(node *ListNode) {

	node.next = c.tail
	node.prev = c.tail.prev

	node.prev.next = node
	c.tail.prev = node
}

func (c *FIFOCache) removeNode(node *ListNode) {

	node.prev.next = node.next
	node.next.prev = node.prev
}

func (c *FIFOCache) incrementrequestCount() {
	c.requestCount = c.requestCount + 1
}

func (c *FIFOCache) incrementhitCount() {
	c.hitCount = c.hitCount + 1
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
	return &ThreadSafeCache{
		cache: cache,
	}
}

func (c *ThreadSafeCache) Get(key string) (interface{}, bool) {
	// TODO: Implement thread-safe get operation
	// Hint: Use read lock for better performance
	defer c.mu.Unlock()
	c.mu.Lock()
	return c.cache.Get(key)
}

func (c *ThreadSafeCache) Put(key string, value interface{}) {
	// TODO: Implement thread-safe put operation
	// Hint: Use write lock
	defer c.mu.Unlock()
	c.mu.Lock()
	c.cache.Put(key, value)
}

func (c *ThreadSafeCache) Delete(key string) bool {
	// TODO: Implement thread-safe delete operation
	defer c.mu.Unlock()
	c.mu.Lock()
	return c.cache.Delete(key)
}

func (c *ThreadSafeCache) Clear() {
	// TODO: Implement thread-safe clear operation
	defer c.mu.Unlock()
	c.mu.Lock()
	c.cache.Clear()
}

func (c *ThreadSafeCache) Size() int {
	// TODO: Implement thread-safe size operation
	defer c.mu.RUnlock()
	c.mu.RLock()
	return c.cache.Size()
}

func (c *ThreadSafeCache) Capacity() int {
	// TODO: Implement thread-safe capacity operation
	defer c.mu.RUnlock()
	c.mu.RLock()
	return c.cache.Capacity()
}

func (c *ThreadSafeCache) HitRate() float64 {
	// TODO: Implement thread-safe hit rate operation
	defer c.mu.RUnlock()
	c.mu.RLock()

	return c.cache.HitRate()
}

//
// Cache Factory Functions
//

// NewCache creates a cache with the specified policy and capacity
func NewCache(policy CachePolicy, capacity int) Cache {
	// TODO: Implement cache factory
	// Should create appropriate cache type based on policy

	if capacity < 0 {
		capacity = 0
	}

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
		return NewFIFOCache(capacity)
	}
}

// NewThreadSafeCacheWithPolicy creates a thread-safe cache with the specified policy
func NewThreadSafeCacheWithPolicy(policy CachePolicy, capacity int) Cache {
	// TODO: Implement thread-safe cache factory
	// Should create cache with policy and wrap it with thread safety
	cache := NewCache(policy, capacity)
	if cache != nil {
		cache = NewThreadSafeCache(cache)
	}
	return cache
}
