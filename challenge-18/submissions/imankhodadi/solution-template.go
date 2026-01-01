package main

import (
	"fmt"
	"math"
)

func main() {
	fmt.Println(CelsiusToFahrenheit(25.0))
	fmt.Println(FahrenheitToCelsius(68.0))
}
func CelsiusToFahrenheit(celsius float64) float64 {
	if err := ValidateCelsius(celsius); err != nil {
		return math.NaN()
	}
	return Round(celsius*9.0/5.0+32.0, 2)
}

func FahrenheitToCelsius(fahrenheit float64) float64 {
	if err := ValidateFahrenheit(fahrenheit); err != nil {
		return math.NaN()
	}
	return Round((fahrenheit-32.0)*5.0/9.0, 2)
}

func ValidateCelsius(celsius float64) error {
	if celsius < -273.15 {
		return fmt.Errorf("temperature below absolute zero: %f°C", celsius)
	}
	return nil
}

func ValidateFahrenheit(fahrenheit float64) error {
	if fahrenheit < -459.67 {
		return fmt.Errorf("temperature below absolute zero: %f°F", fahrenheit)
	}
	return nil
}

func Round(value float64, decimals int) float64 {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return value
	}
	precision := math.Pow10(decimals)
	return math.Round(value*precision) / precision
}