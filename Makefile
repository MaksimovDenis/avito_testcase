# Makefile
.SILENT:

.DEFAULT_GOAL := help
SHELL := /bin/bash

.PHONY: help
help:
	@echo "Usage: make [target]"
	@echo "Targets:"
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  %-20s %s\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

## Build App 
build:
	go mod download && CGD_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/main.go

## Build and start the containers
up:
	docker-compose up --build 

## Stop and remove containers
down:
	docker-compose down

## Restart the containers
restart: down up

## Run tests
test:
	go test ./... -coverprofile cover.out

## Run tests coverage
test-coverage:
	go tool cover -func cover.out | grep total | awk '{print $3}'

## chmod DB Folders
chmod:
	sudo chmod -R 777 ./databaseRedis 
	sudo chmod -R 777 ./database 
	sudo chmod -R 777 ./prometheus 
	sudo chmod -R 777 ./grafana 

## Delete DB Folders
rm:
	sudo rm -rf ./database ./databaseRedis ./grafana ./prometheus

