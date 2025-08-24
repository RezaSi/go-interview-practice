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
SUBMISSION_FILE="$SUBMISSION_DIR/solution.go"

# Check if the submission file exists
if [ ! -f "$SUBMISSION_FILE" ]; then
    echo "Error: Solution file '$SUBMISSION_FILE' not found."
    echo "Note: Package challenges use 'solution.go' instead of 'solution-template.go'"
    exit 1
fi

# Create a temporary directory to avoid modifying the original files
TEMP_DIR=$(mktemp -d)
cp *.go "$TEMP_DIR/"
cp go.mod "$TEMP_DIR/"
if [ -f "go.sum" ]; then
    cp go.sum "$TEMP_DIR/"
fi

# Replace the template file with the submission
cp "$SUBMISSION_FILE" "$TEMP_DIR/solution-template.go"

# Navigate to the temporary directory
cd "$TEMP_DIR"

echo "Running tests for $USERNAME's submission..."
echo "=========================================="

# Initialize Go module and download dependencies
go mod tidy

# Run the tests
TEST_OUTPUT=$(go test -v 2>&1)
TEST_EXIT_CODE=$?

echo "$TEST_OUTPUT"

if [ $TEST_EXIT_CODE -eq 0 ]; then
    echo "=========================================="
    echo "✅ All tests passed! Great job, $USERNAME!"
    
    # Count passed tests
    PASSED_TESTS=$(echo "$TEST_OUTPUT" | grep -c "PASS: Test")
    echo "📊 Passed tests: $PASSED_TESTS"
    
    # Update scoreboard
    cd - > /dev/null
    python3 ../../scripts/update_scoreboard.py "$USERNAME" "challenge-5-generics" $PASSED_TESTS
    
else
    echo "=========================================="
    echo "❌ Some tests failed. Keep working on it!"
    
    # Show failed tests
    FAILED_TESTS=$(echo "$TEST_OUTPUT" | grep -c "FAIL: Test")
    PASSED_TESTS=$(echo "$TEST_OUTPUT" | grep -c "PASS: Test")
    
    echo "📊 Test Results:"
    echo "   ✅ Passed: $PASSED_TESTS"
    echo "   ❌ Failed: $FAILED_TESTS"
    
    cd - > /dev/null
fi

# Cleanup
rm -rf "$TEMP_DIR"

exit $TEST_EXIT_CODE 