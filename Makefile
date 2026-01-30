.PHONY: build test lint clean

build:
	@echo "Building WatchDoc"
	@go build -o watchdoc

test:
	@go test -race -v ./...

lint:
	@go vet ./...
	@if command -v golangci-lint > /dev/null 2>&1; then golangci-lint run; fi

clean:
	@rm -f watchdoc
