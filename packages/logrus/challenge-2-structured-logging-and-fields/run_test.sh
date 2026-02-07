#!/bin/bash

# Script to run tests for a participant's submission

# Function to display usage
usage() {
    echo "Usage: $0"
    exit 1
}

# Verify that we are in a challenge directory
if [ ! -f "solution-template_test.go" ]; then
    echo "Error: solution-template_test.go not found. Please run this script from a challenge directory"
    exit 1
fi

# Prompt for GitHub username
read -p "Enter your GitHub username: " USERNAME

SUBMISSION_DIR="submissions/$USERNAME"
SUBMISSION_FILE="$SUBMISSION_DIR/solution.go"

# Check if the submission file exists
if [ ! -f "$SUBMISSION_FILE" ]; then
    echo "Error: Solution file '$SUBMISSION_FILE' not found"
    echo "Note: Please ensure your solution is named 'solution.go' and placed in a 'submissions/<username>/' directory"
    exit 1
fi

# Create a temporary directory to avoid modifying the original files
TEMP_DIR=$(mktemp -d)

# Copy the participant's solution, test file, and go.mod/go.sum to the temporary directory
cp "$SUBMISSION_FILE" "solution-template_test.go" "go.mod" "go.sum" "$TEMP_DIR/" 2>/dev/null

# The test file expects to be in the `main` package alongside the functions it's testing
mv "$TEMP_DIR/solution.go" "$TEMP_DIR/solution-template.go"

echo "Running tests for user '$USERNAME'..."

# Navigate to the temporary directory
pushd "$TEMP_DIR" > /dev/null || {
    echo "Failed to navigate to temporary directory."
    rm -rf "$TEMP_DIR"
    exit 1
}

# Tidy up dependencies to ensure everything is consistent
echo "Tidying dependencies..."
go mod tidy || {
    echo "Failed to tidy dependencies."
    popd > /dev/null || {
        echo "Failed to return to original directory."
        rm -rf "$TEMP_DIR"
        exit 1
    }
    rm -rf "$TEMP_DIR"
    exit 1
}

# Run the tests with verbosity and coverage
echo "Executing tests..."
go test -v -cover

TEST_EXIT_CODE=$?

# Return to the original directory
popd > /dev/null || {
    echo "Failed to return to original directory."
    rm -rf "$TEMP_DIR"
    exit 1
}

# Clean up the temporary directory
rm -rf "$TEMP_DIR"

if [ $TEST_EXIT_CODE -eq 0 ]; then
    echo "All tests passed!"
else
    echo "Some tests failed."
fi

exit $TEST_EXIT_CODE