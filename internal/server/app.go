package server

import (
	"fmt"
	"log"
	"net"

	"acache/internal/handler"
	pb "acache/pkg/gen/proto"
	"acache/pkg/storage/inmemory"
	"acache/pkg/storage/memcached"
	"acache/pkg/storage/redis"

	"google.golang.org/grpc"
)

type Options struct {
	ServePort int
	CacheType string
	CachePath string
}

func Run(opt *Options) error {
	var cache handler.CacheRecorder

	log.Print("[info] cache init")
	switch opt.CacheType {
	case "inmemory":
		cache = inmemory.New()
	case "redis":
		cache = redis.New(opt.CachePath)
	case "memcached":
		cache = memcached.New(opt.CachePath)
	default:
		cache = inmemory.New()
	}

	log.Print("[info] grpc controllers init")
	controllers := grpc.NewServer()
	pb.RegisterACacheServer(controllers, handler.New(cache))

	sock, err := net.Listen("tcp", fmt.Sprintf(":%d", opt.ServePort))
	if err != nil {
		return err
	}

	log.Print("[info] grpc controllers serve")
	return controllers.Serve(sock)
}
