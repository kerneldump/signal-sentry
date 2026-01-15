.PHONY: all build run run-live run-web run-unified test clean fmt vet

BINARY_NAME=signal-sentry

all: build

build:
	go build -o $(BINARY_NAME) .

# Default legacy run
run:
	go run .

# Run Interactive TUI
run-live:
	go run . -live

# Run Web Server Only
run-web:
	go run . -web

# Run Unified Mode (TUI + Web)
run-unified:
	go run . -live -web

test:
	go test -v ./...

clean:
	rm -f $(BINARY_NAME)
	rm -f stats.log
	rm -f signal-data.json
	rm -f signal-data.csv

fmt:
	go fmt ./...

vet:
	go vet ./...