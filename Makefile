APP := ./cmd/api
BIN := ./bin/app

.PHONY: build run test clean

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

build-mac:
	GOARCH=arm64 GOOS=darwin go build -o ./bin/app-macos $(APP)

start-redis:
	docker compose up redis

stop-redis:
	docker compose down redis
