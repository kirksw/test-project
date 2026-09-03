BIN := hello-fx
GOFLAGS :=

.PHONY: build test lint verify clean

build:
	go build -o $(BIN) .

test:
	go test -race -covermode=atomic -coverprofile=coverage.out ./...

lint:
	golangci-lint run

# Same checks the CI PR gate runs.
verify: lint test build

clean:
	rm -f $(BIN) coverage.out
