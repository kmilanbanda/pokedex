package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt 	time.Time
	val		[]byte
}

type Cache struct {
	entryMap	map[string]cacheEntry
	mu		sync.Mutex
	interval	time.Duration
}

func NewCache(interval time.Duration) *Cache {
	var cache Cache
	cache.entryMap = make(map[string]cacheEntry)
	cache.interval = interval
	go cache.reapLoop()
	return &cache
}

func (c *Cache)Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entryMap[key] = cacheEntry{
		createdAt:	time.Now(),
		val:		val, 
	}
}

func (c *Cache)Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.entryMap[key]
	if !ok {
		return nil, false
	}
	return entry.val, true
}

func (c *Cache)reapLoop() {
	ticker := time.NewTicker(c.interval)

	for range ticker.C {
		c.mu.Lock()
		for key, value := range c.entryMap {
			if time.Since(value.createdAt) > c.interval {
				delete(c.entryMap, key)
			}
		}	
		c.mu.Unlock()
	}

	
}
