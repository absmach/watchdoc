.PHONY: build test lint clean

build:
	@echo "Building WatchDoc"
	@go build -o watchdoc

test:
	@go test -race -v ./...

lint:
	@golangci-lint run

clean:
	@rm -f watchdoc
