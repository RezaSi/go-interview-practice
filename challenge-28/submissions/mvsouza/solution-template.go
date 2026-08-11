package cache

import (
	"container/list"
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

type CacheControl struct {
	hits     int
	misses   int
	capacity int
}

func (c *CacheControl) Capacity() int {
	return c.capacity
}

func (c *CacheControl) HitRate() float64 {
	if c.hits+c.misses == 0 {
		return 0.0
	}
	return float64(c.hits) / (float64(c.hits) + float64(c.misses))
}

//
// LRU Cache Implementation
//

type LRUItem struct {
	key   string
	value any
}

type LRUCache struct {
	values     map[string]*list.Element
	itemsOrder list.List
	CacheControl
}

// NewLRUCache creates a new LRU cache with the specified capacity
func NewLRUCache(capacity int) *LRUCache {
	if capacity <= 0 {
		return nil
	}
	return &LRUCache{
		values:     make(map[string]*list.Element),
		itemsOrder: list.List{},
		CacheControl: CacheControl{
			hits:     0,
			misses:   0,
			capacity: capacity,
		},
	}
}

func (c *LRUCache) Get(key string) (interface{}, bool) {
	if element, ok := c.values[key]; ok {
		c.itemsOrder.MoveToFront(element)
		c.hits++
		item := element.Value.(LRUItem)
		return item.value, true
	}
	c.misses++
	return nil, false
}

func (c *LRUCache) Put(key string, value interface{}) {
	if element, ok := c.values[key]; ok {
		element.Value = LRUItem{key, value}
		c.itemsOrder.MoveToFront(element)
		return
	}
	if c.Size() == c.capacity {
		el := c.itemsOrder.Back()
		if el != nil {
			item := el.Value.(LRUItem)
			c.Delete(item.key)
		}
	}
	c.values[key] = c.itemsOrder.PushFront(LRUItem{key, value})
}

func (c *LRUCache) Delete(key string) bool {
	if element, ok := c.values[key]; ok {
		c.itemsOrder.Remove(element)
		delete(c.values, key)
		return true
	}
	return false
}

func (c *LRUCache) Clear() {
	c.hits = 0
	c.misses = 0
	c.itemsOrder.Init()
	c.values = make(map[string]*list.Element)
}

func (c *LRUCache) Size() int {
	return len(c.values)
}

//
// LFU Cache Implementation
//

type LFUNode struct {
	key   string
	value any
	freq  int
}

type LFUCache struct {
	values    map[string]*list.Element
	freqLists map[int]*list.List
	minFreq   int
	CacheControl
}

// NewLFUCache creates a new LFU cache with the specified capacity
func NewLFUCache(capacity int) *LFUCache {
	if capacity <= 0 {
		return nil
	}
	return &LFUCache{
		values:    make(map[string]*list.Element),
		freqLists: make(map[int]*list.List),
		minFreq:   0,
		CacheControl: CacheControl{
			hits:     0,
			misses:   0,
			capacity: capacity,
		},
	}
}

func (c *LFUCache) updateFreq(element *list.Element) {
	node := element.Value.(*LFUNode)
	c.freqLists[node.freq].Remove(element)
	if c.freqLists[node.freq].Len() == 0 {
		delete(c.freqLists, node.freq)
		if node.freq == c.minFreq {
			c.minFreq++
		}
	}

	node.freq += 1
	if _, exists := c.freqLists[node.freq]; !exists {
		c.freqLists[node.freq] = list.New()
	}
	newEl := c.freqLists[node.freq].PushFront(node)
	c.values[node.key] = newEl
}

func (c *LFUCache) Get(key string) (interface{}, bool) {
	if element, ok := c.values[key]; ok {
		c.updateFreq(element)
		c.hits++
		node := element.Value.(*LFUNode)
		return node.value, true
	}
	c.misses++
	return nil, false
}

func (c *LFUCache) Put(key string, value interface{}) {
	if element, ok := c.values[key]; ok {
		node := element.Value.(*LFUNode)
		node.value = value
		c.updateFreq(element)
		return
	}

	if c.Size() == c.capacity {
		minList := c.freqLists[c.minFreq]
		if minList != nil {
			el := minList.Back()
			if el != nil {
				victim := el.Value.(*LFUNode)
				minList.Remove(el)
				if minList.Len() == 0 {
					delete(c.freqLists, c.minFreq)
				}
				delete(c.values, victim.key)
			}
		}
	}

	node := &LFUNode{
		key:   key,
		value: value,
		freq:  1,
	}
	if _, exists := c.freqLists[1]; !exists {
		c.freqLists[1] = list.New()
	}
	el := c.freqLists[1].PushFront(node)
	c.values[key] = el
	c.minFreq = 1
}

func (c *LFUCache) Delete(key string) bool {
	if element, ok := c.values[key]; ok {
		node := element.Value.(*LFUNode)
		list := c.freqLists[node.freq]
		list.Remove(element)
		if list.Len() == 0 {
			delete(c.freqLists, node.freq)
		}
		delete(c.values, key)

		if len(c.values) == 0 {
			c.minFreq = 0
		} else if node.freq == c.minFreq && list.Len() == 0 {
			min := -1
			for f := range c.freqLists {
				if min == -1 || f < min {
					min = f
				}
			}
			c.minFreq = min
		}
		return true
	}
	return false
}

func (c *LFUCache) Clear() {
	c.values = make(map[string]*list.Element)
	c.freqLists = make(map[int]*list.List)
	c.minFreq = 0
	c.hits = 0
	c.misses = 0
}

func (c *LFUCache) Size() int {
	return len(c.values)
}

//
// FIFO Cache Implementation
//

type fifoElement struct {
	value   any
	element *list.Element
}

type FIFOCache struct {
	values     map[string]fifoElement
	itemsOrder list.List
	CacheControl
}

// NewFIFOCache creates a new FIFO cache with the specified capacity
func NewFIFOCache(capacity int) *FIFOCache {
	if capacity <= 0 {
		return nil
	}
	return &FIFOCache{
		values:     make(map[string]fifoElement),
		itemsOrder: list.List{},
		CacheControl: CacheControl{
			hits:     0,
			misses:   0,
			capacity: capacity,
		},
	}
}

func (c *FIFOCache) Get(key string) (interface{}, bool) {
	if elem, ok := c.values[key]; ok {
		c.hits++
		return elem.value, true
	}
	c.misses++
	return nil, false
}

func (c *FIFOCache) Put(key string, value interface{}) {
	if elem, ok := c.values[key]; ok {
		elem.value = value
		c.values[key] = elem
		return
	}

	if c.Size() == c.capacity {
		el := c.itemsOrder.Front()
		if el != nil {
			oldestKey := el.Value.(string)
			c.itemsOrder.Remove(el)
			delete(c.values, oldestKey)
		}
	}

	// Insert new key at the back
	el := c.itemsOrder.PushBack(key)
	c.values[key] = fifoElement{value: value, element: el}
}

func (c *FIFOCache) Delete(key string) bool {
	if elem, ok := c.values[key]; ok {
		c.itemsOrder.Remove(elem.element)
		delete(c.values, key)
		return true
	}
	return false
}

func (c *FIFOCache) Clear() {
	c.hits = 0
	c.misses = 0
	c.itemsOrder.Init()
	c.values = make(map[string]fifoElement)
}

func (c *FIFOCache) Size() int {
	return len(c.values)
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
