package main

import (
    "fmt"
    "bufio"
    "os"
)

func main(){
    reader := bufio.NewReader(os.Stdin)
    
    name, _ := reader.ReadString('\n')
    
    fmt.Println(MyBestReverseString(name))
}

func ReverseString(s string) string{
    rs:=""
    for i:=len(s)-1; i>=0; i--{
        rs += string(s[i])
    }
    return rs
}

func MyBestReverseString(s string) string{
    runes := []rune(s)
    
    for i,j := 0, len(runes)-1; i<j; i,j = i+1, j-1{
        runes[i], runes[j] = runes[j], runes[i]
    }
    
    return string(runes)
}