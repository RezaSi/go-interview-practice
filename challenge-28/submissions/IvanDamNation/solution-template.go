package cache

import (
	"sync"
	"sync/atomic"
)

// Cache interface defines the contract for all cache implementations
type Cache interface {
	Get(key string) (value any, found bool)
	Put(key string, value any)
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

type LRUCache struct {
	capacity    int
	cache       map[string]*lruNode
	mostRecent  *lruNode
	leastRecent *lruNode
	hits        atomic.Uint64
	misses      atomic.Uint64
}

type lruNode struct {
    key     string
    value   any
    next    *lruNode
    prev    *lruNode
}

// NewLRUCache creates a new LRU cache with the specified capacity
func NewLRUCache(capacity int) *LRUCache {
	if capacity <= 0 {
	    return nil
	}
	
	dummy1, dummy2 := &lruNode{}, &lruNode{} 
	dummy1.next = dummy2
	dummy2.prev = dummy1
	cache := make(map[string]*lruNode, capacity)
	
	return &LRUCache{
	    capacity: capacity,
	    cache: cache,
	    mostRecent: dummy2,
	    leastRecent: dummy1,
	}
}

func (c *LRUCache) Get(key string) (any, bool) {
	node, exists := c.cache[key]
	if !exists {
	    c.misses.Add(1)
	    return nil, false
	}
	
	c.hits.Add(1)
	c.remove(node)
	c.moveToMostRecent(node)
	
	return node.value, true
}

func (c *LRUCache) Put(key string, value any) {
	var node *lruNode
	node, exists := c.cache[key]
	if exists {
	    node.value = value
	    c.remove(node)
	} else {
	    node = &lruNode{
    	    key: key,
    	    value: value,
    	}
    	c.cache[key] = node
	}
	
	c.moveToMostRecent(node)
	if len(c.cache) > c.capacity {
	    c.evict()
	}
}

func (c *LRUCache) Delete(key string) bool {
	node, exists := c.cache[key]
	if !exists {
	    return false
	}
	
	c.remove(node)
	delete(c.cache, key)
	
	return true
}

func (c *LRUCache) Clear() {
	clear(c.cache)
	
	c.hits.Store(0)
	c.misses.Store(0)
	
	c.mostRecent.prev = c.leastRecent
	c.leastRecent.next = c.mostRecent
}

func (c *LRUCache) Size() int {
	return len(c.cache)
}

func (c *LRUCache) Capacity() int {
	return c.capacity
}

func (c *LRUCache) HitRate() float64 {
	h := c.hits.Load()
	m := c.misses.Load()
	
	if h + m == 0 {
	    return 0.0
	}
	
	return float64(h) / float64(h + m)
}

func (c *LRUCache) remove(node *lruNode) {
    node.next.prev = node.prev
    node.prev.next = node.next
    node.prev, node.next = nil, nil
}

func (c *LRUCache) moveToMostRecent(node *lruNode) {
    node.next = c.mostRecent
    node.prev = c.mostRecent.prev
    node.prev.next = node
    c.mostRecent.prev = node
}

func (c *LRUCache) evict() {
    node := c.leastRecent.next
    delete(c.cache, node.key)
    c.remove(node)
}

//
// LFU Cache Implementation
//

type LFUCache struct {
	capacity    int
	minFreq     int
	hits        atomic.Uint64
	misses      atomic.Uint64
	
	cache       map[string]*lfuNode
	freqList    map[int]*freqBucket
	
	mostFreq    *freqBucket
	leastFreq   *freqBucket
}

type freqBucket struct {
    freq        int
    nElems      int
    
    mostRecent  *lfuNode
    leastRecent *lfuNode
    
    next        *freqBucket
    prev        *freqBucket
}

type lfuNode struct {
    key             string
    value           any
    
    freq            int
    parentBucket    *freqBucket
    
    next            *lfuNode
    prev            *lfuNode
}

// NewLFUCache creates a new LFU cache with the specified capacity
func NewLFUCache(capacity int) *LFUCache {
	if capacity <= 0 {
	    return nil
	}
	
	dummyBucket1, dummyBucket2 := initBucket(0), initBucket(0)
	dummyBucket1.next = dummyBucket2
	dummyBucket2.prev = dummyBucket1
	
	return &LFUCache{
	    capacity: capacity,
	    cache: make(map[string]*lfuNode, capacity),
	    freqList: make(map[int]*freqBucket),
	    mostFreq: dummyBucket2,
	    leastFreq: dummyBucket1,
	}
}

func initBucket(freq int) *freqBucket {
    dummy1, dummy2 := &lfuNode{}, &lfuNode{}
	dummy1.next = dummy2
	dummy2.prev = dummy1
	
	return &freqBucket{
	    freq: freq,
	    mostRecent: dummy2,
	    leastRecent: dummy1,
	}
}

func (c *LFUCache) Get(key string) (any, bool) {
	node, exists := c.cache[key]
	if !exists {
	    c.misses.Add(1)
	    return nil, false
	}
	
	c.hits.Add(1)
	c.incrementFreq(node)
	
	return node.value, true
}

func (c *LFUCache) Put(key string, value any) {
	node, exists := c.cache[key]
	if exists {
	    node.value = value
	    c.incrementFreq(node)
	    return
	}
	
	if len(c.cache) == c.capacity {
	    c.evict()
	}
	
	bucket, exists := c.freqList[1]
	if !exists {
	    bucket = initBucket(1)
	    c.freqList[1] = bucket
	    c.putBucketOnRight(bucket, c.leastFreq.next)
	}
	
	node = &lfuNode{
	    key: key,
	    value: value,
	    freq: 1,
	}
	
	c.minFreq = 1
	c.cache[key] = node
	c.putNodeInBucket(bucket, node)
}

func (c *LFUCache) Delete(key string) bool {
	node, exists := c.cache[key]
	if !exists {
	    return false
	}
	
	bucket := node.parentBucket
	c.removeNode(node)
	
	bucket.nElems--
	if bucket.nElems == 0 {
	    c.dropBucket(bucket)
	}
	
	delete(c.cache, key)
	
	return true
}

func (c *LFUCache) Clear() {
	clear(c.cache)
	clear(c.freqList)
	
	c.hits.Store(0)
	c.misses.Store(0)
	
	c.mostFreq.prev = c.leastFreq
	c.leastFreq.next = c.mostFreq
}

func (c *LFUCache) Size() int {
	return len(c.cache)
}

func (c *LFUCache) Capacity() int {
	return c.capacity
}

func (c *LFUCache) HitRate() float64 {
    h := c.hits.Load()
    m := c.misses.Load()
	
	if h + m == 0 {
	    return 0.0
	}
	
	return float64(h) / float64(h + m)
}

func (c *LFUCache) dropBucket(bucket *freqBucket) {
    bucket.next.prev = bucket.prev
    bucket.prev.next = bucket.next
    delete(c.freqList, bucket.freq)
}

func (c *LFUCache) putBucketOnRight(newBucket, bucketLeft *freqBucket) {
    newBucket.next = bucketLeft
    newBucket.prev = bucketLeft.prev
    bucketLeft.prev.next = newBucket
    bucketLeft.prev = newBucket
}

func (c *LFUCache) removeNode(node *lfuNode) {
    node.next.prev = node.prev
    node.prev.next = node.next
    node.next, node.prev, node.parentBucket = nil, nil, nil
}

func (c *LFUCache) putNodeInBucket(bucket *freqBucket, node *lfuNode) {
    bucket.nElems++
    node.parentBucket = bucket
    node.next = bucket.mostRecent
    node.prev = bucket.mostRecent.prev
    bucket.mostRecent.prev.next = node
    bucket.mostRecent.prev = node
}

func (c *LFUCache) incrementFreq(node *lfuNode) {
    bucket := node.parentBucket
    buckLeft := bucket.next
    
    c.removeNode(node)
    
    bucket.nElems--
    if bucket.nElems == 0 {
        c.dropBucket(bucket)
        
        if bucket.freq == c.minFreq {
            c.minFreq++
        }
    }
    
    newBucket, exists := c.freqList[node.freq + 1]
    if !exists {
        newBucket = initBucket(node.freq + 1)
        c.freqList[node.freq + 1] = newBucket
        c.putBucketOnRight(newBucket, buckLeft)
    }
    
    node.freq++
    c.putNodeInBucket(newBucket, node)
}

func (c *LFUCache) evict() {
    bucket := c.freqList[c.minFreq]
    node := bucket.leastRecent.next
    
    c.removeNode(node)
    
    bucket.nElems--
    if bucket.nElems == 0 {
        c.dropBucket(bucket)
    }
    
    delete(c.cache, node.key)
}

//
// FIFO Cache Implementation
//

type FIFOCache struct {
	capacity    int
	cache       map[string]*fifoNode
	mostRecent  *fifoNode
	leastRecent *fifoNode
	hits        atomic.Uint64
	misses      atomic.Uint64
}

type fifoNode struct {
    key     string
    value   any
    next    *fifoNode
    prev    *fifoNode
}

// NewFIFOCache creates a new FIFO cache with the specified capacity
func NewFIFOCache(capacity int) *FIFOCache {
	if capacity <= 0 {
	    return nil
	}
	
	dummy1, dummy2 := &fifoNode{}, &fifoNode{}
	dummy1.next = dummy2
	dummy2.prev = dummy1
	
	return &FIFOCache{
	    capacity: capacity,
	    cache: make(map[string]*fifoNode, capacity),
	    mostRecent: dummy2,
	    leastRecent: dummy1,
	}

}

func (c *FIFOCache) Get(key string) (any, bool) {
	node, exists := c.cache[key]
	if !exists {
	    c.misses.Add(1)
	    return nil, false
	}
	
	c.hits.Add(1)
	
	return node.value, true
}

func (c *FIFOCache) Put(key string, value any) {
	node, exists := c.cache[key]
	if exists {
	    node.value = value
	    return
	}
	
	if len(c.cache) == c.capacity {
	    node = c.leastRecent.next
	    c.remove(node)
	    delete(c.cache, node.key)
	}
	
	node = &fifoNode{
	    key: key,
	    value: value,
	}
	
	c.mostRecent.prev.next = node
	node.prev = c.mostRecent.prev
	c.mostRecent.prev = node
	node.next = c.mostRecent
	
	c.cache[node.key] = node
}

func (c *FIFOCache) Delete(key string) bool {
	node, exists := c.cache[key]
	if !exists {
	    return false
	}
	
	c.remove(node)
	delete(c.cache, node.key)
	
	return true
}

func (c *FIFOCache) Clear() {
	clear(c.cache)
	
	c.hits.Store(0)
	c.misses.Store(0)
	
	c.leastRecent.next = c.mostRecent
	c.mostRecent.prev = c.leastRecent
}

func (c *FIFOCache) Size() int {
	return len(c.cache)
}

func (c *FIFOCache) Capacity() int {
	return c.capacity
}

func (c *FIFOCache) HitRate() float64 {
    h := c.hits.Load()
    m := c.misses.Load()
    
	if h + m == 0 {
	    return 0.0
	}
	
	return float64(h) / float64(h + m)
}

func (c *FIFOCache) remove(node *fifoNode) {
    node.next.prev = node.prev
    node.prev.next = node.next
    node.next, node.prev = nil, nil
}

//
// Thread-Safe Cache Wrapper
//

type ThreadSafeCache struct {
	cache       Cache
	isMutating  bool
	
	mu          sync.RWMutex
}

// NewThreadSafeCache wraps any cache implementation to make it thread-safe
func NewThreadSafeCache(cache Cache) *ThreadSafeCache {
	if cache == nil {
	    return nil
	}
	
	newCache := &ThreadSafeCache{cache: cache}
	
	_, isLRU := cache.(*LRUCache)
	_, isLFU := cache.(*LFUCache)
	if isLRU || isLFU {
	    newCache.isMutating = true
	}
	
	return newCache
}

func (c *ThreadSafeCache) Get(key string) (any, bool) {
	if c.isMutating {
	    c.mu.Lock()
	    defer c.mu.Unlock()
	} else {
	    c.mu.RLock()
	    defer c.mu.RUnlock()
	}
	
	res, ok := c.cache.Get(key)
	
	return res, ok
}

func (c *ThreadSafeCache) Put(key string, value any) {
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
	return c.cache.HitRate()
}

//
// Cache Factory Functions
//

// NewCache creates a cache with the specified policy and capacity
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
		return nil
	}
}

// NewThreadSafeCacheWithPolicy creates a thread-safe cache with the specified policy
func NewThreadSafeCacheWithPolicy(policy CachePolicy, capacity int) Cache {
	newCache := NewCache(policy, capacity)
	if newCache == nil {
	    return nil
	}
	
	return NewThreadSafeCache(newCache)
}
