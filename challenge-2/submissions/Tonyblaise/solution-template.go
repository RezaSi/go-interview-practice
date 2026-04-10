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
	runes := []rune{}

	for _, r := range s {
		runes = append([]rune{r}, runes...)

	}

	return string(runes)
}
