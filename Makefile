APP := ./cmd/api
BIN := ./bin/app

.PHONY: build run test clean redis-start

build:
	go build -o $(BIN) $(APP)

run:
	go run $(APP)

build-tmp:
	go build -o ./tmp/app $(APP)

dev:
	$(shell go env GOPATH)/bin/air

test:
	go test ./...

clean:
	rm -rf ./bin ./api
	rm -rf ./tmp

redis-start:
	docker compose up redis

