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

# Rename solution.go to solution-template.go for the test
mv "$TEMP_DIR/solution.go" "$TEMP_DIR/solution-template.go"

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

# Run the tests
go test -v

TEST_EXIT_CODE=$?

# Return to the original directory
popd > /dev/null

# Clean up the temporary directory
rm -rf "$TEMP_DIR"

exit $TEST_EXIT_CODE

