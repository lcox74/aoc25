default:
    @just --list --unsorted

# Run all tests
test:
    @echo "Running all tests..."
    @go test ./...

# Run Day 1 with real input
[group('day01')]
day01:
    @echo "Running Day 1"
    @go run ./day01 -i ./day01/input.txt

# Run Day 1 with example input
[group('day01')]
day01-test:
    @echo "Running Day 1 Example Test"
    @go test ./day01 

# Run Day 2 with real input
[group('day02')]
day02:
    @echo "Running Day 2"
    @go run ./day02 -i ./day02/input.txt

# Run Day 2 with example input
[group('day02')]
day02-test:
    @echo "Running Day 2 Example Test"
    @go test ./day02 

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
