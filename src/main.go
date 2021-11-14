package main

import (
	"log"
	"net"

	"github.com/pornopatsan/acache/src/api"
	"github.com/pornopatsan/acache/src/server"
	"github.com/pornopatsan/acache/src/storage"
	"google.golang.org/grpc"
)

func main() {
	println("unimplemented")
	s := grpc.NewServer()
	srv := server.Create(storage.CreateDefault())
	api.RegisterACacheServer(s, srv)

	sock, err := net.Listen("tcp", ":8083")
	if err != nil {
		log.Fatal(err)
	}

	if err := s.Serve(sock); err != nil {
		log.Fatal(err)
	}
}
