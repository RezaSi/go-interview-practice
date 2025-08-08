package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestKnapsack(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "Sample test case 1",
			input: `4
2 3
3 4
4 5
5 6
10`,
			expected: "13",
		},
		{
			name: "Sample test case 2",
			input: `3
1 1
2 4
3 5
4`,
			expected: "6",
		},
		{
			name: "Single item",
			input: `1
5 10
5`,
			expected: "10",
		},
		{
			name: "Single item exceeds capacity",
			input: `1
10 20
5`,
			expected: "0",
		},
		{
			name: "Multiple items with optimal selection",
			input: `5
1 1
2 3
3 4
4 5
5 6
7`,
			expected: "9",
		},
		{
			name: "All items fit",
			input: `3
1 2
2 3
3 4
10`,
			expected: "9",
		},
		{
			name: "No items fit",
			input: `3
5 10
6 12
7 14
4`,
			expected: "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("go", "run", "solution-template.go")
			stdin := strings.NewReader(tt.input)
			var stdout, stderr bytes.Buffer
			cmd.Stdin = stdin
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			if err != nil {
				t.Fatalf("Error running the program: %v\nStderr: %s", err, stderr.String())
			}

			output := strings.TrimSpace(stdout.String())
			if output != tt.expected {
				t.Errorf("For input '%s', expected output '%s', got '%s'", tt.input, tt.expected, output)
			}
		})
	}
}
