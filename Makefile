build:
	go build -o bin/ralts cmd/web/main.go

server:
	go run cmd/web/main.go

web:
	cd client; npm start

test:
	go test ./...

db1:
	sudo service postgresql start

redis:
	redis-server --daemonize yes

lint:
	golangci-lint run ./...