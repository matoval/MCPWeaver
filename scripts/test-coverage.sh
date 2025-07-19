#!/bin/bash

# Test Coverage Script for MCPWeaver
# Runs comprehensive tests and generates coverage reports

set -e

echo "🧪 Running MCPWeaver Test Suite..."
echo "=================================="

# Clean up previous coverage files
rm -f coverage.out coverage.html

# Run tests with coverage
echo "Running unit tests..."
go test -v -coverprofile=coverage.out ./internal/... ./tests/unit/...

# Generate coverage report
echo ""
echo "📊 Coverage Summary:"
echo "==================="
go tool cover -func=coverage.out

# Generate HTML coverage report
echo ""
echo "📈 Generating HTML coverage report..."
go tool cover -html=coverage.out -o coverage.html
echo "Coverage report generated: coverage.html"

# Calculate overall coverage percentage
COVERAGE=$(go tool cover -func=coverage.out | tail -1 | awk '{print $3}' | sed 's/%//')
echo ""
echo "🎯 Overall Coverage: ${COVERAGE}%"

# Coverage targets
TARGET=90
if (( $(echo "$COVERAGE >= $TARGET" | bc -l) )); then
    echo "✅ Coverage target achieved! ($COVERAGE% >= $TARGET%)"
    exit 0
else
    echo "⚠️  Coverage below target: $COVERAGE% < $TARGET%"
    echo "   Needs improvement in the following areas:"
    echo ""
    
    # Show low coverage areas
    echo "📉 Components with low coverage:"
    go tool cover -func=coverage.out | awk '$3+0 < 50 { printf "   %-30s %s\n", $1, $3 }'
    
    exit 1
fi