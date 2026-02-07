package pokecache

import (
	"time"
	"sync"
)

type cacheEntry struct {
	createdAt time.Time
	value []byte
}

type Cache struct {
	entries map[string]cacheEntry
	mutex *sync.Mutex
}

func NewCache(interval time.Duration) Cache {
	cache := Cache {
		entries: make(map[string]cacheEntry),
		mutex: &sync.Mutex{},
	}

	go cache.reapLoop(interval)

	return cache
}

func (cache *Cache) Add(key string, value []byte) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cache.entries[key] = cacheEntry {
		createdAt: time.Now().UTC(), 
		value: value,
	}
}

func (cache *Cache) Get(key string) ([]byte, bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	entry, exists := cache.entries[key]
	return entry.value, exists
}

func (cache *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		cache.reap(time.Now().UTC(), interval)
	}
}

func (cache *Cache) reap(now time.Time, interval time.Duration) {
		cache.mutex.Lock()
		defer cache.mutex.Unlock()

		for key, entry := range cache.entries {
			if entry.createdAt.Before(now.Add(-interval)) {
				delete(cache.entries, key)
			}
		}
}
