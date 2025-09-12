package main

import (
	"bufio"
	"fmt"
	"os"
)

func ReverseString(s string) string {
	runes := []rune(s)

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]

	}

	return string(runes)
}

func main() {

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n') // read until newline

	// Remove the trailing newline (if present)
	if len(input) > 0 && input[len(input)-1] == '\n' {
		input = input[:len(input)-1]
	}
	if len(input) > 0 && input[len(input)-1] == '\r' { // handle Windows \r\n
		input = input[:len(input)-1]
	}

	result := ReverseString(input)
	fmt.Print(result)
}
