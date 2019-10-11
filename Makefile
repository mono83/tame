# Makefile configuration
.DEFAULT_GOAL := help
.PHONY: help deps travis build build-all release local test

deps: ## Download dependencies
	@echo "Downloading dependencies"
	go get github.com/mono83/xray
	go get github.com/fatih/color
	go get github.com/spf13/cobra
	go get github.com/PuerkitoBio/goquery
	go get github.com/dsnet/compress

	@echo "Downloading test dependencies"
	go get github.com/stretchr/testify/assert

docker: ## Builds application docker image
	docker build -t tame .

travis: deps test local ## Runs travis tasks
	@echo "Running external tests"
	@release/tame httpbin

local: ## Builds local binary
	@mkdir -p release
	@echo "Building binaries"
	CGO_ENABLED=0 go build -o release/tame app/tame.go

build: deps test local ## Builds native binary

build-all: build ## Builds all binaries
	GOOS="linux" GOARCH="amd64" go build -o release/tame-linux64 app/tame.go
	GOOS="darwin" GOARCH="amd64" go build -o release/tame-darwin64 app/tame.go

release: deps build-all ## Builds release

test:
	@echo "Running tests"
	@go test ./...

help:
	@grep --extended-regexp '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
