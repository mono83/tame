# Makefile configuration
.DEFAULT_GOAL := help
.PHONY: help deps travis release

deps: ## Download dependencies
	@echo "Downloading dependencies"
	go get github.com/spf13/cobra
	go get github.com/mono83/xray
	go get github.com/PuerkitoBio/goquery
	go get github.com/stretchr/testify/assert

travis: deps ## Runs travis tasks
	@mkdir -p release
	go build -o release/tame tame.go
	@echo "Running tests"
	@release/tame httpbin

release: ## Builds release
	@mkdir -p release
	@echo && echo "Building binaries"
	GOOS="linux" GOARCH="amd64" go build -o release/tame-linux64 tame.go
	GOOS="darwin" GOARCH="amd64" go build -o release/tame-darwin64 tame.go
	go build -o release/tame tame.go

help:
	@grep --extended-regexp '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
