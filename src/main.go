package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/pornopatsan/acache/src/api"
	"github.com/pornopatsan/acache/src/server"
	"github.com/pornopatsan/acache/src/storage"
	"google.golang.org/grpc"
)

func main() {
	capacity := flag.Uint("c", 4096, "Capacity of ACache server")
	port := flag.Int("p", 8080, "Port, on which server wiil be located")
	flag.Parse()

	s := grpc.NewServer()
	srv := server.Create(storage.Create(*capacity))
	api.RegisterACacheServer(s, srv)

	sock, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal(err)
	}

	if err := s.Serve(sock); err != nil {
		log.Fatal(err)
	}
}
