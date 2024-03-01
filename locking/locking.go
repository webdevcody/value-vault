package locking

import (
	"sync"
)

type KeyValueStore struct {
	locks map[string]*sync.RWMutex
	mutex sync.Mutex // Mutex to protect access to the locks map
}

func NewKeyValueStore() *KeyValueStore {
	return &KeyValueStore{
		locks: make(map[string]*sync.RWMutex),
	}
}

var kvs = NewKeyValueStore()

func RLock(key string) {
	kvs.mutex.Lock()
	if kvs.locks[key] == nil {
		kvs.locks[key] = &sync.RWMutex{}
	}
	kvs.mutex.Unlock()
	lock := kvs.locks[key]
	lock.RLock()
}

func RUnlock(key string) {
	kvs.mutex.Lock()
	if kvs.locks[key] == nil {
		kvs.locks[key] = &sync.RWMutex{}
	}
	kvs.mutex.Unlock()
	lock := kvs.locks[key]
	lock.RUnlock()
}

func Lock(key string) {
	kvs.mutex.Lock()
	if kvs.locks[key] == nil {
		kvs.locks[key] = &sync.RWMutex{}
	}
	kvs.mutex.Unlock()
	lock := kvs.locks[key]
	lock.Lock()
}

func Unlock(key string) {
	kvs.mutex.Lock()
	if kvs.locks[key] == nil {
		kvs.locks[key] = &sync.RWMutex{}
	}
	kvs.mutex.Unlock()
	lock := kvs.locks[key]
	lock.Unlock()
}
