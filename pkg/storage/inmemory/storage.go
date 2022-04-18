package inmemory

import "sync"

func New() *Cache {
	return &Cache{
		items: make(map[string][]byte),
		len:   0,
	}
}

type Cache struct {
	items map[string][]byte
	len   int64
	sync.Mutex
}

func (c *Cache) Size() int64 {
	return c.len
}

func (c *Cache) Set(key string, value []byte) error {
	c.Lock()
	defer c.Unlock()

	c.items[key] = value
	c.len++
	return nil
}

func (c *Cache) Get(key string) (string, []byte, error) {
	c.Lock()
	defer c.Unlock()

	if value, ok := c.items[key]; !ok {
		return "", nil, KeyNotFoundError{Key: key}
	} else {
		return key, value, nil
	}
}

func (c *Cache) Delete(key string) error {
	c.Lock()
	defer c.Unlock()

	delete(c.items, key)
	c.len--
	return nil
}
