package storage

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	cache := CreateDefault()
	assert.Zero(t, len(cache.items))
	assert.Zero(t, cache.lruIems.Size())
}

func TestSave(t *testing.T) {
	size := uint(32)
	cache := Create(size)

	for i := uint(0); i < size; i++ {
		key := fmt.Sprintf("%d", i)
		value := []byte{byte(i)}
		cache.Save(key, value)
		assert.Equal(t, cache.items[key].Key, key)
		assert.Equal(t, cache.items[key].Value, value)
		assert.Equal(t, cache.lruIems.Size(), i+1)
	}

	for i := uint(size); i < size+8; i++ {
		key := fmt.Sprintf("%d", i)
		removedKey := fmt.Sprintf("%d", i-size)
		value := []byte{byte(i)}
		cache.Save(key, value)
		assert.Equal(t, cache.items[key].Key, key)
		assert.Equal(t, cache.items[key].Value, value)
		assert.Equal(t, cache.lruIems.Size(), size)
		_, ok := cache.items[removedKey]
		assert.False(t, ok)
	}
}
