default:
    @just --list --unsorted

# Run a specific day with real input (e.g., just run day01)
run day:
    @echo "Running {{day}}"
    @go run ./{{day}} -i ./{{day}}/input.txt

# Run tests for a specific day (e.g., just test day01), or all tests if no day specified
test day="":
    @if [ -z "{{day}}" ]; then \
        echo "Running all tests..."; \
        go test ./...; \
    else \
        echo "Running {{day}} Example Test"; \
        go test ./{{day}}; \
    fi

# Format code
[group('dev')]
fmt:
    @echo "Formatting code..."
    @go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.5.0 fmt ./...

# Run linters with auto-fix
[group('dev')]
lint:
    @echo "Running linters with auto-fix..."
    @go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.5.0 run --fix ./...

# Modernize code to latest Go idioms
[group('dev')]
modernize:
    @echo "Modernizing code to latest Go idioms..."
    @go run golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest -fix -test ./...

# Pre-commit: format, lint, modernize, and test
[group('dev')]
pre-commit: fmt lint modernize test
    @echo "Pre-commit checks completed."
