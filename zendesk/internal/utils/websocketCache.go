package utils

import (
	"encoding/json"
	"sync"
	"time"
)

type Cache[Key comparable, T any] interface {
	Get(key Key) (T, bool)
	GetAll() (CacheResponse[Key, T], error)
	Update(key Key, item T)
	Delete(key Key)
	GetCacheAge() *time.Duration
	Reset()
}

type CacheResponse[Key comparable, T any] struct {
	Items map[Key]T `json:"items"`
}

type MemoryCacheInstance[Key comparable, T any] struct {
	mutex       *sync.Mutex
	lastUpdated *time.Time
	items       map[Key]T
}

func (cache *MemoryCacheInstance[Key, T]) Reset() {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cache.lastUpdated = nil
	cache.items = map[Key]T{}
}

func (cache *MemoryCacheInstance[Key, T]) Update(
	key Key,
	update T,
) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	now := time.Now()
	cache.lastUpdated = &now

	cache.items[key] = update
}

func (cache *MemoryCacheInstance[Key, T]) Delete(key Key) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	delete(cache.items, key)
}

func (cache *MemoryCacheInstance[Key, T]) Get(key Key) (T, bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	item, ok := cache.items[key]
	if !ok {
		return *new(T), false
	}

	itemBytes, err := json.Marshal(item)
	if err != nil {
		panic(err)
	}

	target := *new(T)
	if err := json.Unmarshal(itemBytes, &target); err != nil {
		panic(err)
	}

	return item, true
}

func (cache *MemoryCacheInstance[Key, T]) GetAll() (CacheResponse[Key, T], error) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	target := map[Key]T{}

	itemBytes, err := json.Marshal(cache.items)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(itemBytes, &target); err != nil {
		panic(err)
	}

	return CacheResponse[Key, T]{
		Items: target,
	}, nil
}

func (cache *MemoryCacheInstance[Key, T]) GetLastUpdated() *time.Time {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	return cache.lastUpdated
}

func (cache *MemoryCacheInstance[Key, T]) GetCacheAge() *time.Duration {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	if cache.lastUpdated != nil {
		age := time.Since(*cache.lastUpdated)

		return &age
	}

	return nil
}

func NewMemoryCacheInstance[Key comparable, T any]() *MemoryCacheInstance[Key, T] {
	return &MemoryCacheInstance[Key, T]{
		mutex:       &sync.Mutex{},
		items:       map[Key]T{},
		lastUpdated: nil,
	}
}
