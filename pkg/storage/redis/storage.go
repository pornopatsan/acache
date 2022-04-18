package redis

import (
	"context"

	"acache/pkg/util/errors"

	"github.com/go-redis/redis/v8"
)

const (
	defaultPassword = ""
	defaultDB       = 0
)

func New(path string) *Cache {
	return &Cache{client: redis.NewClient(&redis.Options{
		Addr:     path,
		Password: defaultPassword,
		DB:       defaultDB,
	})}
}

type Cache struct {
	client *redis.Client
	len    int64
}

func (c *Cache) Set(key string, value []byte) error {
	ctx := context.Background()
	err := c.client.Set(ctx, key, string(value), 0).Err()
	if err != nil {
		return errors.HandleError(err)
	}

	c.len++
	return nil
}

func (c *Cache) Get(key string) (string, []byte, error) {
	ctx := context.Background()
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return "", nil, errors.HandleError(err)
	}

	c.len++

	return key, []byte(val), nil
}

func (c *Cache) Delete(key string) error {
	ctx := context.Background()
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		return errors.HandleError(err)
	}

	c.len--

	return nil
}

func (c *Cache) Size() int64 {
	return c.len
}
