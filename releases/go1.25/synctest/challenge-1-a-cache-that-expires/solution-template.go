package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type entry struct {
	value   string
	expires time.Time
}

// Cache is a string cache where every entry expires after a fixed TTL
type Cache struct {
	ttl   time.Duration
	mu    sync.Mutex
	items map[string]entry
}

// New returns a cache whose entries live for ttl
func New(ttl time.Duration) *Cache {
	return &Cache{ttl: ttl, items: make(map[string]entry)}
}

// Set stores a value and stamps it with an expiry time
func (c *Cache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = entry{value: value, expires: time.Now().Add(c.ttl)}
}

// Len reports how many entries are currently stored, expired or not
func (c *Cache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.items)
}

// Get returns the value for key. An entry that has passed its expiry time
// counts as a miss, even though it is still in the map.
func (c *Cache) Get(key string) (string, bool) {
	// TODO: return a miss when the entry is missing or has expired
	return "", false
}

// Cleanup deletes expired entries every interval until ctx is done.
// It is meant to run in its own goroutine.
func (c *Cache) Cleanup(ctx context.Context, interval time.Duration) {
	// TODO: tick every interval, delete expired entries, and return when ctx is done
}

func main() {
	c := New(time.Minute)
	c.Set("k", "v")

	v, ok := c.Get("k")
	fmt.Printf("fresh: %q %v\n", v, ok)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go c.Cleanup(ctx, 30*time.Second)

	fmt.Println("entries:", c.Len())
}
