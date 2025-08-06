#!/bin/bash

# Fiber Challenge 1: Basic Routing - Test Runner
set -e

echo "🚀 Fiber Challenge 1: Basic Routing Test Runner"
echo "================================================"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first."
    exit 1
fi

# Check if solution file exists
if [ ! -f "solution-template.go" ]; then
    echo "❌ solution-template.go not found!"
    echo "Please make sure you're in the challenge directory."
    exit 1
fi

# Get username for submission
read -p "Enter your GitHub username (for submission tracking): " username

if [ -z "$username" ]; then
    echo "❌ Username is required for submission tracking."
    exit 1
fi

echo "👤 Testing solution for: $username"
echo ""

# Create temporary directory for testing
temp_dir=$(mktemp -d)
echo "📁 Created temporary test environment: $temp_dir"

# Copy files to temp directory
cp -r . "$temp_dir/"
cd "$temp_dir"

# Initialize go mod if needed
if [ ! -f "go.sum" ]; then
    echo "📦 Installing dependencies..."
    go mod tidy
fi

# Run the tests
echo "🧪 Running tests..."
echo ""

if go test -v; then
    echo ""
    echo "✅ All tests passed! Great job!"
    echo ""
    
    # Create submission directory if it doesn't exist
    submission_dir="../submissions/$username"
    mkdir -p "$submission_dir"
    
    # Copy solution to submissions
    cp solution-template.go "$submission_dir/solution.go"
    
    echo "💾 Solution saved to submissions/$username/solution.go"
    echo ""
    echo "🎉 Challenge completed successfully!"
    echo "Ready to move on to Challenge 2: Middleware"
    
else
    echo ""
    echo "❌ Some tests failed. Please review your implementation and try again."
    echo ""
    echo "💡 Hints:"
    echo "  - Check hints.md for implementation guidance"
    echo "  - Ensure all TODO sections are implemented"
    echo "  - Verify HTTP status codes and JSON responses"
    echo "  - Make sure routes are defined correctly"
    echo ""
fi

# Cleanup
cd - > /dev/null
rm -rf "$temp_dir"

echo "🧹 Cleaned up temporary files"