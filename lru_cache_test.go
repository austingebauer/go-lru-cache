package lru

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCache(t *testing.T) {
	type args struct {
		capacity int
	}
	tests := []struct {
		name    string
		args    args
		want    *Cache
		wantErr bool
	}{
		{
			name: "test construction of a new cache",
			args: args{
				capacity: 1,
			},
			want: &Cache{
				capacity: 1,
				load:     0,
				keyMap:   make(map[interface{}]*lruNode),
				front:    nil,
				rear:     nil,
			},
			wantErr: false,
		},
		{
			name: "test construction of a new cache with capacity of 0",
			args: args{
				capacity: 0,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache, err := NewCache(tt.args.capacity, nil)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.want, cache)
			}
		})
	}
}

func TestLRUCache(t *testing.T) {
	type args struct {
		capacity int
	}

	tests := []struct {
		name       string
		args       args
		operations []string
		opArgs     [][]int
	}{
		{
			name: "puts and gets mixed for eviction",
			args: args{
				capacity: 2,
			},
			operations: []string{
				"Put",
				"Get",
			},
			opArgs: [][]int{
				{2, 1},
				{2, 1},
			},
		},
		{
			name: "puts and gets mixed for eviction",
			args: args{
				capacity: 2,
			},
			operations: []string{
				"Put",
				"Put",
				"Get",
				"Put",
				"Get",
				"Put",
				"Get",
				"Get",
				"Get",
			},
			opArgs: [][]int{
				{1, 1},
				{2, 2},
				{1, 1},
				{3, 3},
				{2, -1},
				{4, 5},
				{1, -1},
				{3, 3},
				{4, 5},
			},
		},
		{
			name: "puts and gets mixed for eviction",
			args: args{
				capacity: 1,
			},
			operations: []string{
				"Put",
				"Get",
				"Put",
				"Get",
				"Get",
			},
			opArgs: [][]int{
				{2, 1},
				{2, 1},
				{3, 2},
				{2, -1},
				{3, 2},
			},
		},
		{
			name: "puts and gets mixed for eviction",
			args: args{
				capacity: 2,
			},
			operations: []string{
				"Put",
				"Put",
				"Get",
			},
			opArgs: [][]int{
				{2, 1},
				{2, 2},
				{2, 2},
			},
		},
		{
			name: "puts and gets mixed for eviction",
			args: args{
				capacity: 2,
			},
			operations: []string{
				"Put",
				"Put",
				"Put",
				"Put",
				"Get",
				"Get",
			},
			opArgs: [][]int{
				{2, 1},
				{1, 1},
				{2, 3},
				{4, 1},
				{1, -1},
				{2, 3},
			},
		},
		{
			name: "puts and gets mixed for eviction",
			args: args{
				capacity: 2,
			},
			operations: []string{
				"Put",
				"Put",
				"Get",
				"Put",
				"Put",
				"Get",
			},
			opArgs: [][]int{
				{2, 1},
				{2, 2},
				{2, 2},
				{1, 1},
				{4, 1},
				{2, -1},
			},
		},
		{
			name: "puts and gets mixed for eviction",
			args: args{
				capacity: 2,
			},
			operations: []string{
				"Get",
				"Put",
				"Get",
				"Put",
				"Put",
				"Get",
				"Get",
			},
			opArgs: [][]int{
				{2, -1},
				{2, 6},
				{1, -1},
				{1, 5},
				{1, 2}, // evicts 2->6
				{1, 2},
				{2, 6},
			},
		},
		{
			name: "puts and gets mixed for eviction",
			args: args{
				capacity: 3,
			},
			operations: []string{
				"Put",
				"Put",
				"Put",
				"Put",
				"Get",
				"Get",
				"Get",
				"Get",
				"Put",
				"Get",
				"Get",
				"Get",
				"Get",
				"Get",
			},
			opArgs: [][]int{
				{1, 1},
				{2, 2},
				{3, 3},
				{4, 4},
				{4, 4},
				{3, 3},
				{2, 2},
				{1, -1},
				{5, 5},
				{1, -1},
				{2, 2},
				{3, 3},
				{4, -1},
				{5, 5},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache, err := NewCache(tt.args.capacity, nil)
			assert.NoError(t, err)
			assert.NotNil(t, cache)

			for i, op := range tt.operations {
				switch op {
				case "Put":
					cache.Put(tt.opArgs[i][0], tt.opArgs[i][1])
				case "Get":
					value, ok := cache.Get(tt.opArgs[i][0])

					// if the k/v pair is not in the cache, assert that
					// Get() returned return nil, false
					if tt.opArgs[i][1] == -1 {
						assert.Nil(t, value)
						assert.False(t, ok)
					} else {
						// otherwise, assert that the correct
						// value is returned for the key
						assert.Equal(t, tt.opArgs[i][1], value)
					}
				}
			}
		})
	}
}

func TestCacheEvictionFunction(t *testing.T) {
	keyVal := ""
	onEvict := func(key, value interface{}) {
		keyVal = key.(string) + value.(string)
	}

	cache, err := NewCache(1, onEvict)
	assert.NoError(t, err)
	cache.Put("ada", "lovelace")
	cache.Put("linus", "torvalds") // evicts ada->lovelace

	// Eviction happened on ada->lovelace, thereby calling onEvict()
	assert.Equal(t, "adalovelace", keyVal)
}
