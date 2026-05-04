APP := ./cmd/api
BIN := ./bin/app

.PHONY: build run test clean

build:
	go build -o $(BIN) $(APP)

run:
	go run $(APP)

dev: build
	$(BIN)

test:
	go test ./...

clean:
	rm -f $(BIN)

build-mac:
	GOARCH=arm64 GOOS=darwin go build -o ./bin/app-macos $(APP)
