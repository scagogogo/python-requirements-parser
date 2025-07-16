# Python Requirements Parser Makefile

.PHONY: help test fmt vet lint build clean docs docs-dev docs-build docs-deploy benchmark examples install-tools

# Default target
help: ## Show this help message
	@echo "Python Requirements Parser - Available commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

# Development commands
fmt: ## Format Go code
	@echo "📝 Formatting Go code..."
	@gofmt -w .
	@echo "✅ Code formatted"

vet: ## Run go vet
	@echo "🔍 Running go vet..."
	@go vet ./...
	@echo "✅ go vet passed"

lint: fmt vet ## Run all linting tools
	@echo "🔍 Running all linting tools..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not installed, skipping"; \
	fi
	@echo "✅ Linting completed"

test: ## Run all tests
	@echo "🧪 Running tests..."
	@go test ./...
	@echo "✅ All tests passed"

test-verbose: ## Run tests with verbose output
	@echo "🧪 Running tests (verbose)..."
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "🧪 Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report generated: coverage.html"

benchmark: ## Run benchmark tests
	@echo "⚡ Running benchmarks..."
	@go test -bench=. -benchmem ./...

# Build commands
build: ## Build the project
	@echo "🔨 Building project..."
	@go build ./...
	@echo "✅ Build completed"

clean: ## Clean build artifacts
	@echo "🧹 Cleaning build artifacts..."
	@go clean ./...
	@rm -f coverage.out coverage.html
	@echo "✅ Clean completed"

# Documentation commands
docs-dev: ## Start documentation development server
	@echo "📖 Starting documentation development server..."
	@cd docs && npm run dev

docs-build: ## Build documentation
	@echo "🔨 Building documentation..."
	@cd docs && npm run build
	@echo "✅ Documentation built"

docs-deploy: ## Deploy documentation to GitHub Pages
	@echo "🚀 Deploying documentation..."
	@cd docs && ./deploy.sh

docs-install: ## Install documentation dependencies
	@echo "📦 Installing documentation dependencies..."
	@cd docs && npm install
	@echo "✅ Documentation dependencies installed"

# Example commands
examples: ## Run all examples
	@echo "🎯 Running examples..."
	@for dir in examples/*/; do \
		if [ -f "$$dir/main.go" ]; then \
			echo "Running $$dir..."; \
			cd "$$dir" && go run main.go && cd ../..; \
		fi \
	done
	@echo "✅ All examples completed"

example-basic: ## Run basic usage example
	@echo "🎯 Running basic usage example..."
	@cd examples/01-basic-usage && go run main.go

example-editor-v2: ## Run version editor V2 example
	@echo "🎯 Running version editor V2 example..."
	@cd examples/07-version-editor-v2 && go run main.go

# Development tools
install-tools: ## Install development tools
	@echo "🔧 Installing development tools..."
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✅ Development tools installed"

# Git hooks
install-hooks: ## Install git hooks
	@echo "🪝 Installing git hooks..."
	@cp scripts/pre-commit.sh .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "✅ Git hooks installed"

# CI/CD simulation
ci: fmt vet test ## Run CI checks locally
	@echo "🤖 Running CI checks..."
	@echo "✅ All CI checks passed"

ci-full: ## Run full CI simulation (like GitHub Actions)
	@echo "🚀 Running full CI simulation..."
	@./scripts/simulate-ci.sh

# Release preparation
pre-release: clean fmt vet test benchmark docs-build ## Prepare for release
	@echo "🚀 Preparing for release..."
	@echo "✅ Release preparation completed"

# Quick development setup
setup: install-tools docs-install install-hooks ## Setup development environment
	@echo "🛠️  Setting up development environment..."
	@echo "✅ Development environment ready"

# Performance testing
perf: ## Run performance tests
	@echo "⚡ Running performance tests..."
	@go test -bench=. -benchmem -cpuprofile=cpu.prof -memprofile=mem.prof ./pkg/...
	@echo "📊 Performance profiles generated: cpu.prof, mem.prof"

# Security check
security: ## Run security checks
	@echo "🔒 Running security checks..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "⚠️  gosec not installed, install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

# Dependency management
deps-update: ## Update dependencies
	@echo "📦 Updating dependencies..."
	@go get -u ./...
	@go mod tidy
	@echo "✅ Dependencies updated"

deps-verify: ## Verify dependencies
	@echo "🔍 Verifying dependencies..."
	@go mod verify
	@echo "✅ Dependencies verified"

# All-in-one commands
check: fmt vet test ## Run all checks
	@echo "✅ All checks completed"

dev: fmt test ## Quick development check
	@echo "✅ Development check completed"
