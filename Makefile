MAKEFLAGS := --no-print-directory --silent
.PHONY: server

default: help

help:
	@echo "Please use 'make <target>' where <target> is one of"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z\._-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)


fmt: ## Format go code
	cd server && go mod tidy
	cd server && go fmt ./...
	cd server && swag fmt
	cd web && npm run fmt

dr: ## Run docker containers
	POSTGRES_PASSWORD="rootroot" docker-compose -f docker-compose.yaml up -d

server: ## Run the server
	cd server && swag init --output swagger -g cmd/qq/main.go
	cd server/cmd/qq && go run .

ui: ## Run the ui
	cd web && npm start

install: # Install dependencies and tools
	cd server && go get
	cd web && npm i

	go install github.com/swaggo/swag/cmd/swag@latest
