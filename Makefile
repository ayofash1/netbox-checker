.PHONY: build test clean

BINARY_NAME=netbox-checker

build:
	go build -o $(BINARY_NAME) ./cmd

test:
	go test ./...

clean:
	rm -f $(BINARY_NAME)