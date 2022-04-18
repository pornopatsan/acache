package main

import (
	pb "acache/pkg/gen/proto"
	"context"
	"flag"
	"google.golang.org/grpc"
	"log"
)

var (
	flagDescriptionGrpcServerPath = "The gRPC server path we will use. Default is \"localhost:8080\""

	grpcPath = flag.String("path", "localhost:8080", flagDescriptionGrpcServerPath)
)

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*grpcPath, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewACacheClient(conn)

	// Set example
	item := &pb.Item{Key: "fff", Value: []byte("FFF")}
	setResp, err := client.Set(context.Background(), item)
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("Received: %+v \n", setResp)

	// Get example
	getResp, err := client.Get(context.Background(), &pb.Key{Key: "fff"})
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("Received: %+v \n", getResp)

	// Size example
	sizeResp1, err := client.Size(context.Background(), &pb.Empty{})
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("Received: %+v \n", sizeResp1)

	// Delete example
	delResp, err := client.Delete(context.Background(), &pb.Key{Key: "fff"})
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("Received: %+v \n", delResp)

	// Size example
	sizeResp2, err := client.Size(context.Background(), &pb.Empty{})
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("Received: %+v \n", sizeResp2)
}
