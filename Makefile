################################################################################
# Code Quality Targets
################################################################################
fmt:
	@go run mvdan.cc/gofumpt@latest -l -w .
.PHONY: fmt

lint:
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.2 run --fix --timeout 5m
.PHONY: lint

