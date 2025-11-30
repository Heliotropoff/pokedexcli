package cache

import (
	"sync"
	"time"
)

type Cache struct {
	CachedData map[string]cacheEntry
	Mutext     sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	newCache := &Cache{
		CachedData: make(map[string]cacheEntry),
	}
	go newCache.reapLoop(interval)
	return newCache
}

func (c *Cache) Add(key string, val []byte) {
	c.Mutext.Lock()
	defer c.Mutext.Unlock()
	c.CachedData[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.Mutext.Lock()
	defer c.Mutext.Unlock()
	data, ok := c.CachedData[key]
	if !ok {
		return []byte{}, ok
	}
	return data.val, ok
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		for key, _ := range c.CachedData {
			timeLive := time.Since(c.CachedData[key].createdAt)
			if timeLive > interval {
				c.Mutext.Lock()
				delete(c.CachedData, key)
				c.Mutext.Unlock()
			}
		}
	}
}
