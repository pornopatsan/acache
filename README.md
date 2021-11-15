

## ACache
API Cache is a simple caching server, using grpc to accept messages. It allows to store key-value pairs, where key is `string` and value is `[]byte`.

**Start server**
```
go run src/main.go -c="4096" -p="8083"
```

**Benchmarks**

Run locally:  
```
go run bench/bench.go -c="4096" -p="8083"
```

Results:
```
Get  /    1 Client  RPS : 8299.20
Get  /   10 Clients RPS : 41231.20
Get  /  100 Clients RPS : 85064.00
Get  / 1000 Clients RPS : 77194.20
Save /    1 Client  RPS : 8741.80
Save /   10 Clients RPS : 33545.20
Save /  100 Clients RPS : 68292.60
Save / 1000 Clients RPS : 68824.40
```
MacBook Pro 2018 15-inch; 2,2 GHz 6-Core Intel Core i7; 16 GB 2400 MHz DDR4

**Realization**

Server: `src/server/server.go`  
Just a simple grpc server, as described in `api/item.proto`  

Storage: `src/server/server.go`  
HashMap over DoubleLinkedList. Values of HashMap are pointing to List Nodes. DoubleLinkedList is used as LruQueue. This gives O(1) time complexity to insert values in front, remove values from tail, geting value from anywhere and moving it to front & removing values.