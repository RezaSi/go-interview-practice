package main

import (
	"bufio"
	"fmt"
	"os"
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
	left := 0
	right := len(s) - 1
	for left < right {
		runes[left], runes[right] = runes[right], runes[left]
		left += 1
		right -= 1
	}
	return string(runes)
}
