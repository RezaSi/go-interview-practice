#!/bin/bash
# Run this challenge's tests, either against the shipped template or a submission.
set -e
if [ ! -f "solution-template_test.go" ]; then
    echo "Error: run this from a challenge directory."
    exit 1
fi
SOLUTION="solution-template.go"
if [ -n "$1" ]; then
    SOLUTION="submissions/$1/solution-template.go"
    [ -f "$SOLUTION" ] || { echo "Error: '$SOLUTION' not found."; exit 1; }
    echo "Running tests for submission: $1"
else
    echo "Running tests against solution-template.go (pass a username to test a submission)."
fi
TEMP_DIR=$(mktemp -d)
trap 'rm -rf "$TEMP_DIR"' EXIT
cp go.mod solution-template_test.go "$TEMP_DIR/"
cp "$SOLUTION" "$TEMP_DIR/solution-template.go"
cd "$TEMP_DIR"
go test -v ./...
