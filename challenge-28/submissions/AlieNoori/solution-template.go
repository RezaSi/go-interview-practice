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

type node struct {
	value any
	next  *node
	prev  *node
}

type LRUCache struct {
	capacity, length int
	hits, misses     int
	head, tail       *node
	lookup           map[string]*node
	reverseLookup    map[*node]string
}

// NewLRUCache creates a new LRU cache with the specified capacity
func NewLRUCache(capacity int) *LRUCache {
	if capacity <= 0 {
		return nil
	}

	return &LRUCache{
		capacity:      capacity,
		length:        0,
		tail:          nil,
		head:          nil,
		lookup:        make(map[string]*node),
		reverseLookup: make(map[*node]string),
	}
}

func (c *LRUCache) Get(key string) (interface{}, bool) {
	n, exists := c.lookup[key]
	if !exists {
		c.misses++
		return nil, false
	}

	c.hits++

	c.detach(n)
	c.prepend(n)

	return n.value, true
}

func (c *LRUCache) Put(key string, value interface{}) {
	// TODO: Implement LRU put operation
	// Should add new item to front and evict least recently used if at capacity
	n, exists := c.lookup[key]
	if !exists {
		n = &node{value: value}
		c.length++
		c.prepend(n)
		c.trimCache()

		c.lookup[key] = n
		c.reverseLookup[n] = key
	} else {
		c.detach(n)
		c.prepend(n)
		n.value = value
	}
}

func (c *LRUCache) Delete(key string) bool {
	node, ok := c.lookup[key]
	if !ok {
		return false
	}

	c.detach(node)

	delete(c.lookup, key)
	delete(c.reverseLookup, node)

	c.length--

	return true
}

func (c *LRUCache) Clear() {
	c.head = nil
	c.tail = nil
	clear(c.lookup)
	clear(c.reverseLookup)
	c.length = 0
}

func (c *LRUCache) Size() int {
	return c.length
}

func (c *LRUCache) Capacity() int {
	return c.capacity
}

func (c *LRUCache) HitRate() float64 {
	total := c.hits + c.misses
	if total == 0 {
		return 0
	}
	return float64(c.hits) / float64(total)
}

func (c *LRUCache) detach(n *node) {
	if n.prev != nil {
		n.prev.next = n.next
	}
	if n.next != nil {
		n.next.prev = n.prev
	}

	if c.head == n {
		c.head = c.head.next
	}

	if c.tail == n {
		c.tail = c.tail.prev
	}

	n.next = nil
	n.prev = nil
}

func (c *LRUCache) prepend(n *node) {
	if c.head == nil {
		c.head, c.tail = n, n
		return
	}

	n.next = c.head
	c.head.prev = n
	c.head = n
}

func (c *LRUCache) trimCache() {
	if c.length <= c.capacity {
		return
	}

	tail := c.tail
	c.detach(c.tail)

	key := c.reverseLookup[tail]
	delete(c.lookup, key)
	delete(c.reverseLookup, tail)
	c.length--
}

//
// LFU Cache Implementation
//

type DLL struct {
	length int
	tail   *node
	head   *node
}

func newDoublyLinkedList() *DLL {
	return &DLL{
		length: 0,
		head:   nil,
		tail:   nil,
	}
}

func (l *DLL) detach(n *node) {
	l.length--
	if n.prev != nil {
		n.prev.next = n.next
	}
	if n.next != nil {
		n.next.prev = n.prev
	}

	if l.head == n {
		l.head = l.head.next
	}

	if l.tail == n {
		l.tail = l.tail.prev
	}

	n.next = nil
	n.prev = nil
}

func (l *DLL) prepend(n *node) {
	l.length++
	if l.head == nil {
		l.head, l.tail = n, n
		return
	}

	n.next = l.head
	l.head.prev = n
	l.head = n
}

func (l *DLL) trimTail() *node {
	if l.tail == nil {
		return nil
	}

	tail := l.tail
	if l.tail.prev != nil {
		l.tail.prev.next = nil
	}
	prev := l.tail.prev
	l.tail.prev = nil
	l.tail = prev
	return tail
}

func (l *DLL) size() int {
	return l.length
}

type LFUCache struct {
	minFreq          int
	hits, misses     int
	capacity, length int
	lookup           map[string]*node
	reverseLookup    map[*node]string
	keyFreq          map[string]int
	freqList         map[int]*DLL
}

// NewLFUCache creates a new LFU cache with the specified capacity
func NewLFUCache(capacity int) *LFUCache {
	if capacity <= 0 {
		return nil
	}

	return &LFUCache{
		length:        0,
		capacity:      capacity,
		minFreq:       0,
		misses:        0,
		hits:          0,
		lookup:        make(map[string]*node),
		reverseLookup: make(map[*node]string),
		keyFreq:       make(map[string]int),
		freqList:      make(map[int]*DLL),
	}
}

func (c *LFUCache) Get(key string) (interface{}, bool) {
	n, exists := c.lookup[key]
	if !exists {
		c.misses++
		return nil, false
	}

	c.hits++
	c.incrementFreq(n)

	return n.value, true
}

func (c *LFUCache) Put(key string, value interface{}) {
	n, exists := c.lookup[key]
	if exists {
		n.value = value
		c.incrementFreq(n)
	} else {
		freq := 1
		c.minFreq = 1
		n := &node{value: value}
		c.keyFreq[key] = freq
		c.lookup[key] = n
		c.reverseLookup[n] = key
		list := c.freqList[freq]
		if list == nil {
			list = newDoublyLinkedList()
			c.freqList[freq] = list
		}
		list.prepend(n)
		c.length++
		c.trimCache()
	}
}

func (c *LFUCache) Delete(key string) bool {
	n, exists := c.lookup[key]
	if !exists {
		return false
	}

	freq := c.keyFreq[key]

	list := c.freqList[freq]
	list.detach(n)

	if list.size() == 0 {
		if c.minFreq == freq {
			c.minFreq++
		}
		delete(c.freqList, freq)
	}

	delete(c.keyFreq, key)
	delete(c.lookup, key)
	delete(c.reverseLookup, n)

	c.length--

	return true
}

func (c *LFUCache) Clear() {
	c.length = 0
	c.minFreq = 0
	clear(c.lookup)
	clear(c.reverseLookup)
	clear(c.keyFreq)
	clear(c.freqList)
}

func (c *LFUCache) Size() int {
	return c.length
}

func (c *LFUCache) Capacity() int {
	return c.capacity
}

func (c *LFUCache) HitRate() float64 {
	total := c.hits + c.misses
	if total == 0 {
		return 0
	}
	return float64(c.hits) / float64(total)
}

func (c *LFUCache) incrementFreq(n *node) {
	key := c.reverseLookup[n]
	freq := c.keyFreq[key]
	list := c.freqList[freq]
	list.detach(n)
	if list.size() == 0 {
		if c.minFreq == freq {
			c.minFreq++
		}
		delete(c.freqList, freq)
	}
	newFreq := freq + 1
	c.keyFreq[key] = newFreq
	list = c.freqList[newFreq]
	if list == nil {
		list = newDoublyLinkedList()
		c.freqList[newFreq] = list
	}
	list.prepend(n)
}

func (c *LFUCache) trimCache() {
	if c.length <= c.capacity {
		return
	}

	list := c.freqList[c.minFreq]
	n := list.trimTail()
	if n != nil {
		delete(c.lookup, c.reverseLookup[n])
		delete(c.keyFreq, c.reverseLookup[n])
		delete(c.reverseLookup, n)
	}
	c.length--
}

//
// FIFO Cache Implementation
//

type FIFOCache struct {
	capacity, length int
	hits, misses     int
	head, tail       *node
	lookup           map[string]*node
	reverseLookup    map[*node]string
}

// NewFIFOCache creates a new FIFO cache with the specified capacity
func NewFIFOCache(capacity int) *FIFOCache {
	if capacity <= 0 {
		return nil
	}
	return &FIFOCache{
		length:        0,
		capacity:      capacity,
		hits:          0,
		misses:        0,
		head:          nil,
		tail:          nil,
		reverseLookup: make(map[*node]string),
		lookup:        make(map[string]*node),
	}
}

func (c *FIFOCache) Get(key string) (interface{}, bool) {
	n, exist := c.lookup[key]
	if !exist {
		c.misses++
		return nil, false
	}

	c.hits++

	return n.value, true
}

func (c *FIFOCache) Put(key string, value interface{}) {
	n, exists := c.lookup[key]
	if exists {
		n.value = value
	} else {
		// Should evict first-in item if at capacity
		n = &node{value: value}

		c.lookup[key] = n
		c.reverseLookup[n] = key

		c.enqueue(n)

		c.trimCache()
	}
}

func (c *FIFOCache) Delete(key string) bool {
	n, exists := c.lookup[key]
	if !exists {
		return false
	}

	delete(c.lookup, key)
	delete(c.reverseLookup, n)

	c.dequeue()
	return true
}

func (c *FIFOCache) Clear() {
	c.length = 0
	c.hits = 0
	c.misses = 0
	c.head = nil
	c.tail = nil
	clear(c.lookup)
	clear(c.reverseLookup)
}

func (c *FIFOCache) Size() int {
	return c.length
}

func (c *FIFOCache) Capacity() int {
	return c.capacity
}

func (c *FIFOCache) HitRate() float64 {
	total := c.hits + c.misses
	if total == 0 {
		return 0
	}
	return float64(c.hits) / float64(total)
}

func (c *FIFOCache) enqueue(n *node) {
	c.length++
	if c.tail == nil {
		c.tail = n
		c.head = n
		return
	}

	n.prev = c.tail
	c.tail.next = n
	c.tail = n
}

func (c *FIFOCache) dequeue() *node {
	if c.head == nil {
		c.tail = nil
		c.length = 0
		return nil
	}

	c.length--
	n := c.head
	c.head = c.head.next
	c.head.prev = nil

	return n
}

func (c *FIFOCache) trimCache() {
	if c.length <= c.capacity {
		return
	}

	n := c.dequeue()
	key := c.reverseLookup[n]
	delete(c.lookup, key)
	delete(c.reverseLookup, n)
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
	if cache == nil {
		return nil
	}
	return &ThreadSafeCache{
		cache: cache,
		mu:    sync.RWMutex{},
	}
}

func (c *ThreadSafeCache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.cache.Get(key)
}

func (c *ThreadSafeCache) Put(key string, value interface{}) {
	c.mu.Lock()
	c.cache.Put(key, value)
	defer c.mu.Unlock()
}

func (c *ThreadSafeCache) Delete(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.cache.Delete(key)
}

func (c *ThreadSafeCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache.Clear()
}

func (c *ThreadSafeCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cache.Size()
}

func (c *ThreadSafeCache) Capacity() int {
	return c.cache.Capacity()
}

func (c *ThreadSafeCache) HitRate() float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cache.HitRate()
}

//
// Cache Factory Functions
//

// NewCache creates a cache with the specified policy and capacity
func NewCache(policy CachePolicy, capacity int) Cache {
	var cache Cache
	switch policy {
	case LRU:
		cache = NewLRUCache(capacity)
	case LFU:
		cache = NewLFUCache(capacity)
	case FIFO:
		cache = NewFIFOCache(capacity)
	default:
		cache = nil
	}

	return cache
}

// NewThreadSafeCacheWithPolicy creates a thread-safe cache with the specified policy
func NewThreadSafeCacheWithPolicy(policy CachePolicy, capacity int) Cache {
	cache := NewCache(policy, capacity)
	if cache == nil {
		return nil
	}
	return NewThreadSafeCache(cache)
}
