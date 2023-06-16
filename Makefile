storage: db1 redis

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

mig_create:
	migrate create -ext sql -dir db/migrations -seq $(change)

mig_up:
	 migrate -source file://db/migrations -database postgres://postgres:postgres@localhost:5432/ralts up 1

lint:
	golangci-lint run ./...

gen_mock:
	mockery --name=$(interface)