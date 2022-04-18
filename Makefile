build: cleaning
	mkdir -p build
	go build -o /build/acache /cmd/acache/main.go

cleaning:
	go fmt ./...
	gofumpt -l -w .
	golangci-lint run ./...

clear:
	rm -r pkg/gen
	rm -r bin

client_run:
	go run cmd/client/main.go

debug: cleaning
	go run ./cmd/acache/main.go

dependencies:
	go mod vendor

docker_compose_up: cleaning dependencies
	docker-compose up

docker_compose_down:
	docker-compose down

generate:
	protoc -I. -I .  \
		--go_out ./pkg/gen/ \
		--go_opt paths=source_relative \
		--go-grpc_out ./pkg/gen/ \
		--go-grpc_opt paths=source_relative \
		./proto/api.proto
