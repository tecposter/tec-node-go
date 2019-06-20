package user

import (
	"sync"
)

type cacheStorage struct {
	sm sync.Map
}

func newCacheStorage() *cacheStorage {
	return &cacheStorage{}
}

func (cache *cacheStorage) set(key, val interface{}) {
	cache.sm.Store(key, val)
}

func (cache *cacheStorage) get(key interface{}) (interface{}, bool) {
	return cache.sm.Load(key)
}

func (cache *cacheStorage) remove(key interface{}) {
	cache.sm.Delete(key)
}
