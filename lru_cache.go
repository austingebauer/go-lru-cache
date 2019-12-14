// Package lru provides a thread-safe least recently used
// (LRU) cache with a fixed capacity.
package lru

import (
	"errors"
	"sync"
)

// Cache is a thread-safe least recently used (LRU) cache with a fixed capacity.
type Cache struct {
	capacity int
	load     int
	keyMap   map[interface{}]*lruNode
	lock     sync.RWMutex

	// font is always the latest used
	front *lruNode
	// rear is always the least recently used (LRU)
	rear *lruNode
}

// lruNode represents a single node in a doubly-linked
// list that's used by the Cache internally.
type lruNode struct {
	prev  *lruNode
	next  *lruNode
	key   interface{}
	value interface{}
}

// NewCache returns a new LRU Cache with the given capacity.
func NewCache(capacity int) (*Cache, error) {
	if capacity < 1 {
		return nil, errors.New("capacity must be greater than 0")
	}

	return &Cache{
		capacity: capacity,
		keyMap:   make(map[interface{}]*lruNode),
		lock:     sync.RWMutex{},
	}, nil
}

// Put inserts a key/value pair into the cache.
// If a value for the given key already exists in the cache, it will be overridden.
func (cache *Cache) Put(key, value interface{}) {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	existingNode, ok := cache.keyMap[key]
	if ok {
		existingNode.value = value
		cache.bringNodeToFront(existingNode)
		return
	}

	node := &lruNode{
		key:   key,
		value: value,
	}
	cache.keyMap[key] = node

	// no need to evict, so place in front of list
	if cache.load < cache.capacity {
		if cache.front == nil && cache.rear == nil {
			cache.front = node
			cache.rear = node
		} else {
			cache.insertInFront(node)
		}

		cache.load = cache.load + 1
		return
	}

	// load is equal to capacity, so need to evict the LRU
	evicted := cache.keyMap[cache.rear.key]
	delete(cache.keyMap, cache.rear.key)

	// a single node is to be evicted
	if evicted.next == nil && evicted.prev == nil {
		cache.front = node
		cache.rear = node
		return
	}

	cache.rear = cache.rear.prev
	cache.rear.next = nil
	cache.insertInFront(node)
}

// Get returns the value stored in the cache for the given key.
// If there is no value cached for the given key, then nil, false is returned.
func (cache *Cache) Get(key int) (interface{}, bool) {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	node, ok := cache.keyMap[key]
	if !ok {
		return nil, false
	}

	// if node is at the front of the list or is the only node in the list
	if cache.front == node || (node.prev == nil && node.next == nil) {
		return node.value, true
	}

	cache.bringNodeToFront(node)
	return node.value, true
}

// insertInFront inserts the passed node into the front of the list
// used by the cache to track the usage of items in the cache.
func (cache *Cache) insertInFront(node *lruNode) {
	cache.front.prev = node
	node.next = cache.front
	cache.front = node
}

// bringNodeToFront brings a node in the list used by the cache to
// track the usage of items in the cache to the front of the list.
func (cache *Cache) bringNodeToFront(node *lruNode) {
	if node == cache.front {
		return
	}

	// node is the last in the list
	if node.next == nil {
		cache.rear = node.prev
	} else {
		// skip next prev to node prev
		node.next.prev = node.prev
	}

	// skip prev next to node next
	node.prev.next = node.next
	node.next = cache.front
	node.prev = nil
	cache.front.prev = node
	cache.front = node
}
