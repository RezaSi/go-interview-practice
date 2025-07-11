#!/bin/bash

# Script to run tests for a participant's submission

# Function to display usage
usage() {
    echo "Usage: $0"
    exit 1
}

# Verify that we are in a challenge directory
if [ ! -f "solution-template_test.go" ]; then
    echo "Error: solution-template_test.go not found. Please run this script from a challenge directory."
    exit 1
fi

# Prompt for GitHub username
read -p "Enter your GitHub username: " USERNAME

SUBMISSION_DIR="submissions/$USERNAME"
SUBMISSION_FILE="$SUBMISSION_DIR/solution-template.go"

# Check if the submission file exists
if [ ! -f "$SUBMISSION_FILE" ]; then
    echo "Error: Solution file '$SUBMISSION_FILE' not found."
    exit 1
fi

# Create a temporary directory to avoid modifying the original files
TEMP_DIR=$(mktemp -d)

# Copy the participant's solution, test file, and go.mod to the temporary directory
cp "$SUBMISSION_FILE" "solution-template_test.go" "$TEMP_DIR/"

# Copy go.mod if it exists
if [ -f "go.mod" ]; then
    cp "go.mod" "$TEMP_DIR/"
fi

echo "Running tests for user '$USERNAME'..."

# Navigate to the temporary directory
pushd "$TEMP_DIR" > /dev/null

# If go.mod exists, use it; otherwise initialize a new module
if [ -f "go.mod" ]; then
    echo "Using existing go.mod file"
    # Update module name to avoid conflicts (macOS compatible)
    sed -i '' 's/^module .*/module challenge/' go.mod
    # Download dependencies
    go mod tidy || {
        echo "Failed to download dependencies."
        popd > /dev/null
        rm -rf "$TEMP_DIR"
        exit 1
    }
else
    # Initialize a new Go module in the temporary directory
    go mod init "challenge" || {
        echo "Failed to initialize Go module."
        popd > /dev/null
        rm -rf "$TEMP_DIR"
        exit 1
    }
fi

echo "Running basic functionality tests..."
# Run the basic tests
go test -v -timeout=30s

TEST_EXIT_CODE=$?

if [ $TEST_EXIT_CODE -eq 0 ]; then
    echo ""
    echo "✓ Basic tests passed! Running performance benchmarks..."
    
    # Run benchmarks
    go test -v -bench=. -benchtime=1s -timeout=60s
    
    BENCH_EXIT_CODE=$?
    
    if [ $BENCH_EXIT_CODE -eq 0 ]; then
        echo ""
        echo "✓ Benchmarks completed! Running race condition tests..."
        
        # Run race condition tests
        go test -v -race -timeout=30s
        
        RACE_EXIT_CODE=$?
        
        if [ $RACE_EXIT_CODE -eq 0 ]; then
            echo ""
            echo "🎉 All tests passed! Your rate limiter implementation is working correctly."
            echo ""
            echo "Summary:"
            echo "✓ Basic functionality tests: PASSED"
            echo "✓ Performance benchmarks: COMPLETED"
            echo "✓ Race condition tests: PASSED"
            echo ""
            echo "Your solution is ready for submission!"
        else
            echo ""
            echo "❌ Race condition tests failed. Please fix thread safety issues."
            TEST_EXIT_CODE=$RACE_EXIT_CODE
        fi
    else
        echo ""
        echo "⚠️  Benchmarks had issues, but basic functionality works."
        echo "Consider optimizing your implementation for better performance."
    fi
else
    echo ""
    echo "❌ Basic tests failed. Please fix the implementation and try again."
fi

# Return to the original directory
popd > /dev/null

# Clean up the temporary directory
rm -rf "$TEMP_DIR"

exit $TEST_EXIT_CODE 