package main

import ("fmt")

func main(){
    var a, b int
    
    fmt.Scanf("%d, %d", &a, &b)
    fmt.Println(Sum(a,b))
    
}

func Sum(a int, b int) int{
    return a+b
}