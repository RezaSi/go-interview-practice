package main

import (
	"fmt"
	"math"
)

func main() {
	// Example usage
	celsius := 0.0
	fahrenheit := CelsiusToFahrenheit(celsius)
	fmt.Printf("%.2f°C is equal to %.2f°F\n", celsius, fahrenheit)

	fahrenheit = 68.0
	celsius = FahrenheitToCelsius(fahrenheit)
	fmt.Printf("%.2f°F is equal to %.2f°C\n", fahrenheit, celsius)
}

// CelsiusToFahrenheit converts a temperature from Celsius to Fahrenheit
// Formula: F = C × 9/5 + 32
func CelsiusToFahrenheit(celsius float64) float64 {
	if celsius < -273.15 {
		fmt.Printf("temperature below absolute zero %fC", celsius)
	}
	f := celsius*9.0/5.0 + 32
	fahrenheit := Round(f, 2)
	return fahrenheit
}

// FahrenheitToCelsius converts a temperature from Fahrenheit to Celsius
// Formula: C = (F - 32) × 5/9
func FahrenheitToCelsius(fahrenheit float64) float64 {
	if fahrenheit < -459.67 {
		fmt.Printf("temperature below absolute zero %fF", fahrenheit)
	}
	c := (fahrenheit - 32) * 5.0 / 9.0
	celsius := Round(c, 2)
	return celsius
}

// Round rounds a float64 value to the specified number of decimal places
func Round(value float64, decimals int) float64 {
	precision := math.Pow10(decimals)
	return math.Round(value*precision) / precision
}
