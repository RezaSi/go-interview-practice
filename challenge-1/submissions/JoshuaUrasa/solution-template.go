package main

import "fmt"


func Sum(a int , b int) int{
    return a +b
}

func main(){
    fmt.Println(Sum(2,3))
    fmt.Println(Sum(-5,10))
}