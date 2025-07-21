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
echo "Using temporary directory: $TEMP_DIR"

# Copy all files to temp directory
cp -r . "$TEMP_DIR/"
cd "$TEMP_DIR"

# Replace solution-template.go with the submission
if [ -f "solution-template.go" ]; then
    rm solution-template.go
fi
cp "$SUBMISSION_FILE" solution-template.go

echo "Testing submission for $USERNAME..."
echo "================================"

# Initialize go module if needed
if [ ! -f "go.mod" ]; then
    go mod init cobra-challenge-2
    go mod tidy
fi

# Download dependencies
echo "Downloading dependencies..."
go mod tidy

# Run the tests
echo "Running tests..."
echo "=================="

# Run basic compilation test
echo "1. Testing compilation..."
if ! go build -o filecli .; then
    echo "❌ Compilation failed"
    cd - > /dev/null
    rm -rf "$TEMP_DIR"
    exit 1
fi
echo "✅ Compilation successful"

# Run unit tests
echo ""
echo "2. Running unit tests..."
TEST_OUTPUT=$(go test -v 2>&1)
TEST_EXIT_CODE=$?

echo "$TEST_OUTPUT"

if [ $TEST_EXIT_CODE -eq 0 ]; then
    echo "✅ All unit tests passed"
else
    echo "❌ Some unit tests failed"
fi

# Run functional tests with the built binary
echo ""
echo "3. Running functional tests..."

# Test help command
echo "Testing help functionality..."
if ./filecli --help > /dev/null 2>&1; then
    echo "✅ Help command works"
else
    echo "❌ Help command failed"
fi

# Test global verbose flag
echo "Testing global verbose flag..."
if ./filecli --verbose list --help > /dev/null 2>&1; then
    echo "✅ Global verbose flag works"
else
    echo "❌ Global verbose flag failed"
fi

# Test list command with different formats
echo "Testing list command..."
if ./filecli list --help > /dev/null 2>&1; then
    echo "✅ List command help works"
else
    echo "❌ List command help failed"
fi

# Test copy command argument validation
echo "Testing copy command argument validation..."
if ./filecli copy 2>&1 | grep -q "accepts 2 arg"; then
    echo "✅ Copy command validates arguments correctly"
else
    echo "❌ Copy command argument validation failed"
fi

# Test delete command required flag
echo "Testing delete command required flag..."
if ./filecli delete test.txt 2>&1 | grep -q "required flag"; then
    echo "✅ Delete command requires --force flag"
else
    echo "❌ Delete command force flag validation failed"
fi

# Test create command required flag
echo "Testing create command required flag..."
if ./filecli create 2>&1 | grep -q "required flag"; then
    echo "✅ Create command requires --name flag"
else
    echo "❌ Create command name flag validation failed"
fi

# Create a test file for more advanced testing
echo "test content" > test.txt

# Test actual file operations if implemented
echo "Testing file operations..."

# Test list command on current directory
if ./filecli list . > /dev/null 2>&1; then
    echo "✅ List command can list current directory"
else
    echo "⚠️  List command implementation may be incomplete"
fi

# Test list command with JSON format
if ./filecli list --format json . > /dev/null 2>&1; then
    echo "✅ List command supports JSON format"
else
    echo "⚠️  List command JSON format may be incomplete"
fi

# Test create command with valid flags
if ./filecli create --name "testfile.txt" --size 10 > /dev/null 2>&1; then
    echo "✅ Create command works with valid flags"
else
    echo "⚠️  Create command implementation may be incomplete"
fi

# Test copy command with valid arguments
if ./filecli copy test.txt test_copy.txt > /dev/null 2>&1; then
    echo "✅ Copy command works with valid arguments"
else
    echo "⚠️  Copy command implementation may be incomplete"
fi

# Test delete command with force flag
if ./filecli delete test_copy.txt --force > /dev/null 2>&1; then
    echo "✅ Delete command works with --force flag"
else
    echo "⚠️  Delete command implementation may be incomplete"
fi

echo ""
echo "================================"
echo "Test Summary:"
echo "================================"

# Count results
PASSED=$(echo "$TEST_OUTPUT" | grep -c "PASS:")
FAILED=$(echo "$TEST_OUTPUT" | grep -c "FAIL:")
TOTAL_TESTS=$(echo "$TEST_OUTPUT" | grep -c "=== RUN")

echo "Unit Tests: $PASSED passed, $FAILED failed out of $TOTAL_TESTS total"

if [ $TEST_EXIT_CODE -eq 0 ]; then
    echo "Overall Result: ✅ PASSED"
    echo ""
    echo "🎉 Congratulations! Your solution passes all tests."
    echo "Your CLI demonstrates proper flag and argument handling!"
else
    echo "Overall Result: ❌ FAILED"
    echo ""
    echo "❗ Your solution needs some work. Please review the test output above."
    echo "Common issues to check:"
    echo "  - Are all commands properly implemented?"
    echo "  - Are flags correctly defined and validated?"
    echo "  - Are argument validators properly set?"
    echo "  - Do error messages match expectations?"
fi

echo ""
echo "📚 Next Steps:"
echo "  - Review the hints.md file for guidance"
echo "  - Check the learning.md for detailed explanations"
echo "  - Test your CLI manually with different flag combinations"
echo "  - Make sure all TODOs in the template are completed"

# Cleanup
cd - > /dev/null
rm -rf "$TEMP_DIR"

exit $TEST_EXIT_CODE 