package main

import (
	"fmt"
	"math"
)

func main() {
	celsius := 25.0
	fahrenheit := CelsiusToFahrenheit(celsius)
	fmt.Printf("%.2f°C is equal to %.2f°F\n", celsius, fahrenheit)

	fahrenheit = 68.0
	celsius = FahrenheitToCelsius(fahrenheit)
	fmt.Printf("%.2f°F is equal to %.2f°C\n", fahrenheit, celsius)
}

func CelsiusToFahrenheit(celsius float64) float64 {return celsius*9/5 + 32}

func FahrenheitToCelsius(fahrenheit float64) float64 {
	return (fahrenheit - 32) * 5 / 9
}

func Round(value float64, decimals int) float64 {
	precision := math.Pow10(decimals)
	return math.Round(value*precision) / precision
}
