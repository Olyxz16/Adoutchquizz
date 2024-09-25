# Simple Makefile for a Go project

# Build the application
all: build

build:
	@templ generate
	@go build -o build/main.exe cmd/main.go

build-linux:
	@templ generate
	@set GOOS=linux && go build -o build/main cmd/main.go

build-all:
	@npm run build
	@templ generate
	@go build -o build/main.exe cmd/main.go

dev: watch-tailwind air
watch-tailwind:
	@npm run watch
air: 
	@air

# Run the application
run:
	@npm run build
	@templ generate
	@go run cmd/main.go

install:
	@-mkdir .\static\js
	@curl -L https://unpkg.com/htmx.org@latest/dist/htmx.min.js -o ./static/js/htmx.min.js
	@go install github.com/a-h/templ/cmd/templ@latest
	@go mod download
	@npm i

init-db:
	@createdb -U postgres Adoutchquizz

# Create DB container
docker-run:
	@if docker compose up 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Test the application
test:
	@echo "Testing..."
	@go test ./tests -v

.PHONY: all build run test clean
