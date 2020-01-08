// Package lru provides a thread-safe least recently used
// (LRU) cache with a fixed capacity.
package lru

import (
	"errors"
	"sync"
)

// Cache is a thread-safe least recently used (LRU) cache with a fixed capacity.
type Cache struct {
	capacity  int
	load      int
	keyMap    map[interface{}]*lruNode
	lock      sync.RWMutex
	onEvicted func(key, value interface{})

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
	value interface{}
	// store the key to reverse lookup entry in keyMap during eviction
	key interface{}
}

// NewCache returns a new LRU Cache with the given capacity and eviction function.
// When the Cache evicts a key/value pair, the passed eviction function will be
// called with the evicted key/value pair as its arguments.
func NewCache(capacity int, onEvicted func(key, value interface{})) (*Cache, error) {
	if capacity < 1 {
		return nil, errors.New("capacity must be greater than 0")
	}

	return &Cache{
		capacity:  capacity,
		keyMap:    make(map[interface{}]*lruNode),
		lock:      sync.RWMutex{},
		onEvicted: onEvicted,
	}, nil
}

// Put inserts a key/value pair into the cache.
// If a value for the given key already exists in the cache, it will be overridden.
func (c *Cache) Put(key, value interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()

	existingNode, ok := c.keyMap[key]
	if ok {
		existingNode.value = value
		c.bringNodeToFront(existingNode)
		return
	}

	node := &lruNode{
		key:   key,
		value: value,
	}
	c.keyMap[key] = node

	// no need to evict, so place in front of list
	if c.load < c.capacity {
		if c.front == nil && c.rear == nil {
			c.front = node
			c.rear = node
		} else {
			c.insertInFront(node)
		}

		c.load = c.load + 1
		return
	}

	// load is equal to capacity, so need to evict the LRU
	delete(c.keyMap, c.rear.key)

	// call eviction function supplied in cache construction
	if c.onEvicted != nil {
		c.onEvicted(c.rear.key, c.rear.value)
	}

	// a single node is to be evicted
	if c.rear.next == nil && c.rear.prev == nil {
		c.front = node
		c.rear = node
		return
	}

	c.rear = c.rear.prev
	c.rear.next = nil
	c.insertInFront(node)
}

// Get returns the value stored in the cache for the given key.
// If there is no value cached for the given key, then nil, false is returned.
func (c *Cache) Get(key int) (interface{}, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	node, ok := c.keyMap[key]
	if !ok {
		return nil, false
	}

	// if node is at the front of the list or is the only node in the list
	if c.front == node || (node.prev == nil && node.next == nil) {
		return node.value, true
	}

	c.bringNodeToFront(node)
	return node.value, true
}

// insertInFront inserts the passed node into the front of the list
// used by the cache to track the usage of items in the cache.
func (c *Cache) insertInFront(node *lruNode) {
	c.front.prev = node
	node.next = c.front
	c.front = node
}

// Purge completely clears the cache.
// After a call to Purge, the length of the cache is 0.
func (c *Cache) Purge() {
	c.lock.Lock()
	defer c.lock.Unlock()

	// call eviction function and delete each key in the cache
	for k, v := range c.keyMap {
		if c.onEvicted != nil {
			c.onEvicted(k, v)
		}

		delete(c.keyMap, k)
	}

	// reset the lru list
	c.front = nil
	c.rear = nil

	// reset the load
	c.load = 0
}

// Len returns the number of key/value pairs that have been inserted into the cache.
func (c *Cache) Len() int {
	return len(c.keyMap)
}

// bringNodeToFront brings a node in the list used by the cache to
// track the usage of items in the cache to the front of the list.
func (c *Cache) bringNodeToFront(node *lruNode) {
	if node == c.front {
		return
	}

	// node is the last in the list
	if node.next == nil {
		c.rear = node.prev
	} else {
		// skip next prev to node prev
		node.next.prev = node.prev
	}

	// skip prev next to node next
	node.prev.next = node.next
	node.next = c.front
	node.prev = nil
	c.front.prev = node
	c.front = node
}
