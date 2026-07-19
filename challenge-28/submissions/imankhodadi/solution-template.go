// Challenge 28: Cache Implementation with Multiple Eviction Policies
package cache // fixed for submission

import (
	"fmt"
	"sync"
)

type Cache interface {
	Get(key string) (value interface{}, found bool)
	Put(key string, value interface{})
	Delete(key string) bool
	Clear()
	Size() int     // current entries
	Capacity() int // fixed maximum
	HitRate() float64
}
type CacheMetrics struct {
	hits      int64 //requests served from cache
	misses    int64 //requests not in cache
	evictions int64 //Number of items evicted
}
type CachePolicy int

const (
	LRU CachePolicy = iota
	LFU
	FIFO
)

type Node struct {
	key       string
	value     interface{}
	prev      *Node
	next      *Node
	frequency int
}

/*
LRU (Least Recently Used): LRU evicts the item that was accessed least recently

	Operating system page replacement
	CPU cache management
	Web browser cache
	Database buffer pools

Advantages:

	Good temporal locality performance
	Intuitive eviction strategy
	Works well for most general-purpose scenarios

Disadvantages:

	Doesn't consider access frequency
	Can be affected by sequential scans that destroy cache locality
*/
type LRUCache struct {
	capacity int
	size     int
	cache    map[string]*Node
	head     *Node // Most recently used
	tail     *Node // Least recently used
	metrics  CacheMetrics
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{capacity: capacity, cache: map[string]*Node{}, metrics: CacheMetrics{}}
}
func (c *LRUCache) moveToFront(node *Node) {
	if node == c.head {
		return
	}
	if node.prev != nil {
		node.prev.next = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	}
	if node == c.tail {
		c.tail = node.next
	}
	node.prev = c.head
	node.next = nil
	if c.head != nil {
		c.head.next = node
	}
	c.head = node
	if c.tail == nil {
		c.tail = node
	}
}
func (c *LRUCache) Get(key string) (interface{}, bool) {
	node, ok := c.cache[key]
	if !ok {
		c.metrics.misses++
		return nil, false
	}
	c.metrics.hits++
	c.moveToFront(node)
	return node.value, true
}
func (c *LRUCache) Put(key string, value interface{}) {
	if c.capacity <= 0 {
		return
	}
	if node, ok := c.cache[key]; ok {
		node.value = value
		c.moveToFront(node)
		return
	}
	node := &Node{key: key, value: value}
	c.cache[key] = node
	if c.head == nil {
		c.head = node
		c.tail = node
		c.size = 1
		return
	}
	node.prev = c.head
	c.head.next = node
	c.head = node
	if c.size < c.capacity {
		c.size++
		return
	}
	victim := c.tail
	delete(c.cache, victim.key)
	c.tail = victim.next
	if c.tail != nil {
		c.tail.prev = nil
	}
	c.metrics.evictions++
}

func (c *LRUCache) Delete(key string) bool {
	node, ok := c.cache[key]
	if !ok {
		return false
	}
	c.size--
	if node.prev != nil {
		node.prev.next = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	}
	if node == c.head {
		c.head = node.prev
	}
	if node == c.tail {
		c.tail = node.next
	}
	delete(c.cache, key)
	return true
}
func (c *LRUCache) Clear() {
	clear(c.cache)
	c.head = nil
	c.size = 0
	c.tail = nil

}
func (c *LRUCache) Size() int {
	return c.size
}
func (c *LRUCache) Capacity() int {
	return c.capacity
}
func (c *LRUCache) HitRate() float64 {
	total := c.metrics.hits + c.metrics.misses
	if total == 0 {
		return 0
	}
	return float64(c.metrics.hits) / float64(total)
}

/*
LFU Cache Implementation: LFU evicts the item that has been accessed the fewest times.
Maintain a frequency counter for each item. Use a min-heap or frequency buckets for efficient eviction
On access: increment frequency counter
On eviction: remove item with lowest frequency

	Long-running applications with stable access patterns
	Scientific computing with repeated data access
	CDN systems

Advantages:

	Excellent for workloads with clear hot data
	Adapts well to changing access patterns over time
	Good for scenarios where some data is accessed much more frequently

Disadvantages

	More complex implementation
	New items are immediately evicted if cache is full
	Frequency counts can become stale over time
*/
type FreqGroup struct {
	freq int
	head *Node
	tail *Node
}
type LFUCache struct {
	capacity   int
	size       int
	minFreq    int
	cache      map[string]*Node
	freqGroups map[int]*FreqGroup // frequency -> list of nodes
	metrics    CacheMetrics
}

func NewLFUCache(capacity int) *LFUCache {
	return &LFUCache{
		capacity:   capacity,
		cache:      make(map[string]*Node),
		freqGroups: make(map[int]*FreqGroup),
	}
}
func (c *LFUCache) evictLFU() {
	minGroup := c.freqGroups[c.minFreq]
	if minGroup == nil || minGroup.tail == nil {
		return
	}
	victim := minGroup.tail
	delete(c.cache, victim.key)
	minGroup.tail = victim.next
	if minGroup.tail != nil {
		minGroup.tail.prev = nil
	} else {
		minGroup.head = nil
	}
	c.metrics.evictions++
	c.size--
	if minGroup.head == nil && minGroup.tail == nil {
		delete(c.freqGroups, c.minFreq)
	}
}
func (c *LFUCache) promote(node *Node) {
	oldFreq := node.frequency
	oldGroup := c.freqGroups[oldFreq]
	if node.prev != nil {
		node.prev.next = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	}
	if oldGroup.head == node {
		oldGroup.head = node.prev
	}
	if oldGroup.tail == node {
		oldGroup.tail = node.next
	}
	node.prev = nil
	node.next = nil
	node.frequency++

	group := c.freqGroups[node.frequency]
	if group == nil {
		group = &FreqGroup{freq: node.frequency}
		c.freqGroups[node.frequency] = group
	}
	if group.head == nil {
		group.head = node
		group.tail = node
	} else {
		group.head.next = node
		node.prev = group.head
		group.head = node
	}
	if oldFreq == c.minFreq && oldGroup.head == nil && oldGroup.tail == nil {
		c.minFreq++
	}
	c.pruneFreqGroup(oldFreq)
}

func (c *LFUCache) Get(key string) (interface{}, bool) {
	node, ok := c.cache[key]
	if !ok {
		c.metrics.misses++
		return nil, false
	}
	c.metrics.hits++
	c.promote(node)
	return node.value, true
}

func (c *LFUCache) Put(key string, value interface{}) {
	if c.capacity <= 0 {
		return
	}
	if node, ok := c.cache[key]; ok {
		node.value = value
		c.promote(node)
		return
	}
	if c.size == c.capacity {
		c.evictLFU()
	}
	node := &Node{key: key, value: value, frequency: 1}
	c.cache[key] = node
	group := c.freqGroups[1]
	if group == nil {
		group = &FreqGroup{freq: 1}
		c.freqGroups[1] = group
	}
	if group.head == nil {
		group.head = node
		group.tail = node
	} else {
		group.head.next = node
		node.prev = group.head
		group.head = node
	}
	c.minFreq = 1
	c.size++
}

func (c *LFUCache) Delete(key string) bool {
	node, ok := c.cache[key]
	if !ok {
		return false
	}
	c.size--
	group := c.freqGroups[node.frequency]
	Prev, Next := node.prev, node.next
	if Prev != nil {
		Prev.next = Next
	} else if group != nil {
		group.tail = Next
	}
	if Next != nil {
		Next.prev = Prev
	} else if group != nil {
		group.head = Prev
	}
	node.next = nil
	node.prev = nil
	node.value = nil
	delete(c.cache, key)
	c.recomputeMinFreq()
	return true
}
func (c *LFUCache) pruneFreqGroup(freq int) {
	group := c.freqGroups[freq]
	if group != nil && group.head == nil && group.tail == nil {
		delete(c.freqGroups, freq)
	}
}

func (c *LFUCache) recomputeMinFreq() {
	c.minFreq = 0
	for freq, group := range c.freqGroups {
		if group == nil || group.head == nil {
			delete(c.freqGroups, freq)
			continue
		}
		if c.minFreq == 0 || freq < c.minFreq {
			c.minFreq = freq
		}
	}
}

func (c *LFUCache) Clear() {
	clear(c.freqGroups)
	clear(c.cache)
	c.size = 0
	c.minFreq = 0
}
func (c *LFUCache) Size() int {
	return c.size
}
func (c *LFUCache) Capacity() int {
	return c.capacity
}
func (c *LFUCache) HitRate() float64 {
	total := c.metrics.hits + c.metrics.misses
	if total == 0 {
		return 0
	}
	return float64(c.metrics.hits) / float64(total)
}

/*
FIFO evicts the oldest item in the cache, regardless of access patterns.
Simple caching scenarios When access patterns are unknown

	Embedded systems with memory constraints

Advantages:

	Simple to implement and understand
	Predictable behavior
	No access pattern tracking needed

Disadvantages:

	Ignores access patterns completely
	May evict frequently used items
	Generally poor cache hit rates
*/
type FIFOCache struct {
	capacity int
	size     int
	cache    map[string]*Node
	head     *Node // Newest item
	tail     *Node // Oldest item
	metrics  CacheMetrics
}

func NewFIFOCache(capacity int) *FIFOCache {
	return &FIFOCache{
		capacity: capacity,
		cache:    make(map[string]*Node),
	}
}
func (c *FIFOCache) Get(key string) (interface{}, bool) {
	node, ok := c.cache[key]
	if !ok {
		c.metrics.misses++
		return nil, false
	}
	c.metrics.hits++
	return node.value, true
}

func (c *FIFOCache) Put(key string, value interface{}) {
	if c.capacity <= 0 {
		return
	}
	if existing, ok := c.cache[key]; ok {
		existing.value = value
		return
	}
	node := &Node{key: key, value: value}
	c.cache[key] = node
	if c.head == nil {
		c.head = node
		c.tail = node
		c.size = 1
		return
	}
	c.head.next = node
	node.prev = c.head
	c.head = node
	if c.size < c.capacity {
		c.size++
		return
	}
	victim := c.tail
	delete(c.cache, victim.key)
	c.tail = victim.next
	if c.tail != nil {
		c.tail.prev = nil
	}
	c.metrics.evictions++
}

func (c *FIFOCache) Delete(key string) bool {
	node, ok := c.cache[key]
	if !ok {
		return false
	}
	c.size--
	if node.prev != nil {
		node.prev.next = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	}
	if node == c.head {
		c.head = node.prev
	}
	if node == c.tail {
		c.tail = node.next
	}
	delete(c.cache, key)
	return true
}
func (c *FIFOCache) Clear() {
	clear(c.cache)
	c.head = nil
	c.tail = nil
	c.size = 0
}
func (c *FIFOCache) Size() int {
	return c.size
}

func (c *FIFOCache) Capacity() int {
	return c.capacity
}
func (c *FIFOCache) HitRate() float64 {
	total := c.metrics.hits + c.metrics.misses
	if total == 0 {
		return 0
	}
	return float64(c.metrics.hits) / float64(total)
}

// Thread-Safe Cache Wrapper
type ThreadSafeCache struct {
	cache Cache
	mu    sync.RWMutex
}

func NewThreadSafeCache(cache Cache) *ThreadSafeCache {
	return &ThreadSafeCache{cache: cache}
}
func (c *ThreadSafeCache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.cache.Get(key)
}
func (c *ThreadSafeCache) Put(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache.Put(key, value)
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
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cache.Capacity()
}
func (c *ThreadSafeCache) HitRate() float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cache.HitRate()
}

// Cache Factory Functions
func NewCache(policy CachePolicy, capacity int) Cache {
	switch policy {
	case LRU:
		return NewLRUCache(capacity)
	case LFU:
		return NewLFUCache(capacity)
	case FIFO:
		return NewFIFOCache(capacity)
	default:
		return nil
	}
}

// NewThreadSafeCacheWithPolicy creates a thread-safe cache with the specified policy
func NewThreadSafeCacheWithPolicy(policy CachePolicy, capacity int) Cache {
	cache := NewCache(policy, capacity)
	if cache == nil {
		return nil
	}
	return NewThreadSafeCache(cache)
}
func main() {
	cache := NewLFUCache(2)
	cache.Put("a", 1)
	cache.Put("b", 2)
	// Access "a" multiple times to increase its frequency
	cache.Get("a")
	cache.Get("a")
	// Now "a" has frequency 3, "b" has frequency 1
	// Add new item, should evict "b" (least frequently used)
	cache.Put("c", 3)
	// "b" should be evicted
	_, found := cache.Get("b")
	if found {
		fmt.Println("Expected 'b' to be evicted (least frequently used)")
	}
	// "a" and "c" should still be present
	value, found := cache.Get("a")
	if !found || value != 1 {
		fmt.Printf("Expected 'a' to be present with value 1, got (%v, %v)", value, found)
	}
	value, found = cache.Get("c")
	if !found || value != 3 {
		fmt.Printf("Expected 'c' to be present with value 3, got (%v, %v)", value, found)
	}
}
