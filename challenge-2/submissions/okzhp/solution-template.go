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

		// 调用ReverseString函数
		output := ReverseString(input)

		// 打印结果
		fmt.Println(output)
	}
}

// ReverseString返回反转后的字符串
func ReverseString(s string) string {
	runes := []rune(s)
	length := len(runes)
	for i := 0; i < length/2; i++ {
		tmp := runes[i]
		runes[i] = runes[length-i-1]
		runes[length-i-1] = tmp
	}
	return string(runes)
}
