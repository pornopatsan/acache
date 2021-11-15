package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/pornopatsan/acache/src/api"
	"google.golang.org/grpc"
)

var DATUM = []byte{ // All values are equal and 256 bytes long
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, // 16 byte
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, // 32 byte
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, // 64 byte
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, // 64 byte
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, // 128 byte
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, // 256 byte
}

func getKey(i int) string {
	return fmt.Sprintf("%02d%04d%010d", i, i*i, i*i*i) // All keys are 16 bytes long
}

func randomKey(cacheCapacity int) string {
	return getKey(rand.Int() % cacheCapacity)
}

func printResult(str string, rps float64) {
	fmt.Printf("%s RPS : %5.2f\n", str, rps)
}

func prepareCache(conn *grpc.ClientConn, capacity int) {
	client := api.NewACacheClient(conn)
	for i := 0; i < capacity; i++ {
		client.Save(context.Background(), &api.Item{Key: getKey(i), Value: DATUM})
	}
}

func testGetSingleClient(conn *grpc.ClientConn, cacheCapacity int) float64 {
	iterations := 0
	client := api.NewACacheClient(conn)
	timer := time.NewTimer(5 * time.Second)
	for {
		select {
		case <-timer.C:
			return float64(iterations) / 5.0
		default:
			{
				_, err := client.Get(context.Background(), &api.Key{Key: randomKey(cacheCapacity)})
				if err != nil {
					log.Fatal("Get request failed: ", err)
				}
				iterations++
			}
		}
	}
}

func testSaveSingleClient(conn *grpc.ClientConn, cacheCapacity int) float64 {
	iterations := 0
	client := api.NewACacheClient(conn)
	timer := time.NewTimer(5 * time.Second)
	for {
		select {
		case <-timer.C:
			return float64(iterations) / 5.0
		default:
			{
				_, err := client.Save(
					context.Background(),
					&api.Item{Key: randomKey(2 * cacheCapacity), Value: DATUM},
				)
				if err != nil {
					log.Fatal("Save request failed: ", err)
				}
				iterations++
			}
		}
	}
}

type singleClientTest = func(*grpc.ClientConn, int) float64

func testMultClients(conn *grpc.ClientConn, cacheCapacity int, clients int, method singleClientTest) float64 {
	resChan := make(chan float64, clients)
	for i := 0; i < clients; i++ {
		go func() {
			res := method(conn, cacheCapacity)
			resChan <- res
		}()
	}
	res := float64(0.0)
	for i := 0; i < clients; i++ {
		res += <-resChan
	}
	return res
}

func main() {
	var capacity int
	flag.IntVar(&capacity, "c", 4096, "Capacity of ACache server (to utilize it propperly)")
	port := flag.Int("p", 8080, "Port, on which server is located")
	flag.Parse()

	conn, err := grpc.Dial(fmt.Sprintf(":%d", *port), grpc.WithInsecure())
	if err != nil {
		log.Fatal("Cannof find server: ", err)
	}

	prepareCache(conn, capacity)
	printResult("Get  /    1 Client ", testMultClients(conn, capacity, 1, testGetSingleClient))
	printResult("Get  /   10 Clients", testMultClients(conn, capacity, 10, testGetSingleClient))
	printResult("Get  /  100 Clients", testMultClients(conn, capacity, 100, testGetSingleClient))
	printResult("Get  / 1000 Clients", testMultClients(conn, capacity, 1000, testGetSingleClient))
	printResult("Save /    1 Client ", testMultClients(conn, capacity, 1, testSaveSingleClient))
	printResult("Save /   10 Clients", testMultClients(conn, capacity, 10, testSaveSingleClient))
	printResult("Save /  100 Clients", testMultClients(conn, capacity, 100, testSaveSingleClient))
	printResult("Save / 1000 Clients", testMultClients(conn, capacity, 1000, testSaveSingleClient))
}
