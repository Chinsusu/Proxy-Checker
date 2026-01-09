.PHONY: build test lint clean dev

build:
	wails build

dev:
	wails dev

test:
	go test ./...

lint:
	golangci-lint run

clean:
	rm -rf build/bin
	rm -rf frontend/dist
