.PHONY: default run build test docs clean docker-build docker-run docker-stop docker-clean

# Variables
APP_NAME=opportunities
ENTRY_POINT=cmd/server/main.go
DOCKER_IMAGE=opportunities-api

# Default task: generates documentation and runs the application
default: run-with-docs

# --- Local Development ---

# Runs the application without regenerating swagger documentation
run:
	@go run $(ENTRY_POINT)

# Generates swagger documentation and runs the application
run-with-docs:
	@swag init -g $(ENTRY_POINT) --parseInternal
	@go run $(ENTRY_POINT)

# Optimized build: generates the binary in the root directory
build:
	@swag init -g $(ENTRY_POINT) --parseInternal
	@go build -o $(APP_NAME) $(ENTRY_POINT)

# Runs tests in all packages recursively
test:
	@go test -v ./internal/... ./config/...

# Generates Swagger documentation only
docs:
	@swag init -g $(ENTRY_POINT) --parseInternal

# Removes binaries and temporary folders
clean:
	@rm -f $(APP_NAME)
	@rm -rf ./docs
	@rm -rf ./db/*.db

# --- Docker ---

# Builds the Docker image using the multi-stage Dockerfile
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE) .

# Runs the container with port mapping and volume persistence for SQLite
docker-run:
	@echo "Running container on port 8080..."
	@docker run --name $(APP_NAME) -p 8080:8080 -v $(shell pwd)/db:/root/db $(DOCKER_IMAGE)

# Stops and removes the application container
docker-stop:
	@echo "Stopping and removing container..."
	@docker stop $(APP_NAME) || true
	@docker rm $(APP_NAME) || true

# Removes the Docker image and cleans up associated containers
docker-clean: docker-stop
	@echo "Cleaning up Docker images..."
	@docker rmi $(DOCKER_IMAGE) || true