package main

import (
	"flag"
	"log"

	"acache/pkg/util/errors"

	"acache/internal/server"
)

const (
	flagDescriptionCachePath  = "The cache path we will use. Default is \"localhost:11211\""
	flagDescriptionCacheType  = "The cache to be used for the handler. Default is \"memcached\""
	flagDescriptionServerPort = "Port, on which handler will be located. Default is \":8080\""
)

func main() {
	var (
		servePort = flag.Int("p", 8080, flagDescriptionServerPort)
		cacheType = flag.String("cache_type", "memcached", flagDescriptionCacheType)
		cachePath = flag.String("cache_path", "localhost:11211", flagDescriptionCachePath)
	)
	flag.Parse()

	log.Printf("[info] acache server init")
	log.Printf("[info] serve port: %v", *servePort)
	err := server.Run(&server.Options{
		ServePort: *servePort,
		CacheType: *cacheType,
		CachePath: *cachePath,
	})
	if err != nil {
		log.Fatal(errors.HandleError(err))
	}
}
