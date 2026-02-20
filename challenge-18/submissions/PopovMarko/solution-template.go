package main

import (
	"fmt"
	"math"
	"os"
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
	if celsius < -273.15 {
		fmt.Printf("temperature below absolute zero %.2f\n", celsius)
		os.Exit(1)
	}
	f := celsius*9.0/5.0 + 32

	return round(f, 2)
}

// FahrenheitToCelsius converts a temperature from Fahrenheit to Celsius
// Formula: C = (F - 32) × 5/9
func FahrenheitToCelsius(fahrenheit float64) float64 {
	if fahrenheit < -459.67 {
		fmt.Printf("temperature below absolute zero %.2f\n", fahrenheit)
	}
	c := (fahrenheit - 32) * 5.0 / 9.0
	return round(c, 2)
}

// round rounds a float64 value to the specified number of decimal places
func round(value float64, decimals int) float64 {
	precision := math.Pow10(decimals)
	return math.Round(value*precision) / precision
}
