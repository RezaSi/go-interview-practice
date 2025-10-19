package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		input := scanner.Text()
		output := ReverseString(input)
		fmt.Println(output)
	}
}

func ReverseString(s string) string {
	runes := []rune(s)
	builder := strings.Builder{}

	for i := len(runes) - 1; i >= 0; i-- {
		builder.WriteRune(runes[i])
	}

	return builder.String()
}
