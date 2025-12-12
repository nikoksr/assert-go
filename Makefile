################################################################################
# Code Quality Targets
################################################################################
fmt:
	@goimports -w .
	@golines --shorten-comments -m 120 -w .
	@gofumpt -w -l .
	@gci write -s standard -s default -s "prefix(github.com/nikoksr/assert-go)" .
.PHONY: fmt

lint:
	@golangci-lint run ./...
.PHONY: lint

################################################################################
# Test Targets
################################################################################
test:
	@go test -v -race -coverprofile=coverage.out ./...
.PHONY: test

test-debug:
	@go test -v -race -tags assertdebug ./...
.PHONY: test-debug

test-coverage:
	@go test -race -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
.PHONY: test-coverage

bench:
	@go test -bench=. -benchmem ./...
.PHONY: bench

bench-debug:
	@go test -bench=. -benchmem -tags assertdebug ./...
.PHONY: bench-debug

################################################################################
# Convenience Targets
################################################################################
check: fmt lint test
	@echo "✅ All checks passed!"
.PHONY: check

ci: lint test test-debug
	@echo "✅ CI checks passed!"
.PHONY: ci

