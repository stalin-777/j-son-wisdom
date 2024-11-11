#
# Makefile for Application
#

LINTER_VER ?= v1.61.0-alpine

## lint: Running linting.
lint:
	docker run --rm \
		-v $(shell pwd):/app \
		-v ~/.cache/golangci-lint/:/root/.cache \
		-v ${shell go env GOPATH}/pkg:/go/pkg \
		-w /app \
		golangci/golangci-lint:$(LINTER_VER) \
		golangci-lint run ./...

## test: Running unit tests.
test:
	go test -race --short ./...

## coverage: Running uint tests with coverage.
test.coverage:
	go test -race -covermode=atomic -coverprofile=coverage.out ./...

start:
	docker-compose up -d --build --force-recreate
# Removing containers with forced stopping
down:
	docker-compose down
# Stopping containers
stop:
	docker-compose stop
# Recreating containers
recreate:
	docker-compose down
	docker-compose up -d --build --force-recreate
