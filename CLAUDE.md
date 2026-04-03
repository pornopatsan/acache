# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

ACache is a gRPC-based in-memory LRU cache server written in Go. It exposes Save, Get, and Remove RPCs defined in Protocol Buffers.

## Commands

```bash
# Run server (defaults: capacity=4096, port=8080)
go run src/main.go -c="4096" -p="8083"

# Run all tests
go test -race -timeout 60s ./...

# Run benchmarks (requires a running server)
go run bench/bench.go -c="4096" -p="8083"

# Regenerate protobuf (requires protoc, protoc-gen-go@v1.36.11, protoc-gen-go-grpc@v1.6.1)
protoc --go_out=. --go-grpc_out=. api/proto/item.proto
```

## Key Interfaces

- `Cache` (storage.go): `Save`, `Get`, `Remove` — implemented by `LruCache`
- `LruQueue` (dlist.go): `PushFront`, `PopBack`, `MoveFront`, `Remove` — implemented by `DLinkedList`
