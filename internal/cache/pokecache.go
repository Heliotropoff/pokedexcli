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
	return &Cache{
		CachedData: make(map[string]cacheEntry),
	}
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
	data, ok := c.CachedData[key]
	if !ok {
		return []byte{}, ok
	}
	return data.val, ok
}
