package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestReverseString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Simple word", "hello", "olleh"},
		{"Sentence with spaces", "Go is fun!", "!nuf si oG"},
		{"Empty string", "", ""},
		{"Palindrome", "madam", "madam"},
		{"Special characters", "12345!@#$%", "%$#@!54321"},
		{"Mixed case", "GoLang", "gnaLoG"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("go", "run", "submissions/tonybui98/solution-template.go")
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
