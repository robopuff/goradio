.DEFAULT_GOAL := help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "[32m%-40s[0m %s\n", $$1, $$2}'

run: ## Run code without building executable
	go run main.go

.PHONY: build
build: ## Build an executable
	go build -o build/goradio

move-executable: ## Move executable to ~/.local/bin
	mv build/goradio ~/.local/bin/goradio

install: build move-executable ## Build and install executable
