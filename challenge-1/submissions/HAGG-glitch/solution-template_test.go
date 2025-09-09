package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestSum(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Positive numbers", "2, 3", "5"},
		{"Zero values", "0, 0", "0"},
		{"Negative numbers", "-2, -3", "-5"},
		{"Mixed signs", "-5, 10", "5"},
		{"Large numbers", "1000000000, 1000000000", "2000000000"},
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
