package server

import (
	"context"
	"log"
	"sync"

	"github.com/pornopatsan/acache/src/api"
	"github.com/pornopatsan/acache/src/storage"
)

type ACacheServer struct {
	api.UnimplementedACacheServer
	cache   storage.Cache
	cacheMx sync.Mutex
}

func Create(cache storage.Cache) *ACacheServer {
	return &ACacheServer{
		cache: cache,
	}
}

func (self *ACacheServer) Save(ctx context.Context, item *api.Item) (*api.Response, error) {
	log.Printf("Save %s\n", item.Key)
	self.cacheMx.Lock()
	defer self.cacheMx.Unlock()
	if err := self.cache.Save(item.Key, item.Value); err != nil {
		return &api.Response{Status: api.Status_UNKNOWN_ERROR}, nil
	}
	return &api.Response{Status: api.Status_OK}, nil
}

func (self *ACacheServer) Get(ctx context.Context, key *api.Key) (*api.ItemResponse, error) {
	log.Printf("Get %s\n", key.Key)
	self.cacheMx.Lock()
	defer self.cacheMx.Unlock()
	item, err := self.cache.Get(key.Key)
	if err != nil {
		return &api.ItemResponse{Status: api.Status_KEY_NOT_FOUND}, nil
	}
	return &api.ItemResponse{
		Item: &api.Item{
			Key:   item.Key,
			Value: item.Value,
		},
		Status: api.Status_OK,
	}, nil
}

func (self *ACacheServer) Remove(ctx context.Context, key *api.Key) (*api.Response, error) {
	log.Printf("Remove %s\n", key.Key)
	self.cacheMx.Lock()
	defer self.cacheMx.Unlock()
	if err := self.cache.Remove(key.Key); err != nil {
		return &api.Response{Status: api.Status_KEY_NOT_FOUND}, nil
	}
	return &api.Response{Status: api.Status_OK}, nil
}
