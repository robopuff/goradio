.DEFAULT_GOAL := help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "[32m%-40s[0m %s\n", $$1, $$2}'

run: ## Run debug code without building executable
	go run main.go -d

.PHONY: build
build: ## Build an executable
	go build -o build/goradio

install: ## Install executable
	 go install

clear-build: ## Clear content of build directory
	rm build/goradio*

release: ## Build for windows, linux and darwin
	for os in darwin linux windows; do \
		env GOOS=$$os GOARCH=amd64 go build -o build/goradio-$$os-amd64; \
	done; \
	mv build/goradio-windows-amd64 build/goradio-windows-amd64.exe