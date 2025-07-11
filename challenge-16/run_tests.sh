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

# Copy the participant's solution and the test file to the temporary directory
cp "$SUBMISSION_FILE" "solution-template_test.go" "$TEMP_DIR/"

echo "Running tests for user '$USERNAME'..."

# Navigate to the temporary directory
pushd "$TEMP_DIR" > /dev/null

# Initialize a new Go module in the temporary directory
go mod init "challenge" || {
  echo "Failed to initialize Go module."
  popd > /dev/null
  rm -rf "$TEMP_DIR"
  exit 1
}

# Run the tests
go test -v

TEST_EXIT_CODE=$?

# Return to the original directory
popd > /dev/null

# Clean up the temporary directory
rm -rf "$TEMP_DIR"

exit $TEST_EXIT_CODE 