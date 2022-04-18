package inmemory

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	// Error messages
	errCantGet       = "can't get an item from the cache"
	errCantSet       = "can't set an item to the cache"
	errCantDelete    = "can't delete an item from the cache"
	errKeyNotEqual   = "the key is not equal to expected"
	errValueNotEqual = "the value is not equal to expected"
)

var (
	// Examples values
	exampleKey   = "golang"
	exampleValue = []byte("1.18.1")
)

func TestCache_Get(t *testing.T) {
	cache := New()

	_ = cache.Set(exampleKey, exampleValue)

	k, v, err := cache.Get("golang")
	if !assert.NoError(t, err) {
		panic(errCantGet)
	}

	assert.Equal(t, exampleKey, k, errKeyNotEqual)
	assert.Equal(t, exampleValue, v, errValueNotEqual)
}

func TestCache_Set(t *testing.T) {
	cache := New()

	if err := cache.Set(exampleKey, exampleValue); !assert.NoError(t, err) {
		panic(errCantSet)
	}

	if _, _, err := cache.Get(exampleKey); !assert.NoError(t, err) {
		panic(errCantGet)
	}
}

func TestCache_Delete(t *testing.T) {
	cache := New()

	_ = cache.Set(exampleKey, exampleValue)

	if err := cache.Delete(exampleKey); !assert.NoError(t, err) {
		panic(errCantDelete)
	}
}

func TestCache_Size(t *testing.T) {
	cache := New()

	assert.Equal(t, uint(0), cache.Size())

	for i := 0; i < 10; i++ {
		key := fmt.Sprint("golang", rune(i))
		value := []byte(fmt.Sprint("1.18.1", rune(i)))

		_ = cache.Set(key, value)
	}

	assert.Equal(t, uint(10), cache.Size())

	_ = cache.Delete("golang2")

	assert.Equal(t, uint(9), cache.Size())
}
