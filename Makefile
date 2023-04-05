build:
	go build -o bin/ralts cmd/web/main.go

run:
	go run cmd/web/main.go

test:
	go test ./...