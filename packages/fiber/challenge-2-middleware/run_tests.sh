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

# Copy the participant's solution, test file, and go.mod to the temporary directory
cp "$SUBMISSION_FILE" "solution-template_test.go" "$TEMP_DIR/"

# Copy go.mod if it exists
if [ -f "go.mod" ]; then
    cp "go.mod" "$TEMP_DIR/"
fi

# Copy go.sum if it exists
if [ -f "go.sum" ]; then
    cp "go.sum" "$TEMP_DIR/"
fi

# Change to the temporary directory
cd "$TEMP_DIR" || exit 1

# Rename the solution file to match the expected name
mv "solution.go" "solution-template.go"

echo "Running tests for $USERNAME's solution..."

# Run the tests
go test -v

# Capture the exit code
EXIT_CODE=$?

# Return to the original directory
cd - > /dev/null

# Clean up the temporary directory
rm -rf "$TEMP_DIR"

# Check if tests passed
if [ $EXIT_CODE -eq 0 ]; then
    echo "✅ All tests passed! Great job, $USERNAME!"
    echo "Your solution is ready for submission."
else
    echo "❌ Some tests failed. Please review your implementation and try again."
    exit 1
fi