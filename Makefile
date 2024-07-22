# make file
# Makefile for a Go project

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=go-s3

# Default target: build the project
all: build

# Build the project
build:
	$(GOBUILD) -o $(BINARY_NAME) -v cmd/main.go

# Run tests 
test:
	$(GOTEST) -v ./... -args -access=${ACCESS_KEY} -secret=${SECRET_KEY} -bucket=${BUCKET} -region=${REGION}

# Clean up build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Install dependencies
deps:
	$(GOGET) -v ./...

# Format code
fmt:
	$(GOCMD) fmt ./...

# Run the application
run:
	./$(BINARY_NAME) -access=${ACCESS_KEY} -secret=${SECRET_KEY} -bucket=${BUCKET} -region=${REGION}

.PHONY: all build test clean deps fmt run
