.PHONY: all build run test clean

BINARY_NAME=tmobile-stats

all: build

build:
	go build -o $(BINARY_NAME) main.go

run:
	go run main.go

test:
	go test -v ./...

clean:
	rm -f $(BINARY_NAME)
