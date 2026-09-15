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
    for i,j := 0, len(runeSlice)-1; i<j; i,j = i+1,j-1 {
        runeSlice[i], runeSlice[j] = runeSlice[j], runeSlice[i]
    }
	return string(runeSlice)
}
