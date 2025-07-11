# Hints for Challenge 28: Cache Implementation with Multiple Eviction Policies

## Hint 1: LRU Cache Design Pattern
LRU (Least Recently Used) requires O(1) access and O(1) eviction:

**Data Structure Choice:**
- **Hash Map**: O(1) access to cache nodes by key
- **Doubly Linked List**: O(1) insertion/deletion at any position
- **Dummy Head/Tail**: Simplifies edge cases in linked list operations

**Key Components:**
- Node structure with `prev`/`next` pointers
- Map for quick lookups: `map[string]*Node`
- Dummy nodes to avoid null checks

```go
type LRUNode struct {
    key, value interface{}
    prev, next *LRUNode
}

// Always insert after head, remove before tail
head.next = cache.tail  // Initialize empty list
tail.prev = cache.head
```

## Hint 2: LRU Operations Logic
Understanding the core LRU operations:

**Get Operation:**
- If key exists: move node to head (mark as recently used), return value
- If key doesn't exist: return cache miss

**Put Operation:**
- If key exists: update value, move to head
- If key doesn't exist: create new node, add to head
- If at capacity: remove tail node first (LRU), then add new

**Essential Helper Methods:**
```go
func moveToHead(node) {
    removeNode(node)  // Unlink from current position
    addToHead(node)   // Insert after dummy head
}

func addToHead(node) {
    node.next = head.next
    node.prev = head
    head.next.prev = node
    head.next = node
}

func removeTail() {
    lru := tail.prev  // Get LRU node
    removeNode(lru)   // Unlink it
    return lru
}
```

## Hint 3: LFU Cache Design
LFU (Least Frequently Used) requires tracking access frequency:

**Key Concept:** Group nodes by frequency in separate doubly-linked lists
- `freqMap[1]` → list of nodes accessed 1 time
- `freqMap[2]` → list of nodes accessed 2 times  
- Track `minFreq` to know which frequency list to evict from

**Data Structure:**
- Nodes have `freq` field to track access count
- `freqMap` maps frequency → dummy head of node list
- When frequency increases, move node to new frequency list

```go
type LFUNode struct {
    key, value interface{}
    freq       int
    prev, next *LFUNode
}

// Move node from freq N to freq N+1
func updateFreq(node) {
    removeFromFreqList(node, oldFreq)
    node.freq++
    addToFreqList(node, newFreq)
}
```
```

## Hint 4: LFU Cache - Get and Put Implementation
Handle frequency updates and eviction:
```go
func (c *LFUCache) Get(key string) (interface{}, bool) {
    if node, exists := c.cache[key]; exists {
        c.updateFreq(node)
        c.hits++
        return node.value, true
    }
    c.misses++
    return nil, false
}

func (c *LFUCache) updateFreq(node *LFUNode) {
    oldFreq := node.freq
    newFreq := oldFreq + 1
    
    c.removeFromFreqList(node)
    node.freq = newFreq
    c.addToFreqList(node, newFreq)
    
    // Update minFreq if necessary
    if oldFreq == c.minFreq && c.isFreqListEmpty(oldFreq) {
        c.minFreq++
    }
}

func (c *LFUCache) isFreqListEmpty(freq int) bool {
    head := c.freqMap[freq]
    return head != nil && head.next == head
}

func (c *LFUCache) Put(key string, value interface{}) {
    if c.capacity == 0 {
        return
    }
    
    if node, exists := c.cache[key]; exists {
        node.value = value
        c.updateFreq(node)
        return
    }
    
    if c.size >= c.capacity {
        c.evict()
    }
    
    newNode := &LFUNode{key: key, value: value, freq: 1}
    c.cache[key] = newNode
    c.addToFreqList(newNode, 1)
    c.minFreq = 1
    c.size++
}

func (c *LFUCache) evict() {
    head := c.getFreqList(c.minFreq)
    victim := head.prev
    c.removeFromFreqList(victim)
    delete(c.cache, victim.key)
    c.size--
}
```

## Hint 5: FIFO Cache - Simple Queue Structure
Implement FIFO using a slice as a queue:
```go
type FIFOCache struct {
    capacity int
    cache    map[string]interface{}
    order    []string
    hits     int64
    misses   int64
}

func NewFIFOCache(capacity int) *FIFOCache {
    return &FIFOCache{
        capacity: capacity,
        cache:    make(map[string]interface{}),
        order:    make([]string, 0, capacity),
    }
}

func (c *FIFOCache) Get(key string) (interface{}, bool) {
    if value, exists := c.cache[key]; exists {
        c.hits++
        return value, true
    }
    c.misses++
    return nil, false
}

func (c *FIFOCache) Put(key string, value interface{}) {
    if _, exists := c.cache[key]; exists {
        c.cache[key] = value
        return
    }
    
    if len(c.cache) >= c.capacity {
        // Remove oldest (first in)
        oldest := c.order[0]
        delete(c.cache, oldest)
        c.order = c.order[1:]
    }
    
    c.cache[key] = value
    c.order = append(c.order, key)
}

func (c *FIFOCache) Delete(key string) bool {
    if _, exists := c.cache[key]; exists {
        delete(c.cache, key)
        
        // Remove from order slice
        for i, k := range c.order {
            if k == key {
                c.order = append(c.order[:i], c.order[i+1:]...)
                break
            }
        }
        return true
    }
    return false
}
```

## Hint 6: Thread-Safe Wrapper
Add thread safety with read-write mutex:
```go
import "sync"

type ThreadSafeCache struct {
    cache Cache
    mutex sync.RWMutex
}

func NewThreadSafeCache(cache Cache) *ThreadSafeCache {
    return &ThreadSafeCache{
        cache: cache,
    }
}

func (c *ThreadSafeCache) Get(key string) (interface{}, bool) {
    c.mutex.RLock()
    defer c.mutex.RUnlock()
    return c.cache.Get(key)
}

func (c *ThreadSafeCache) Put(key string, value interface{}) {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    c.cache.Put(key, value)
}

func (c *ThreadSafeCache) Delete(key string) bool {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    return c.cache.Delete(key)
}

func (c *ThreadSafeCache) Clear() {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    c.cache.Clear()
}

func (c *ThreadSafeCache) Size() int {
    c.mutex.RLock()
    defer c.mutex.RUnlock()
    return c.cache.Size()
}

func (c *ThreadSafeCache) Capacity() int {
    c.mutex.RLock()
    defer c.mutex.RUnlock()
    return c.cache.Capacity()
}

func (c *ThreadSafeCache) HitRate() float64 {
    c.mutex.RLock()
    defer c.mutex.RUnlock()
    return c.cache.HitRate()
}
```

## Hint 7: Cache Factory Pattern
Implement factory pattern for cache creation:
```go
type CachePolicy int

const (
    LRU CachePolicy = iota
    LFU
    FIFO
)

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

func NewThreadSafeCacheWithPolicy(policy CachePolicy, capacity int) Cache {
    baseCache := NewCache(policy, capacity)
    return NewThreadSafeCache(baseCache)
}

// Utility function for testing
func BenchmarkCache(cache Cache, operations int) {
    start := time.Now()
    
    for i := 0; i < operations; i++ {
        key := fmt.Sprintf("key_%d", i%1000)
        
        if i%3 == 0 {
            cache.Put(key, i)
        } else {
            cache.Get(key)
        }
    }
    
    duration := time.Since(start)
    fmt.Printf("Completed %d operations in %v\n", operations, duration)
    fmt.Printf("Hit rate: %.2f%%\n", cache.HitRate()*100)
}
```

## Key Cache Implementation Concepts:
- **O(1) Operations**: Use hash maps for constant-time access
- **Doubly Linked List**: Efficient insertion/deletion for LRU
- **Frequency Buckets**: Group nodes by frequency for LFU
- **Thread Safety**: Use RWMutex for concurrent access
- **Memory Management**: Proper cleanup to prevent leaks
- **Factory Pattern**: Flexible cache creation and configuration 