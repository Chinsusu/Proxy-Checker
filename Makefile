.PHONY: build test lint clean frontend

# Build production binary with embedded frontend
build: frontend
	go build -o bin/proxy-checker .

# Build frontend static files
frontend:
	cd frontend && npm install && npm run build

test:
	go test ./...

lint:
	golangci-lint run

clean:
	rm -rf bin
	rm -rf frontend/dist
	rm -rf frontend/node_modules
