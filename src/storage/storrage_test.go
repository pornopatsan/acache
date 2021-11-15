package storage

import (
	"fmt"
	"testing"

	"github.com/pornopatsan/acache/src/lib/dlist"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	cache := CreateDefault()
	assert.Zero(t, len(cache.items))
	assert.Zero(t, cache.lruIems.Size())
	assert.Zero(t, cache.Size())
}

func TestSave(t *testing.T) {
	capacity := uint(32)
	cache := Create(capacity)

	// Save some
	for i := uint(0); i < capacity; i++ {
		key := fmt.Sprintf("%d", i)
		value := []byte{byte(i)}
		cache.Save(key, value)
		assert.Equal(t, cache.items[key].Key, key)
		assert.Equal(t, cache.items[key].Value, value)
		assert.Equal(t, cache.Size(), i+1)
	}

	// Save with overflow
	for i := uint(capacity); i < capacity+8; i++ {
		key := fmt.Sprintf("%d", i)
		removedKey := fmt.Sprintf("%d", i-capacity)
		value := []byte{byte(i)}
		cache.Save(key, value)
		assert.Equal(t, cache.items[key].Key, key)
		assert.Equal(t, cache.items[key].Value, value)
		assert.Equal(t, cache.Size(), capacity)
		_, ok := cache.items[removedKey]
		assert.False(t, ok)
	}
}

func TestRemove(t *testing.T) {
	capacity := uint(32)
	cache := Create(capacity)

	{ // Remove from empty storage
		key := fmt.Sprintf("k1")
		err := cache.Remove(key)
		assert.IsType(t, err, KeyNotFoundError{})
		_, ok := cache.items[key]
		assert.False(t, ok)
		assert.Equal(t, uint(0), cache.lruIems.Size())
	}

	{ // Save & Remove single
		key := fmt.Sprintf("k2")
		cache.Save(key, []byte{1, 2, 3, 4, 5})
		err := cache.Remove(key)
		assert.Nil(t, err)
		_, ok := cache.items[key]
		assert.False(t, ok)
		assert.Equal(t, uint(0), cache.Size())
	}

	// Save some and then remove
	for i := uint(0); i < capacity/2; i++ {
		key := fmt.Sprintf("%d", i)
		value := []byte{byte(i)}
		cache.Save(key, value)
	}
	assert.Equal(t, uint(capacity/2), cache.Size())
	for i := uint(0); i < capacity/2; i++ {
		key := fmt.Sprintf("%d", i)
		err := cache.Remove(key)
		assert.Nil(t, err)
		_, ok := cache.items[key]
		assert.False(t, ok)
		assert.Equal(t, uint(capacity/2-i-1), cache.Size())
	}
}

func TestGet(t *testing.T) {
	capacity := uint(32)
	cache := Create(capacity)
	nodes := make([]*dlist.Node, 2*capacity)

	{ // Get from empty
		key := fmt.Sprintf("k1")
		item, err := cache.Get(key)
		assert.IsType(t, err, KeyNotFoundError{})
		assert.Nil(t, item)
	}

	// Save some and then try Get
	for i := uint(0); i < capacity; i++ {
		key := fmt.Sprintf("%d", i)
		value := []byte{byte(i)}
		cache.Save(key, value)
		nodes[i] = cache.items[key]
	}
	for i := uint(0); i < capacity; i++ {
		key := fmt.Sprintf("%d", i)
		node, err := cache.Get(key)
		assert.Nil(t, err)
		assert.Equal(t, node.Key, key)
		assert.Equal(t, node.Value, []byte{byte(i)})
		nodes[i] = cache.items[key]
	}

	// Test lru update. Add capacity-2 objects and test, that only these two left from orginals
	cache.Get("5")
	cache.Get("10")

	// Test overflow
	for i := uint(capacity); i < 2*capacity; i++ {
		if i-capacity == 5 || i-capacity == 10 {
			continue
		}
		key := fmt.Sprintf("%d", i)
		value := []byte{byte(i)}
		cache.Save(key, value)
		nodes[i] = cache.items[key]
	}
	for i := uint(capacity); i < 2*capacity; i++ {
		if i-capacity == 5 || i-capacity == 10 {
			continue
		}
		removedKey := fmt.Sprintf("%d", i-capacity)
		key := fmt.Sprintf("%d", i)

		_, err := cache.Get(removedKey)
		assert.IsType(t, KeyNotFoundError{}, err)

		node, err := cache.Get(key)
		assert.Nil(t, err)
		assert.Equal(t, node.Key, key)
		assert.Equal(t, node.Value, []byte{byte(i)})
		nodes[i] = cache.items[key]
	}
	{ // Test saved two values
		_, err := cache.Get("5")
		assert.Nil(t, err)
		nodes[5] = cache.items["5"]
	}
	{
		_, err := cache.Get("10")
		assert.Nil(t, err)
		nodes[10] = cache.items["10"]
	}
}
