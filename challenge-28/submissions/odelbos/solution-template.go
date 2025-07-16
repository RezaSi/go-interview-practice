package cache

import (
	"sync"
	"container/list"
	"slices"
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

type lruItem struct {
	key   string
	value any
}

type LRUCache struct {
	capacity int
	cache    map[string]*list.Element
	list     *list.List
	hits     int
	misses   int
	mu       sync.RWMutex
}

// NewLRUCache creates a new LRU cache with the specified capacity
func NewLRUCache(capacity int) *LRUCache {
	if capacity < 1 {
		return nil
	}
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*list.Element),
		list:  list.New(),
	}
}

func (c *LRUCache) Get(key string) (interface{}, bool) {
	c.mu.Lock()                     // Fix DATA RACE error within the tests
	defer c.mu.Unlock()

	if item, ok := c.cache[key]; ok {
		c.list.MoveToFront(item)
		c.hits++
		return item.Value.(*lruItem).value, true
	}
	c.misses++
	return nil, false
}

func (c *LRUCache) Put(key string, value interface{}) {
	if item, ok := c.cache[key]; ok {
		c.list.MoveToFront(item)
		item.Value.(*lruItem).value = value
		return
	}

	if len(c.cache) >= c.capacity {
		back := c.list.Back()
		if back != nil {
			backItem := back.Value.(*lruItem)
			delete(c.cache, backItem.key)
			c.list.Remove(back)
		}
	}

	item := c.list.PushFront(&lruItem{key, value})
	c.cache[key] = item
}

func (c *LRUCache) Delete(key string) bool {
	if item, ok := c.cache[key]; ok {
		c.list.Remove(item)
		delete(c.cache, key)
		return true
	}
	return false
}

func (c *LRUCache) Clear() {
	c.cache = make(map[string]*list.Element)
	c.list.Init()
	c.hits = 0
	c.misses = 0
}

func (c *LRUCache) Size() int {
	return len(c.cache)
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

//
// LFU Cache Implementation
//

type lfuItem struct {
	key      string
	value    any
	freq     int
	node     *list.Element
}

type LFUCache struct {
	capacity int
	cache    map[string]*lfuItem
	freqs    map[int]*list.List
	minFreq  int
	hits     int
	misses   int
}

// NewLFUCache creates a new LFU cache with the specified capacity
func NewLFUCache(capacity int) *LFUCache {
	return &LFUCache{
		capacity: capacity,
		cache:    make(map[string]*lfuItem),
		freqs:    make(map[int]*list.List),
	}
}

func (c *LFUCache) Get(key string) (interface{}, bool) {
	if item, ok := c.cache[key]; ok {
		c.hits++
		c.increment(item)
		return item.value, true
	}
	c.misses++
	return nil, false
}

func (c *LFUCache) Put(key string, value interface{}) {
	if c.capacity == 0 {
		return
	}

	if item, ok := c.cache[key]; ok {
		item.value = value
		c.increment(item)
		return
	}

	if len(c.cache) >= c.capacity {
		c.evict()
	}

	item := &lfuItem{key: key, value: value, freq:  1}
	if c.freqs[1] == nil {
		c.freqs[1] = list.New()
	}
	item.node = c.freqs[1].PushBack(item)
	c.cache[key] = item
	c.minFreq = 1
}

func (c *LFUCache) Delete(key string) bool {
	item, ok := c.cache[key]
	if !ok {
		return false
	}
	c.remove(item)
	return true
}

func (c *LFUCache) Clear() {
	c.cache = make(map[string]*lfuItem)
	c.freqs = make(map[int]*list.List)
	c.minFreq = 0
	c.hits = 0
	c.misses = 0
}

func (c *LFUCache) Size() int {
	return len(c.cache)
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

func (c *LFUCache) increment(item *lfuItem) {
	freq := item.freq
	c.freqs[freq].Remove(item.node)
	if c.freqs[freq].Len() == 0 {
		delete(c.freqs, freq)
		if c.minFreq == freq {
			c.minFreq++
		}
	}

	item.freq++
	if c.freqs[item.freq] == nil {
		c.freqs[item.freq] = list.New()
	}
	item.node = c.freqs[item.freq].PushBack(item)
}

func (c *LFUCache) evict() {
	lfuList := c.freqs[c.minFreq]
	if lfuList == nil {
		return
	}
	front := lfuList.Front()
	if front == nil {
		return
	}
	item := front.Value.(*lfuItem)
	c.remove(item)
}

func (c *LFUCache) remove(entry *lfuItem) {
	freq := entry.freq
	c.freqs[freq].Remove(entry.node)
	if c.freqs[freq].Len() == 0 {
		delete(c.freqs, freq)
		if c.minFreq == freq {
			c.minFreq++
		}
	}
	delete(c.cache, entry.key)
}

//
// FIFO Cache Implementation
//

type fifoItem struct {
    key   string
    value any
}

type FIFOCache struct {
    capacity int
    queue    []fifoItem
    items    map[string]any
    hits     int
    misses   int
}

// NewFIFOCache creates a new FIFO cache with the specified capacity
func NewFIFOCache(capacity int) *FIFOCache {
    return &FIFOCache{
        capacity: capacity,
        queue:    make([]fifoItem, 0, capacity),
        items:    make(map[string]any),
    }
}

func (c *FIFOCache) Get(key string) (interface{}, bool) {
    val, ok := c.items[key]
    if ok {
        c.hits++
        return val, true
    }
    c.misses++
    return nil, false
}

func (c *FIFOCache) Put(key string, value interface{}) {
    if _, ok := c.items[key]; ok {
        c.items[key] = value
        return
    }
    if len(c.queue) >= c.capacity {
        old := c.queue[0]
        c.queue = c.queue[1:]
        delete(c.items, old.key)
    }
    c.queue = append(c.queue, fifoItem{key, value})
    c.items[key] = value
}

func (c *FIFOCache) Delete(key string) bool {
    if _, ok := c.items[key]; ! ok {
        return false
    }
    delete(c.items, key)
    for i, item := range c.queue {
        if item.key == key {
            c.queue = slices.Delete(c.queue, i, i + 1)
            break
        }
    }
    return true
}

func (c *FIFOCache) Clear() {
    c.queue = make([]fifoItem, 0, c.capacity)
    c.items = make(map[string]any)
    c.hits = 0
    c.misses = 0
}

func (c *FIFOCache) Size() int {
	return len(c.items)
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

//
// Thread-Safe Cache Wrapper
//

type ThreadSafeCache struct {
	cache Cache
	mu    sync.RWMutex
}

// NewThreadSafeCache wraps any cache implementation to make it thread-safe
func NewThreadSafeCache(cache Cache) *ThreadSafeCache {
	if cache == nil {
		return nil
	}
	return &ThreadSafeCache{cache: cache}
}

func (c *ThreadSafeCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
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
