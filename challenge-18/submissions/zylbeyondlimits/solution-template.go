package main

import (
	"fmt"
	"math"
)

func main() {
	// Example usage
	celsius := 25.0
	fahrenheit := CelsiusToFahrenheit(celsius)
	fmt.Printf("%.2f°C is equal to %.2f°F\n", celsius, fahrenheit)

	fahrenheit = 68.0
	celsius = FahrenheitToCelsius(fahrenheit)
	fmt.Printf("%.2f°F is equal to %.2f°C\n", fahrenheit, celsius)
}

// CelsiusToFahrenheit converts a temperature from Celsius to Fahrenheit
// Formula: F = C × 9/5 + 32
func CelsiusToFahrenheit(celsius float64) float64 {
	// 使用 9.0/5.0 确保进行的是浮点数除法（结果为 1.8）
	// 如果使用 9/5 (整数除法)，结果会变成 1，导致计算错误
	f := celsius*(9.0/5.0) + 32
	return Round(f, 2)
}

// FahrenheitToCelsius converts a temperature from Fahrenheit to Celsius
// Formula: C = (F - 32) × 5/9
func FahrenheitToCelsius(fahrenheit float64) float64 {
	// 同样使用 5.0/9.0 确保浮点数精度
	c := (fahrenheit - 32) * (5.0 / 9.0)
	return Round(c, 2)
}

// Round rounds a float64 value to the specified number of decimal places
func Round(value float64, decimals int) float64 {
	precision := math.Pow10(decimals)
	return math.Round(value*precision) / precision
}