.PHONY: all build run test clean

BINARY_NAME=tmobile-stats

all: build

build:
	go build -o $(BINARY_NAME) .

run:
	go run .

test:
	go test -v ./...

clean:
	rm -f $(BINARY_NAME)
