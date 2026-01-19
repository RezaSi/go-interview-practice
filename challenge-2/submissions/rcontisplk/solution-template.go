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

	runeSlice := []rune(s)
	var backwardsFirst []rune

	for i := len(s); i > 0; i-- {
		backwardsFirst = append(backwardsFirst, runeSlice[i-1])
	}

	backwardsStringFirst := string(backwardsFirst)

	return backwardsStringFirst
}
