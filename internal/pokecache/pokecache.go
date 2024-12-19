package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createAt time.Time
	val      []byte
}

type Cache struct {
	data map[string]cacheEntry
	mu   sync.RWMutex
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		data: make(map[string]cacheEntry),
	}

	c.ReadLoop(interval)

	return c
}

func (cache *Cache) Add(key string, val []byte) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	cache.data[key] = cacheEntry{
		createAt: time.Now(),
		val:      val,
	}
}

func (cache *Cache) Get(key string) ([]byte, bool) {
	cache.mu.RLock()
	defer cache.mu.RUnlock()
	entry, ok := cache.data[key]
	if !ok {
		return nil, false
	}

	return entry.val, true
}

func (cache *Cache) ReadLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		for range ticker.C {
			for key, entry := range cache.data {
				cache.mu.Lock()
				if time.Since(entry.createAt) > interval {
					delete(cache.data, key)
				}
				cache.mu.Unlock()
			}
		}
	}()
}
