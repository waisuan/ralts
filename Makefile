build:
	go build -o bin/ralts cmd/web/main.go

server:
	go run cmd/web/main.go

test:
	go test ./...

lint:
	golangci-lint run ./...