package pokecache

import (
	"sync"
	"time"
)

var Cache *SafeMap

type SafeMap struct {
	mu   sync.Mutex
	data map[string]cacheEntry
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		data: make(map[string]cacheEntry),
	}
}

func (sm *SafeMap) Add(key string, val []byte) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.data[key] = cacheEntry{createdAt: time.Now(), val: val}
}

func (sm *SafeMap) Get(key string) ([]byte, bool) {
	var byteArr []byte
	sm.mu.Lock()
	defer sm.mu.Unlock()
	value, exists := sm.data[key]
	if exists {
		return value.val, true
	}
	return byteArr, false
}

func (sm *SafeMap) Delete(key string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.data, key)
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) {
	Cache = NewSafeMap()
	ReapLoop(Cache, interval)
}

func ReapLoop(Cache *SafeMap, intervalMinute time.Duration) {
	ticker := time.NewTicker(intervalMinute * time.Minute)
	defer ticker.Stop()

	currentTime := time.Now()
	clearCacheTime := currentTime.Add(-intervalMinute * time.Minute)
	for range ticker.C {
		for key, value := range Cache.data {
			if value.createdAt.Before(clearCacheTime) {
				Cache.Delete(key)
			}
		}
	}
}
