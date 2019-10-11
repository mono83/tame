# Makefile configuration
.DEFAULT_GOAL := help
.PHONY: help deps travis build release local

deps: ## Download dependencies
	@echo "Downloading dependencies"
	go get github.com/spf13/cobra
	go get github.com/mono83/xray
	go get github.com/PuerkitoBio/goquery
	go get github.com/dsnet/compress
	go get github.com/stretchr/testify/assert

travis: deps test local ## Runs travis tasks
	@echo "Running external tests"
	@release/tame httpbin

local: ## Builds local binary
	@mkdir -p release
	@echo && echo "Building binaries"
	go build -o release/tame tame/tame.go

build: local ## Builds binaries
	GOOS="linux" GOARCH="amd64" go build -o release/tame-linux64 tame/tame.go
	GOOS="darwin" GOARCH="amd64" go build -o release/tame-darwin64 tame/tame.go

release: deps build ## Builds release

test:
	@echo "Running tests"
	@go test ./...

help:
	@grep --extended-regexp '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
