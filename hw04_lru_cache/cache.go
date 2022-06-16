package hw04lrucache

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	lock     sync.Mutex
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	if li, ok := c.items[key]; ok {
		li.Value = cacheItem{key: key, value: value}
		c.queue.MoveToFront(li)
		return true
	}
	if c.queue.Len() >= c.capacity {
		purge := c.queue.Back()
		delete(c.items, purge.Value.(cacheItem).key)
		c.queue.Remove(purge)
	}
	li := c.queue.PushFront(cacheItem{key: key, value: value})
	c.items[key] = li
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if li, ok := c.items[key]; ok {
		c.queue.MoveToFront(li)
		return li.Value.(cacheItem).value, true
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
