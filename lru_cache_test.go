package lru

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLRUCache_1(t *testing.T) {
	cache := NewCache(2)
	cache.Put(1, 1)
	cache.Put(2, 2)
	assert.Equal(t, 1, cache.Get(1))
	cache.Put(3, 3)
	assert.Equal(t, -1, cache.Get(2))
	cache.Put(4, 4)
	assert.Equal(t, -1, cache.Get(1))
	assert.Equal(t, 3, cache.Get(3))
	assert.Equal(t, 4, cache.Get(4))
}

func TestLRUCache_2(t *testing.T) {
	cache := NewCache(1)
	cache.Put(2, 1)
	assert.Equal(t, 1, cache.Get(2))
}

func TestLRUCache_3(t *testing.T) {
	cache := NewCache(1)
	cache.Put(2, 1)
	assert.Equal(t, 1, cache.Get(2))

	cache.Put(3, 2)
	assert.Equal(t, -1, cache.Get(2))
	assert.Equal(t, 2, cache.Get(3))
}

func TestLRUCache_4(t *testing.T) {
	cache := NewCache(2)
	cache.Put(2, 1)
	cache.Put(2, 2)
	assert.Equal(t, 2, cache.Get(2))
}

func TestLRUCache_5(t *testing.T) {
	cache := NewCache(2)
	cache.Put(2, 1)
	cache.Put(1, 1)
	cache.Put(2, 3)
	cache.Put(4, 1)
	assert.Equal(t, -1, cache.Get(1))
	assert.Equal(t, 3, cache.Get(2))
}

func TestLRUCache_6(t *testing.T) {
	cache := NewCache(2)
	cache.Put(2, 1)
	cache.Put(2, 2)
	assert.Equal(t, 2, cache.Get(2))
	cache.Put(1, 1)
	cache.Put(4, 1)
	assert.Equal(t, -1, cache.Get(2))
}

func TestLRUCache_7(t *testing.T) {
	cache := NewCache(2)
	assert.Equal(t, -1, cache.Get(2))
	cache.Put(2, 6)
	assert.Equal(t, -1, cache.Get(1))
	cache.Put(1, 5)
	cache.Put(1, 2) // evicts 2->6
	assert.Equal(t, 2, cache.Get(1))
	assert.Equal(t, 6, cache.Get(2))
}

func TestLRUCache_8(t *testing.T) {
	cache := NewCache(3)
	cache.Put(1, 1)
	cache.Put(2, 2)
	cache.Put(3, 3)
	cache.Put(4, 4)
	assert.Equal(t, 4, cache.Get(4))
	assert.Equal(t, 3, cache.Get(3))
	assert.Equal(t, 2, cache.Get(2))
	assert.Equal(t, -1, cache.Get(1))
	cache.Put(5, 5)
	assert.Equal(t, -1, cache.Get(1))
	assert.Equal(t, 2, cache.Get(2))
	assert.Equal(t, 3, cache.Get(3))
	assert.Equal(t, -1, cache.Get(4))
	assert.Equal(t, 5, cache.Get(5))
}
