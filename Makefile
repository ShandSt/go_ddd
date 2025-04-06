.PHONY: build test run docker-build docker-run lint

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=main
BINARY_UNIX=$(BINARY_NAME)_unix

# Docker parameters
DOCKER_IMAGE=ddd
DOCKER_TAG=latest

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/main.go

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/main.go
	./$(BINARY_NAME)

docker-build:
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run:
	docker run -p 8080:8080 $(DOCKER_IMAGE):$(DOCKER_TAG)

lint:
	golangci-lint run

deps:
	$(GOGET) -v -t -d ./...

swagger:
	swag init -g cmd/main.go -o docs 