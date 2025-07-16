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

func NewCache(interval time.Duration) *SafeMap {
	Cache = NewSafeMap()
	go Cache.ReapLoop(interval)
	return Cache
}

func (Cache *SafeMap) ReapLoop(intervalDur time.Duration) {
	//fmt.Println("Inside ReapLoop logic")
	ticker := time.NewTicker(intervalDur)
	//fmt.Printf("Initialized ticker with value: %v\n", ticker.C)
	defer ticker.Stop()

	for range ticker.C {
		currentTime := time.Now()
		clearCacheTime := currentTime.Add(-intervalDur)
		//fmt.Printf("Looping over ticker %v\n", ticker.C)
		for key, value := range Cache.data {
			//fmt.Printf("Checking if value.createdAt: %v is before clearCacheTime: %v\n", value.createdAt, clearCacheTime)
			if value.createdAt.Before(clearCacheTime) {
				//fmt.Printf("Clearing cache for url: %s\n", key)
				Cache.Delete(key)
			}
		}
	}
}
