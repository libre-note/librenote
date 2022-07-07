#!/usr/bin/env bash

BUILD_VERSION := $(shell git describe --always --tags)
BUILD_TIME=$(shell date '+%Y%m%d-%H%M%S')
DOCKER_IMAGE_NAME="hrshadhin/librenote"
BINARY_NAME=librenote
BIN_OUT_DIR=bin

MIGRATION_PATH_PG="infrastructure/db/migrations/pgsql"
MIGRATION_PATH_MYSQL="infrastructure/db/migrations/mysql"
MIGRATION_PATH_SQLITE="infrastructure/db/migrations/sqlite"

export PATH=$(shell go env GOPATH)/bin:$(shell echo $$PATH)

.PHONY: all

all: build test.unit ## Build binary (with unit tests)

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

setup:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.44.0
	go install golang.org/x/tools/cmd/goimports@v0.1.9

lint: setup build ## Run lint checks
	$(shell go env GOPATH)/bin/golangci-lint run

fmt: setup ## Refactor go files with gofmt and goimports
	find . -name '*.go' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done

test-unit:  ## Run unit tests
	go test -v -coverprofile=coverage.txt -covermode=atomic -cover ./app/...

test-integration:  ## Run integration tests
	go test -v -tags=integration ./it -count=1

clean: ## Cleans output directory
	$(shell rm -rf $(BIN_OUT_DIR)/*)
	$(shell rm -rf ./*.db ./it/*.db coverage.txt)

build-deps:
	go mod vendor

build: clean build-deps ## Build binary
	go build -v -ldflags="-w -s -X librenote/app.Version=${BUILD_VERSION} -X librenote/app.BuildTime=${BUILD_TIME}" -o $(BIN_OUT_DIR)/$(BINARY_NAME)

run: build ## Build and run binary
	./$(BIN_OUT_DIR)/$(BINARY_NAME)

serve: build ## Run http server
	./$(BIN_OUT_DIR)/$(BINARY_NAME) serve

swagger: ## Creates swagger documentation as html file
	go install github.com/swaggo/swag/cmd/swag@v1.7.9-p1
	$(shell go env GOPATH)/bin/swag init -g _doc/api.go -o _doc

migrate-up-pgsql: build ## Run migration postgresql
	./$(BIN_OUT_DIR)/$(BINARY_NAME) migrate -p ${MIGRATION_PATH_PG} up

migrate-down-pgsql: build ## Revert migration postgresql
	./$(BIN_OUT_DIR)/$(BINARY_NAME) migrate -p ${MIGRATION_PATH_PG} down

migrate-up-mysql: build ## Run migration mysql
	./$(BIN_OUT_DIR)/$(BINARY_NAME) migrate -p ${MIGRATION_PATH_MYSQL} up

migrate-down-mysql: build ## Revert migration mysql
	./$(BIN_OUT_DIR)/$(BINARY_NAME) migrate -p ${MIGRATION_PATH_MYSQL} down

migrate-up-sqlite: build ## Run migration sqlite
	./$(BIN_OUT_DIR)/$(BINARY_NAME) migrate -p ${MIGRATION_PATH_SQLITE} up

migrate-down-sqlite: build ## Revert migration sqlite
	./$(BIN_OUT_DIR)/$(BINARY_NAME) migrate -p ${MIGRATION_PATH_SQLITE} down

docker-build: ## Build docker image
	docker build --build-arg BUILD_VERSION=${BUILD_VERSION} --build-arg BUILD_TIME=${BUILD_TIME} --tag ${DOCKER_IMAGE_NAME} .

docker-push: ## Push docker image
	docker push

docker-run: ## Run docker image with sqlite
	mkdir -p data
	sudo chown -R 1000:1000 data
	docker run --name librenote_server --rm -it -p 8000:8000 \
		-v $$(pwd)/config.yml:/app/config.yml \
		-v $$(pwd)/infrastructure/db/migrations/sqlite:/app/migrations \
		-v $$(pwd)/data:/persist \
		$(DOCKER_IMAGE_NAME):latest

docker-migrate: ## Run migrations inside dorker
	docker exec librenote_server /app/librenote migrate -p /app/migrations up
