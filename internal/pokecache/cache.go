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

func NewCache(entryLifespan time.Duration) Cache {
	cache := Cache {
		entries: make(map[string]cacheEntry),
		mutex: &sync.Mutex{},
	}

	go cache.reapLoop(entryLifespan)

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

func (cache *Cache) reapLoop(entryLifespan time.Duration) {
	ticker := time.NewTicker(entryLifespan)
	defer ticker.Stop()

	for range ticker.C {
		cache.reap(time.Now().UTC(), entryLifespan)
	}
}

func (cache *Cache) reap(now time.Time, entryLifespan time.Duration) {
		cache.mutex.Lock()
		defer cache.mutex.Unlock()

		for key, entry := range cache.entries {
			if entry.createdAt.Before(now.Add(-entryLifespan)) {
				delete(cache.entries, key)
			}
		}
}
