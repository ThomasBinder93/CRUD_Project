# Makefile for CRUD Project

.PHONY: help build run test clean lint fmt docker-build docker-run migrate

help:
	@echo "Available targets:"
	@echo "  make build        - Build the application"
	@echo "  make run          - Run the application"
	@echo "  make test         - Run tests with coverage"
	@echo "  make test-unit    - Run unit tests only"
	@echo "  make test-e2e     - Run end-to-end tests (requires Node.js and Playwright)"
	@echo "  make test-selenium - Run Selenium browser tests (requires Python and Selenium)"
	@echo "  make lint         - Run linter"
	@echo "  make fmt          - Format code"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-run   - Run Docker container"
	@echo "  make docker-stop  - Stop Docker container"

build:
	go build -o crud-app .

run: build
	./crud-app

test:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-unit:
	go test -v ./database ./handlers ./models ./config

test-e2e:
	npm test

test-selenium:
	python -m unittest discover -s tests/selenium -p "test_*.py" -v

lint:
	golangci-lint run ./...

fmt:
	go fmt ./...
	goimports -w .

clean:
	rm -f crud-app crud.db coverage.out coverage.html
	go clean

docker-build:
	docker build -t crud-demo:latest .

docker-run: docker-build
	docker run -p 8080:8080 -v $(PWD):/app crud-demo:latest

docker-stop:
	docker stop crud-demo || true

.DEFAULT_GOAL := help
