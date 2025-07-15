package pokecache

import (
	"time"
)

type SafeMap struct {
	mu sync.Mutex
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
	sm.data[key] = {createdAt: time.Now(), val: val}
}

func (sm *SafeMap) Get(key string) ([]byte, bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	value, exists := sm.data[key]
	if exists {
		return value, true
	}
	return []byte, false
}

func (sm *SafeMap) Delete(key string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.data, key)
}

type Cache struct {
	map[string]cacheEntry
}

type cacheEntry struct {
	createdAt time.Time
	val []byte
}

func NewCache(interval time.Duration) {
	Cache = NewSafeMap()
	Cache.ReapLoop()
}

func ReapLoop(intervalMinute int) {
	ticker := time.NewTicker(intervalMinute * time.Minute)
	defer ticker.Stop()

	currentTime := time.Now()
	clearCacheTime := currentTime.Add(-intervalMinute * time.Minute)
	for range ticker.C {
		for key, value := range Cache{
			if value.createdAt.Before(clearCacheTime) {
				Cache.Delete(key)
			}
		}
	}
}