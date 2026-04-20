// Package cache provides implementations of common cache eviction policies:
// LRU (Least Recently Used), LFU (Least Frequently Used), and FIFO (First In First Out).
// All implementations satisfy the Cache interface and can be wrapped with ThreadSafeCache
// for concurrent use. Every operation runs in O(1) time.
package cache

import (
	"reflect"
	"sync"
)

// Cache defines the contract that all cache implementations must satisfy.
type Cache interface {
	// Get retrieves a value by key. Returns the value and true if found,
	// or nil and false if the key does not exist.
	Get(key string) (value interface{}, found bool)

	// Put stores a key-value pair. If the cache is at capacity it evicts
	// one entry according to the implementation's policy before inserting.
	Put(key string, value interface{})

	// Delete removes the entry for key. Returns true if the key existed.
	Delete(key string) bool

	// Clear removes all entries from the cache.
	Clear()

	// Size returns the current number of items in the cache.
	Size() int

	// Capacity returns the maximum number of items the cache can hold.
	Capacity() int

	// HitRate returns the cache hit rate as a float between 0 and 1.
	HitRate() float64
}

// CachePolicy selects the eviction strategy used when constructing a cache.
type CachePolicy int

const (
	LRU  CachePolicy = iota // Evicts the least recently accessed item.
	LFU                     // Evicts the least frequently accessed item.
	FIFO                    // Evicts the item that was inserted first.
)

// ─── LRU Cache ───────────────────────────────────────────────────────────────

// Metrics holds access counters used to compute the hit rate.
type Metrics struct {
	hits      int
	misses    int
	evictions int
}

// LRUCache implements Cache with the Least Recently Used eviction policy.
// It keeps a doubly-linked list (most-recent at head, least-recent at tail)
// backed by a hash map for O(1) node lookup and repositioning.
type LRUCache struct {
	capacity int
	size     int
	cache    map[string]*LRUNode
	metrics  Metrics
	head     *LRUNode // most recently used
	tail     *LRUNode // least recently used — next eviction victim
}

// LRUNode is a doubly-linked list node that stores a single cache entry.
type LRUNode struct {
	key   string
	value interface{}
	next  *LRUNode
	prev  *LRUNode
}

// NewLRUCache returns a new LRUCache with the given capacity.
// Returns nil if capacity is less than 1.
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

func newLRUNode(key string, value interface{}) *LRUNode {
	return &LRUNode{
		key:   key,
		value: value,
	}
}

// evictLRUNode removes the tail node (least recently used entry).
func (c *LRUCache) evictLRUNode() bool {
	if c.size == 0 {
		return false
	}
	node := c.tail
	return c.removeLRUNode(node)
}

// removeLRUNode unlinks node from the doubly-linked list and deletes its map entry.
func (c *LRUCache) removeLRUNode(node *LRUNode) bool {
	if c.size == 0 || node == nil {
		return false
	}

	if c.size == 1 {
		// Last node in list — reset both sentinels.
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
	// Middle node: relink its neighbors directly.
	node.prev.next = node.next
	node.next.prev = node.prev
	node.next = nil
	node.prev = nil
	delete(c.cache, node.key)
	c.size--
	return true
}

// moveToFront repositions node at the head of the list, marking it as most recently used.
// If the node is already tracked in the cache map it is unlinked from its current
// position before being prepended.
func (c *LRUCache) moveToFront(node *LRUNode) bool {
	if node == nil {
		return false
	}
	if c.head == node {
		// Already at the front — nothing to do.
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

// Get returns the value for key and moves the accessed node to the MRU position.
func (c *LRUCache) Get(key string) (interface{}, bool) {
	if node, exists := c.cache[key]; exists {
		c.moveToFront(node)
		c.metrics.hits++
		return node.value, true
	}
	c.metrics.misses++
	return nil, false
}

// Put inserts or updates a key-value pair.
// If the cache is full, the least recently used entry is evicted before insertion.
// Updating an existing key refreshes its recency without changing the cache size.
func (c *LRUCache) Put(key string, value interface{}) {
	if node, exists := c.cache[key]; exists {
		// Update value in-place and refresh its recency position.
		node.value = value
		c.moveToFront(node)
		return
	}
	node := newLRUNode(key, value)
	if c.size == c.capacity {
		c.evictLRUNode()
		c.metrics.evictions++
	}
	c.moveToFront(node)
}

// Delete removes the entry for key and returns true if the key existed.
func (c *LRUCache) Delete(key string) bool {
	if node, exists := c.cache[key]; exists {
		c.removeLRUNode(node)
		return true
	}
	return false
}

// Clear removes all entries and resets metrics.
func (c *LRUCache) Clear() {
	c.size = 0
	c.cache = make(map[string]*LRUNode, c.capacity)
	c.head = nil
	c.tail = nil
	c.metrics = Metrics{}
}

// Size returns the number of entries currently in the cache.
func (c *LRUCache) Size() int {
	return c.size
}

// Capacity returns the maximum number of entries the cache can hold.
func (c *LRUCache) Capacity() int {
	return c.capacity
}

// HitRate returns the fraction of Get calls that found a key (between 0.0 and 1.0).
// Returns 0.0 when no Get calls have been made yet.
func (c *LRUCache) HitRate() float64 {
	requests := c.metrics.misses + c.metrics.hits
	if requests == 0 {
		return 0.0
	}
	return float64(c.metrics.hits) / float64(requests)
}

// ─── LFU Cache ───────────────────────────────────────────────────────────────

// LFUCache implements Cache with the Least Frequently Used eviction policy.
// Each entry tracks its access count. When the cache is full the entry with the
// lowest count is evicted; ties are broken by evicting the oldest entry within
// that frequency bucket (LRU order within the same frequency).
//
// Internally, entries are organised into per-frequency doubly-linked lists
// (FreqGroup). minFreq records the current minimum for O(1) victim selection.
type LFUCache struct {
	capacity   int
	size       int
	cache      map[string]*LFUNode
	freqGroups map[int]*FreqGroup // frequency → bucket of nodes sharing that count
	minFreq    int
	metrics    Metrics
}

// LFUNode is a doubly-linked list node used inside a FreqGroup.
type LFUNode struct {
	key   string
	value interface{}
	freq  int // access count; incremented on every Get or Put-update
	next  *LFUNode
	prev  *LFUNode
}

// FreqGroup is a doubly-linked list of all nodes that share the same access
// frequency. New nodes are prepended to the head; eviction removes from the
// tail so the oldest entry within a tied frequency is evicted first.
type FreqGroup struct {
	freq int
	size int
	head *LFUNode
	tail *LFUNode
}

// NewLFUCache returns a new LFUCache with the given capacity.
// Returns nil if capacity is less than or equal to 0.
func NewLFUCache(capacity int) *LFUCache {
	if capacity <= 0 {
		return nil
	}
	return &LFUCache{
		capacity:   capacity,
		cache:      make(map[string]*LFUNode, capacity),
		freqGroups: make(map[int]*FreqGroup),
	}
}

func newLFUNode(key string, value interface{}) *LFUNode {
	return &LFUNode{
		key:   key,
		value: value,
		freq:  1, // every new entry starts with frequency 1
	}
}

func newFreqGroup(freq int) *FreqGroup {
	return &FreqGroup{
		freq: freq,
	}
}

// evictLFUNode removes the tail of the minFreq bucket (least frequent, oldest on ties).
// If the bucket becomes empty after removal, minFreq is incremented.
func (c *LFUCache) evictLFUNode() bool {
	if bucket, exists := c.freqGroups[c.minFreq]; exists {
		c.removeLFUNode(bucket.tail)
		// Advance minFreq if the bucket was just emptied.
		if _, exists = c.freqGroups[c.minFreq]; !exists {
			c.minFreq++
		}
		return true
	}
	return false
}

// removeLFUNode unlinks node from its FreqGroup list, removes it from the cache map,
// and decrements the cache size. The FreqGroup is deleted when it becomes empty.
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
			// Removing the only node — drop the entire bucket.
			bucket.head = nil
			bucket.tail = nil
			delete(c.freqGroups, bucket.freq)
			delete(c.cache, node.key)
			c.size--
			return true
		}
		if bucket.head == node {
			bucket.head = node.next
			bucket.head.prev = nil
			node.next = nil
			delete(c.cache, node.key)
			bucket.size--
			c.size--
			return true
		}
		if bucket.tail == node {
			bucket.tail = node.prev
			bucket.tail.next = nil
			node.prev = nil
			delete(c.cache, node.key)
			bucket.size--
			c.size--
			return true
		}
		// Middle node: relink its neighbors directly.
		node.prev.next = node.next
		node.next.prev = node.prev
		node.next = nil
		node.prev = nil
		delete(c.cache, node.key)
		bucket.size--
		c.size--
		return true
	}
	return false
}

// addLFUNode prepends node to the head of its FreqGroup, creating the group if needed.
// It also registers the node in the cache map and increments the cache size.
func (c *LFUCache) addLFUNode(node *LFUNode) {
	if bucket, exists := c.freqGroups[node.freq]; exists {
		// Prepend: new node becomes the most-recent in this frequency bucket.
		node.next = bucket.head
		bucket.head = node
		node.next.prev = node
		node.prev = nil
		c.cache[node.key] = node
		bucket.size++
		c.size++
		return
	}
	// First node for this frequency — create a new bucket.
	bucket := newFreqGroup(node.freq)
	bucket.head = node
	bucket.tail = node
	node.next = nil
	node.prev = nil
	c.freqGroups[bucket.freq] = bucket
	c.cache[node.key] = node
	bucket.size++
	c.size++
}

// Get returns the value for key and increments its access frequency.
// The node is moved from its current frequency bucket to the freq+1 bucket.
func (c *LFUCache) Get(key string) (interface{}, bool) {

	if node, exists := c.cache[key]; exists {
		// Promote node: remove from current bucket, bump frequency, re-insert.
		c.removeLFUNode(node)
		if _, exists = c.freqGroups[node.freq]; !exists && node.freq == c.minFreq {
			// The minFreq bucket became empty; the promoted frequency is the new minimum.
			c.minFreq++
		}
		node.freq++
		c.addLFUNode(node)
		c.metrics.hits++
		return node.value, true
	}
	c.metrics.misses++
	return nil, false
}

// Put inserts or updates a key-value pair.
// Updating an existing key re-inserts it with an incremented frequency.
// Inserting a new key resets minFreq to 1.
// If the cache is full, the least-frequent (oldest on ties) entry is evicted first.
func (c *LFUCache) Put(key string, value interface{}) {

	var (
		node   *LFUNode
		exists bool
	)

	if node, exists = c.cache[key]; exists {
		// Update existing entry: remove from old freq bucket, bump frequency.
		oldFreq := node.freq
		c.removeLFUNode(node)
		if _, stillExists := c.freqGroups[oldFreq]; !stillExists && oldFreq == c.minFreq {
			c.minFreq++
		}
		node.value = value
		node.freq++
		c.addLFUNode(node)
		return
	} else {
		if c.size == c.capacity {
			if c.evictLFUNode() {
				c.metrics.evictions++
			}
		}
		// New entry: starts at frequency 1; a new entry is always the new minimum.
		node = newLFUNode(key, value)
		c.minFreq = 1
	}
	c.addLFUNode(node)
}

// Delete removes the entry for key and returns true if the key existed.
// minFreq is advanced if the deleted node's bucket was the minimum and is now empty.
func (c *LFUCache) Delete(key string) bool {
	if node, exists := c.cache[key]; exists {
		c.removeLFUNode(node)
		if _, exists = c.freqGroups[node.freq]; !exists && node.freq == c.minFreq {
			c.minFreq++
			return true
		}
	}
	return false
}

// Clear removes all entries, resets frequency tracking, and resets metrics.
func (c *LFUCache) Clear() {
	c.size = 0
	c.cache = make(map[string]*LFUNode, c.capacity)
	c.freqGroups = make(map[int]*FreqGroup)
	c.minFreq = 0
	c.metrics = Metrics{}
}

// Size returns the number of entries currently in the cache.
func (c *LFUCache) Size() int {
	return c.size
}

// Capacity returns the maximum number of entries the cache can hold.
func (c *LFUCache) Capacity() int {
	return c.capacity
}

// HitRate returns the fraction of Get calls that found a key (between 0.0 and 1.0).
// Returns 0.0 when no Get calls have been made yet.
func (c *LFUCache) HitRate() float64 {
	total := c.metrics.hits + c.metrics.misses
	if total > 0 {
		return float64(c.metrics.hits) / float64(total)
	}
	return 0.0
}

// ─── FIFO Cache ──────────────────────────────────────────────────────────────

// FIFOCache implements Cache with the First In First Out eviction policy.
// Entries are evicted in insertion order regardless of how often they are accessed.
// Internally uses a doubly-linked list: new entries are prepended to the head and
// eviction removes from the tail (the oldest entry).
type FIFOCache struct {
	capacity int
	size     int
	cache    map[string]*FIFONode
	metrics  Metrics
	head     *FIFONode // most recently inserted
	tail     *FIFONode // oldest inserted — next eviction victim
}

// FIFONode is a doubly-linked list node that stores a single cache entry.
type FIFONode struct {
	key   string
	value interface{}
	next  *FIFONode
	prev  *FIFONode
}

// NewFIFOCache returns a new FIFOCache with the given capacity.
// Returns nil if capacity is less than or equal to 0.
func NewFIFOCache(capacity int) *FIFOCache {
	if capacity <= 0 {
		return nil
	}
	return &FIFOCache{
		capacity: capacity,
		cache:    make(map[string]*FIFONode, capacity),
		metrics:  Metrics{},
	}
}

func newFIFONode(key string, value interface{}) *FIFONode {
	return &FIFONode{
		key:   key,
		value: value,
	}
}

// evictFIFONode removes the oldest entry (the tail of the list).
func (c *FIFOCache) evictFIFONode() bool {
	if c.size == 0 {
		return false
	}
	node := c.tail
	return c.removeFIFONode(node)
}

// removeFIFONode unlinks node from the doubly-linked list and deletes its map entry.
func (c *FIFOCache) removeFIFONode(node *FIFONode) bool {
	if node == nil || c.size == 0 {
		return false
	}
	if c.size == 1 {
		// Last node — reset both sentinels.
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
	// Middle node: relink its neighbors directly.
	node.prev.next = node.next
	node.next.prev = node.prev
	node.next = nil
	node.prev = nil
	delete(c.cache, node.key)
	c.size--
	return true
}

// Get returns the value for key.
// In FIFO semantics, reads do not affect eviction order.
func (c *FIFOCache) Get(key string) (interface{}, bool) {
	if node, exists := c.cache[key]; exists {
		c.metrics.hits++
		return node.value, true
	}
	c.metrics.misses++
	return nil, false
}

// Put inserts or updates a key-value pair.
// Updating an existing key changes its value but preserves its position in the
// eviction queue. If the cache is full, the oldest entry is evicted first.
func (c *FIFOCache) Put(key string, value interface{}) {
	if node, exists := c.cache[key]; exists {
		// Update value only — eviction order is unchanged for existing keys.
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
	// Prepend to head; the new entry is the most recently inserted.
	node.next = c.head
	c.head = node
	node.next.prev = node
	c.cache[key] = node
	c.size++
}

// Delete removes the entry for key and returns true if the key existed.
func (c *FIFOCache) Delete(key string) bool {
	if node, exists := c.cache[key]; exists {
		return c.removeFIFONode(node)
	}
	return false
}

// Clear removes all entries and resets metrics.
func (c *FIFOCache) Clear() {
	c.cache = make(map[string]*FIFONode)
	c.head = nil
	c.tail = nil
	c.size = 0
	c.metrics = Metrics{}
}

// Size returns the number of entries currently in the cache.
func (c *FIFOCache) Size() int {
	return c.size
}

// Capacity returns the maximum number of entries the cache can hold.
func (c *FIFOCache) Capacity() int {
	return c.capacity
}

// HitRate returns the fraction of Get calls that found a key (between 0.0 and 1.0).
// Returns 0.0 when no Get calls have been made yet.
func (c *FIFOCache) HitRate() float64 {
	total := c.metrics.hits + c.metrics.misses
	if total > 0 {
		return float64(c.metrics.hits) / float64(total)
	}
	return 0.0
}

// ─── Thread-Safe Wrapper ─────────────────────────────────────────────────────

// ThreadSafeCache wraps any Cache implementation with a mutex so that all
// operations are safe for concurrent use by multiple goroutines.
// A full write lock is used for every method because LRU and LFU mutate internal
// state (node repositioning) even during reads.
type ThreadSafeCache struct {
	cache Cache
	mu    sync.Mutex
}

// NewThreadSafeCache wraps cache with a mutex. Returns nil if cache is nil.
func NewThreadSafeCache(cache Cache) *ThreadSafeCache {
	if isNilCache(cache) {
		return nil
	}
	return &ThreadSafeCache{
		cache: cache,
	}
}

func isNilCache(cache Cache) bool {
	if cache == nil {
		return true
	}
	value := reflect.ValueOf(cache)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
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
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.cache.Size()
}

func (c *ThreadSafeCache) Capacity() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.cache.Capacity()
}

func (c *ThreadSafeCache) HitRate() float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.cache.HitRate()
}

// ─── Cache Factory ────────────────────────────────────────────────────────────

// NewCache creates and returns a Cache using the given eviction policy and capacity.
// Unrecognised policies fall back to LRU.
func NewCache(policy CachePolicy, capacity int) Cache {
	if capacity <= 0 {
		return nil
	}
	switch policy {
	case LRU:
		return NewLRUCache(capacity)
	case LFU:
		return NewLFUCache(capacity)
	case FIFO:
		return NewFIFOCache(capacity)
	default:
		return NewLRUCache(capacity)
	}
}

// NewThreadSafeCacheWithPolicy creates a thread-safe Cache using the given
// eviction policy and capacity. Unrecognised policies fall back to LRU.
func NewThreadSafeCacheWithPolicy(policy CachePolicy, capacity int) Cache {
	if capacity <= 0 {
		return nil
	}
	switch policy {
	case LRU:
		return NewThreadSafeCache(NewLRUCache(capacity))
	case LFU:
		return NewThreadSafeCache(NewLFUCache(capacity))
	case FIFO:
		return NewThreadSafeCache(NewFIFOCache(capacity))
	default:
		return NewThreadSafeCache(NewLRUCache(capacity))
	}
}
