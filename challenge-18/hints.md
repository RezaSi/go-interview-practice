# Hints for Temperature Converter

## Hint 1: Temperature Conversion Formulas
Key conversion formulas:
- Celsius to Fahrenheit: `F = C × 9/5 + 32`
- Fahrenheit to Celsius: `C = (F - 32) × 5/9`
- Celsius to Kelvin: `K = C + 273.15`
- Kelvin to Celsius: `C = K - 273.15`

## Hint 2: Function Structure
Create separate conversion functions:
```go
func CelsiusToFahrenheit(celsius float64) float64 {
    return celsius*9/5 + 32
}

func FahrenheitToCelsius(fahrenheit float64) float64 {
    return (fahrenheit - 32) * 5 / 9
}
```

## Hint 3: Kelvin Conversions
Remember absolute zero for Kelvin:
```go
func CelsiusToKelvin(celsius float64) float64 {
    return celsius + 273.15
}

func KelvinToCelsius(kelvin float64) float64 {
    return kelvin - 273.15
}
```

## Hint 4: Chain Conversions
For Fahrenheit ↔ Kelvin, go through Celsius:
```go
func FahrenheitToKelvin(fahrenheit float64) float64 {
    celsius := FahrenheitToCelsius(fahrenheit)
    return CelsiusToKelvin(celsius)
}
```

## Hint 5: Input Validation
Validate temperatures against physical limits:
```go
func isValidKelvin(kelvin float64) bool {
    return kelvin >= 0 // Absolute zero
}

func isValidCelsius(celsius float64) bool {
    return celsius >= -273.15 // Absolute zero in Celsius
}
``` 