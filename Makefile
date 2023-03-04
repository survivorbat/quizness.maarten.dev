MAKEFLAGS := --no-print-directory --silent
.PHONY: server

default: help

help:
	@echo "Please use 'make <target>' where <target> is one of"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z\._-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

t: test
test: fmt ## Run unit tests, alias: t
	cd server && go test ./... -timeout=60s -parallel=10 --cover

fmt: ## Format go code
	cd server && go mod tidy
	cd server && go fmt ./...

dr: ## Run docker containers
	docker-compose up -d

server: ## Run the server
	cd server/cmd/qq && \
		DB_CONNECTION_STRING="host=localhost user=postgres password=rootroot dbname=qq" \
		AUTH_CREDENTIALS_PATH="../../credentials.json" \
		JWT_SECRET="abc" \
		AUTH_REDIRECT_URL="http://localhost:8000/auth" go run .
