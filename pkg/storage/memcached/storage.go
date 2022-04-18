package memcached

import (
	"acache/pkg/util/errors"
	"github.com/bradfitz/gomemcache/memcache"
)

func New(path string) *Cache {
	return &Cache{client: memcache.New(path)}
}

type Cache struct {
	client *memcache.Client
	len    int64
}

func (c *Cache) Set(key string, value []byte) error {
	err := c.client.Set(&memcache.Item{Key: key, Value: value})
	if err != nil {
		return errors.HandleError(err)
	}

	c.len++

	return nil
}

func (c *Cache) Get(key string) (string, []byte, error) {
	item, err := c.client.Get(key)
	if err != nil {
		return "", nil, errors.HandleError(err)
	}
	return item.Key, item.Value, nil
}

func (c *Cache) Delete(key string) error {
	err := c.client.Delete(key)
	if err != nil {
		return errors.HandleError(err)
	}

	c.len--

	return nil
}

func (c *Cache) Size() int64 {
	return c.len
}
