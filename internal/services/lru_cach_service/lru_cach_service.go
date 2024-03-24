package lru_cach_service

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	sync.RWMutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (lruCache *lruCache) Set(key Key, value interface{}) bool {
	lruCache.Lock()
	defer lruCache.Unlock()

	if item, ok := lruCache.items[key]; ok {
		item.Value = value
		lruCache.queue.MoveToFront(item)
		lruCache.deleteOldExtraItemIfExist()
		return true
	}

	item := lruCache.queue.PushFront(value)
	item.ExternalID = key
	lruCache.items[key] = item
	lruCache.deleteOldExtraItemIfExist()

	return false
}

func (lruCache *lruCache) Get(key Key) (interface{}, bool) {
	lruCache.RLock()
	defer lruCache.RUnlock()

	if item, ok := lruCache.items[key]; ok {
		lruCache.queue.MoveToFront(item)
		return item.Value, true
	}

	return nil, false
}

func (lruCache *lruCache) Clear() {
	lruCache.queue = NewList()
	lruCache.items = make(map[Key]*ListItem, lruCache.capacity)
}

func (lruCache *lruCache) deleteOldExtraItemIfExist() {
	if len(lruCache.items) > lruCache.capacity {
		delete(lruCache.items, lruCache.queue.Back().ExternalID)
	}

	if lruCache.queue.Len() > lruCache.capacity {
		lruCache.queue.Remove(lruCache.queue.Back())
	}
}
