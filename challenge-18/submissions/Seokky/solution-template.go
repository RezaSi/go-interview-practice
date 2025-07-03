package main

import "math"

func main() {
}

func CelsiusToFahrenheit(celsius float64) float64 {
	return Round((celsius * 9/5) + 32, 2)
}

func FahrenheitToCelsius(fahrenheit float64) float64 {
	return Round((fahrenheit - 32) * 5/9, 2)
}

func Round(value float64, decimals int) float64 {
	precision := math.Pow10(decimals)
	return math.Round(value*precision) / precision
}
