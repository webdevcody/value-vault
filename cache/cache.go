package cache

import (
	"sync"
)

var cache map[string][]byte = make(map[string][]byte)
var dirty map[string]bool = make(map[string]bool)
var cacheMutex = sync.RWMutex{}
var dirtyMutex = sync.RWMutex{}

func StoreIntoCache(key string, value []byte) {
	cacheMutex.Lock()
	cache[key] = value
	cacheMutex.Unlock()
}

func GetFromCache(key string) []byte {
	cacheMutex.RLock()
	value := cache[key]
	cacheMutex.RUnlock()
	return value
}

func IsDirty(key string) bool {
	dirtyMutex.RLock()
	value := dirty[key]
	dirtyMutex.RUnlock()
	return value
}

func SetDirty(key string) {
	dirtyMutex.Lock()
	dirty[key] = true
	dirtyMutex.Unlock()
}

func UnsetDirty(key string) {
	dirtyMutex.Lock()
	delete(dirty, key)
	dirtyMutex.Unlock()
}
