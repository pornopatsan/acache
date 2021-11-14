package storage

import (
	"fmt"

	"github.com/pornopatsan/acache/src/lib/dlist"
)

const (
	DEFAULT_SIZE = uint(4096)
)

type Item struct {
	Key   string
	Value []byte
}

type Cache interface {
	Save(string, []byte) error
	Get(string) (Item, error)
	Remove(string) error
}

type LruCache struct {
	items   map[string]dlist.Node
	lruIems dlist.LruQueue
	size    uint
}

func Create(size uint) *LruCache {
	return &LruCache{
		items:   make(map[string]dlist.Node),
		lruIems: dlist.Create(),
		size:    size,
	}
}

func CreateDefault() *LruCache {
	return Create(DEFAULT_SIZE)
}

func (self *LruCache) removeOldest() error {
	node, err := self.lruIems.PopBack()
	if err != nil {
		return err
	}
	delete(self.items, node.Key)
	return nil
}

func (self *LruCache) Save(key string, value []byte) error {
	node, exists := self.items[key]
	if !exists {
		if self.lruIems.Size() >= self.size {
			if err := self.removeOldest(); err != nil {
				return err
			}
		}
		newNode := dlist.CreateNode(key, value)
		self.items[key] = *newNode
		self.lruIems.PushFront(newNode)
	} else {
		self.lruIems.MoveFront(&node)
		node.Value = value
	}
	return nil
}

func (self *LruCache) Get(key string) (Item, error) {
	node, exists := self.items[key]
	if !exists {
		return Item{}, KeyNotFoundError{key: key}
	}
	self.lruIems.MoveFront(&node)
	return Item{Key: key, Value: node.Value}, nil
}

func (self *LruCache) Remove(key string) error {
	node, exists := self.items[key]
	if !exists {
		return KeyNotFoundError{key: key}
	}
	self.lruIems.Remove(&node)
	delete(self.items, key)
	return nil
}

type KeyAlreadyExsitsError struct {
	key string
}

func (self KeyAlreadyExsitsError) Error() string {
	return fmt.Sprintf("Key `%s` already exists", self.key)
}

type KeyNotFoundError struct {
	key string
}

func (self KeyNotFoundError) Error() string {
	return fmt.Sprintf("Key `%s` not found", self.key)
}
