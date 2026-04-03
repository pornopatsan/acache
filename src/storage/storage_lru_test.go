package storage

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLRUEvictionOrder(t *testing.T) {
	cache := Create(8)

	// Fill cache with keys "0" through "7"
	for i := 0; i < 8; i++ {
		require.NoError(t, cache.Save(fmt.Sprintf("%d", i), []byte{byte(i)}))
	}
	assert.Equal(t, uint(8), cache.Size())

	// Access key "0" — moves it to front (most recently used)
	item, err := cache.Get("0")
	require.NoError(t, err)
	assert.Equal(t, []byte{byte(0)}, item.Value)

	// Save key "8" — must evict key "1" (oldest untouched), NOT key "0"
	require.NoError(t, cache.Save("8", []byte{8}))
	assert.Equal(t, uint(8), cache.Size())

	_, err = cache.Get("1") // key "1" should be evicted
	assert.IsType(t, KeyNotFoundError{}, err)

	item, err = cache.Get("0") // key "0" should still be present
	require.NoError(t, err)
	assert.Equal(t, []byte{byte(0)}, item.Value)

	// Save key "9" — must evict key "2" (next oldest untouched)
	require.NoError(t, cache.Save("9", []byte{9}))
	_, err = cache.Get("2")
	assert.IsType(t, KeyNotFoundError{}, err)

	// Validate internal list structure
	require.NoError(t, cache.lruIems.Validate())
}

func TestSaveExistingMovesToFront(t *testing.T) {
	cache := Create(4)

	require.NoError(t, cache.Save("A", []byte{1}))
	require.NoError(t, cache.Save("B", []byte{2}))
	require.NoError(t, cache.Save("C", []byte{3}))
	require.NoError(t, cache.Save("D", []byte{4}))

	// Re-save "A" with updated value — should move to front, not create duplicate
	require.NoError(t, cache.Save("A", []byte{10}))
	assert.Equal(t, uint(4), cache.Size())

	// Save "E" — should evict "B" (oldest since A was refreshed)
	require.NoError(t, cache.Save("E", []byte{5}))
	assert.Equal(t, uint(4), cache.Size())

	_, err := cache.Get("B")
	assert.IsType(t, KeyNotFoundError{}, err, "B should be evicted, not A")

	item, err := cache.Get("A")
	require.NoError(t, err)
	assert.Equal(t, []byte{10}, item.Value, "A should have updated value")

	require.NoError(t, cache.lruIems.Validate())
}

func TestLRUMixedOperations(t *testing.T) {
	cache := Create(4)

	// Fill: [D, C, B, A] (front to back)
	require.NoError(t, cache.Save("A", []byte{1}))
	require.NoError(t, cache.Save("B", []byte{2}))
	require.NoError(t, cache.Save("C", []byte{3}))
	require.NoError(t, cache.Save("D", []byte{4}))

	// Get "A" → moves to front: [A, D, C, B]
	_, err := cache.Get("A")
	require.NoError(t, err)

	// Remove "C" → [A, D, B], size=3
	require.NoError(t, cache.Remove("C"))
	assert.Equal(t, uint(3), cache.Size())

	// Save "E" → [E, A, D, B], size=4
	require.NoError(t, cache.Save("E", []byte{5}))
	assert.Equal(t, uint(4), cache.Size())

	// Save "F" → must evict "B" (back): [F, E, A, D]
	require.NoError(t, cache.Save("F", []byte{6}))
	assert.Equal(t, uint(4), cache.Size())

	_, err = cache.Get("B")
	assert.IsType(t, KeyNotFoundError{}, err, "B should be evicted")

	// All of D, A, E, F should be present
	for _, key := range []string{"D", "A", "E", "F"} {
		_, err := cache.Get(key)
		require.NoError(t, err, "key %s should be present", key)
	}

	require.NoError(t, cache.lruIems.Validate())
}

// referenceLRU is a simple slice-based LRU for verifying eviction order.
type referenceLRU struct {
	keys     []string
	capacity int
}

func (r *referenceLRU) touch(key string) {
	for i, k := range r.keys {
		if k == key {
			r.keys = append(r.keys[:i], r.keys[i+1:]...)
			break
		}
	}
	r.keys = append([]string{key}, r.keys...)
	if len(r.keys) > r.capacity {
		r.keys = r.keys[:r.capacity]
	}
}

func (r *referenceLRU) remove(key string) {
	for i, k := range r.keys {
		if k == key {
			r.keys = append(r.keys[:i], r.keys[i+1:]...)
			return
		}
	}
}

func (r *referenceLRU) contains(key string) bool {
	for _, k := range r.keys {
		if k == key {
			return true
		}
	}
	return false
}

func (r *referenceLRU) evictedOnSave(key string) string {
	// If key already exists, no eviction
	if r.contains(key) {
		return ""
	}
	if len(r.keys) < r.capacity {
		return ""
	}
	return r.keys[len(r.keys)-1]
}

func TestLRUEvictionOrderRandomized(t *testing.T) {
	seeds := []int64{42, 123, 9999, 314159, 271828}

	for _, seed := range seeds {
		t.Run(fmt.Sprintf("seed_%d", seed), func(t *testing.T) {
			const capacity = 32
			const ops = 500
			rng := rand.New(rand.NewSource(seed))

			cache := Create(capacity)
			ref := &referenceLRU{capacity: capacity}

			allKeys := make([]string, 0, capacity*2)
			for i := 0; i < capacity*2; i++ {
				allKeys = append(allKeys, fmt.Sprintf("key-%d", i))
			}

			for i := 0; i < ops; i++ {
				op := rng.Intn(3) // 0=Save, 1=Get, 2=Remove
				keyIdx := rng.Intn(len(allKeys))
				key := allKeys[keyIdx]

				switch op {
				case 0: // Save
					evicted := ref.evictedOnSave(key)
					value := []byte(fmt.Sprintf("val-%d-%d", seed, i))
					require.NoError(t, cache.Save(key, value))
					ref.touch(key)

					if evicted != "" {
						_, err := cache.Get(evicted)
						assert.IsType(t, KeyNotFoundError{}, err,
							"op %d: key %q should have been evicted by saving %q", i, evicted, key)
						// Re-touch in ref since Get of missing key doesn't change state
					}

				case 1: // Get
					item, err := cache.Get(key)
					if ref.contains(key) {
						require.NoError(t, err, "op %d: key %q should be in cache", i, key)
						assert.Equal(t, key, item.Key)
						ref.touch(key)
					} else {
						assert.IsType(t, KeyNotFoundError{}, err,
							"op %d: key %q should NOT be in cache", i, key)
					}

				case 2: // Remove
					err := cache.Remove(key)
					if ref.contains(key) {
						require.NoError(t, err, "op %d: key %q should be removable", i, key)
						ref.remove(key)
					} else {
						assert.IsType(t, KeyNotFoundError{}, err)
					}
				}
			}

			require.NoError(t, cache.lruIems.Validate())
			assert.Equal(t, uint(len(ref.keys)), cache.Size())
		})
	}
}
