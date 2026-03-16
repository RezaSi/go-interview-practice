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

type DLL struct {
    Left, Right *DLL
    Val interface{}
    Frequency int // for LFU
}

func (node *DLL) Delete() {
    nl := node.Left
    nr := node.Right
    if nl != nil {
        nl.Right = nr
    }
    if nr != nil {
        nr.Left = nl
    }
}

func (node *DLL) Swap(to *DLL) {
    if node.Left == to {
        node.Left = to.Left
        to.Right = node.Right
        to.Left = node
        node.Right = to
    } else {
        toLeft := to.Left
        toRight := to.Right
        
        to.Left = node.Left
        to.Right = node.Right
        
        node.Left = toLeft
        node.Right = toRight
    }

    if node.Left != nil {
        node.Left.Right = node
    }

    if to.Right != nil {
        to.Right.Left = to
    }
}

//
// LRU Cache Implementation
//

type LRUCache struct {
    capacity, size, hit, miss int
    top, bottom *DLL
	cache map[string]*DLL // key -> DLL node
	revMap map[*DLL]string // DLL node -> key (for delete)
}

// NewLRUCache creates a new LRU cache with the specified capacity
func NewLRUCache(capacity int) *LRUCache {
    // Note: negative capacity is normalized to 0
    capacity = max(capacity, 0)
	return &LRUCache{
	    capacity: capacity,
	    cache: make(map[string]*DLL, capacity),
	    revMap: make(map[*DLL]string, capacity),
	}
}

func (c *LRUCache) moveToTop(node *DLL) {
    if c.top == nil {
        c.top = node
        c.bottom = node
        return
    }
    if node == c.top {
		return
	}
    if node.Left != nil {
        node.Left.Right = node.Right
        
        if node == c.bottom {
            c.bottom = node.Left
        }
    }
    if node.Right != nil {
        node.Right.Left = node.Left
    }
    node.Left = nil
    node.Right = c.top
    c.top.Left = node
    c.top = node
}

func (c *LRUCache) tryToEvict() {
    for c.size > c.capacity && c.bottom != nil {
        bottomLeft := c.bottom.Left
        bottomKey := c.revMap[c.bottom]
        delete(c.revMap, c.bottom)
        delete(c.cache, bottomKey)
        c.bottom.Delete()
        c.bottom = bottomLeft
        c.size--
	}
	if c.bottom == nil {
	    c.top = nil
	}
}

func (c *LRUCache) Get(key string) (interface{}, bool) {
	node, exists := c.cache[key]
	if !exists {
	    c.miss++
	    return nil, false
	}
	c.moveToTop(node)
	c.hit++
	return node.Val, true
}

func (c *LRUCache) Put(key string, value interface{}) {
	node, exists := c.cache[key]
	if exists {
	    node.Val = value
	    c.moveToTop(node)
	    return
	}
	node = &DLL{Val: value}
	c.moveToTop(node)
	c.cache[key] = node
	c.revMap[node] = key
	c.size++
    c.tryToEvict()
}

func (c *LRUCache) Delete(key string) bool {
	node, exists := c.cache[key]
	if !exists {
	    return false
	}
	if node == c.bottom {
	    c.bottom = node.Left
	}
	if node == c.top {
	    c.top = node.Right
	}
	delete(c.revMap, node)
    delete(c.cache, key)
	node.Delete()
	c.size--
	return true
}

func (c *LRUCache) Clear() {
    c.cache = make(map[string]*DLL, c.capacity)
    c.revMap = make(map[*DLL]string, c.capacity)
    c.size = 0
    c.top = nil
    c.bottom = nil
    c.hit = 0
    c.miss = 0
}

func (c *LRUCache) Size() int {
	return c.size
}

func (c *LRUCache) Capacity() int {
	return c.capacity
}

func (c *LRUCache) HitRate() float64 {
    if c.hit + c.miss == 0 {
        return 0.0
    }
	return float64(c.hit) / float64(c.hit + c.miss)
}

//
// LFU Cache Implementation
//

type LFUCache struct {
	capacity, size, hit, miss int
    top, bottom *DLL
	cache map[string]*DLL // key -> DLL node
	revMap map[*DLL]string // DLL node -> key (for delete)
}

// NewLFUCache creates a new LFU cache with the specified capacity
func NewLFUCache(capacity int) *LFUCache {
	// Note: negative capacity is normalized to 0
    capacity = max(capacity, 0)
	return &LFUCache{
	    capacity: capacity,
	    cache: make(map[string]*DLL, capacity),
	    revMap: make(map[*DLL]string, capacity),
	}
}

func (c *LFUCache) tryToUpByFrequency(node *DLL) {
    if node == c.bottom && node.Left != nil {
        c.bottom = node.Left
    }
    for node.Left != nil && node.Frequency >= node.Left.Frequency {
        node.Swap(node.Left)
    }
}

func (c *LFUCache) tryToEvict() {
    for c.size >= c.capacity && c.bottom != nil {
        bottomLeft := c.bottom.Left
        bottomKey := c.revMap[c.bottom]
        delete(c.revMap, c.bottom)
        delete(c.cache, bottomKey)
        c.bottom.Delete()
        c.bottom = bottomLeft
        c.size--
	}
	if c.bottom == nil {
	    c.top = nil
	}
}

func (c *LFUCache) Get(key string) (interface{}, bool) {
	node, exists := c.cache[key]
	if !exists {
	    c.miss++
	    return nil, false
	}
	node.Frequency++
	c.tryToUpByFrequency(node)
	c.hit++
	return node.Val, true
}

func (c *LFUCache) Put(key string, value interface{}) {
    node, exists := c.cache[key]
    if exists {
        node.Val = value
        return
    }
	c.tryToEvict()
    node = &DLL{
        Val: value,
        Frequency: 1,
    }
    if c.bottom == nil {
        c.top = node
        c.bottom = node
    } else {
        node.Left = c.bottom
        c.bottom.Right = node
        c.bottom = node
    }
	c.cache[key] = node
	c.revMap[node] = key
	c.size++
	c.tryToUpByFrequency(node)
}

func (c *LFUCache) Delete(key string) bool {
	node, exists := c.cache[key]
	if !exists {
	    return false
	}
	if node == c.bottom {
	    c.bottom = node.Left
	}
	if node == c.top {
	    c.top = node.Right
	}
	delete(c.revMap, node)
    delete(c.cache, key)
	node.Delete()
	c.size--
	return true
}

func (c *LFUCache) Clear() {
    c.cache = make(map[string]*DLL, c.capacity)
    c.revMap = make(map[*DLL]string, c.capacity)
    c.size = 0
    c.top = nil
    c.bottom = nil
    c.hit = 0
    c.miss = 0
}

func (c *LFUCache) Size() int {
	return c.size
}

func (c *LFUCache) Capacity() int {
	return c.capacity
}

func (c *LFUCache) HitRate() float64 {
	if c.hit + c.miss == 0 {
        return 0.0
    }
	return float64(c.hit) / float64(c.hit + c.miss)
}

//
// FIFO Cache Implementation
//

type FIFOCache struct {
	capacity, size, hit, miss int
	queue []string
	cache map[string]interface{}
}

// NewFIFOCache creates a new FIFO cache with the specified capacity
func NewFIFOCache(capacity int) *FIFOCache {
	// Note: negative capacity is normalized to 0
    capacity = max(capacity, 0)
	return &FIFOCache{
	    capacity: capacity,
	    queue: make([]string, 0),
	    cache: make(map[string]interface{}, capacity),
	}
}

func (c *FIFOCache) tryToEvict() {
    for (c.size >= c.capacity && len(c.queue) > 0) {
        topKey := c.queue[0]
        c.queue = c.queue[1:]
        if _, found := c.cache[topKey]; !found {
            continue
        }
        delete(c.cache, topKey)
        c.size--
    }
}

func (c *FIFOCache) Get(key string) (interface{}, bool) {
	val, exists := c.cache[key]
	if !exists {
	    c.miss++
	    return nil, false
	}
	c.hit++
	return val, true
}

func (c *FIFOCache) Put(key string, value interface{}) {
	_, exists := c.cache[key]
	if exists {
	    c.cache[key] = value
	    return
	}
	c.tryToEvict()
	c.cache[key] = value
	c.queue = append(c.queue, key)
	c.size++
}

func (c *FIFOCache) Delete(key string) bool {
	if _, found := c.cache[key]; !found {
	    return false
	}
	delete(c.cache, key)
	c.size--
	return true
}

func (c *FIFOCache) Clear() {
	c.size = 0
	c.hit = 0
	c.miss = 0
	c.queue = make([]string, 0)
	c.cache = make(map[string]interface{}, c.capacity)
}

func (c *FIFOCache) Size() int {
	return c.size
}

func (c *FIFOCache) Capacity() int {
	return c.capacity
}

func (c *FIFOCache) HitRate() float64 {
	if c.hit + c.miss == 0 {
        return 0.0
    }
	return float64(c.hit) / float64(c.hit + c.miss)
}

//
// Thread-Safe Cache Wrapper
//

type ThreadSafeCache struct {
	cache Cache
	mu    sync.RWMutex
}

// NewThreadSafeCache wraps any cache implementation to make it thread-safe
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

// NewThreadSafeCacheWithPolicy creates a thread-safe cache with the specified policy
func NewThreadSafeCacheWithPolicy(policy CachePolicy, capacity int) Cache {
	return NewThreadSafeCache(NewCache(policy, capacity))
}
