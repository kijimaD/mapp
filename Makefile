.DEFAULT_GOAL := help

.PHONY: lint
lint: ## Run lint
	docker run --rm -v ${PWD}:/app -w /app golangci/golangci-lint:v1.51.2 golangci-lint run -v
	docker run --rm -v ${PWD}:/app -w /app golang:1.20 go vet ./...

.PHONY: help
help: ## Show help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
