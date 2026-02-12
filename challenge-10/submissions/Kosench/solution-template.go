// Package challenge10 contains the solution for Challenge 10.
package challenge10

import (
	"errors"
	"fmt"
	"math"
	"sort"
	// Add any necessary imports here
)

// pi is a constant representing the mathematical constant π (pi)
const pi = math.Pi

// Shape defines the interface for geometric shapes.
// All shapes must implement Area, Perimeter, and String methods.
type Shape interface {
	Area() float64
	Perimeter() float64
	fmt.Stringer // Includes String() string method
}

// Rectangle represents a rectangle with width and height.
type Rectangle struct {
	Width  float64
	Height float64
}

// NewRectangle creates a new Rectangle with validation.
// Returns an error if width or height is not positive.
func NewRectangle(width, height float64) (*Rectangle, error) {
	if width <= 0 {
		return nil, errors.New("width must be positive")
	}
	if height <= 0 {
		return nil, errors.New("height must be positive")
	}
	return &Rectangle{
		Width:  width,
		Height: height,
	}, nil
}

// Area calculates and returns the area of the rectangle (width × height).
func (r *Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Perimeter calculates and returns the perimeter of the rectangle.
func (r *Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// String returns a string representation of the rectangle.
func (r *Rectangle) String() string {
	return fmt.Sprintf("Rectangle(width=%.2f, height=%.2f)", r.Width, r.Height)
}

// Circle represents a circle with a radius.
type Circle struct {
	Radius float64
}

// NewCircle creates a new Circle with validation
func NewCircle(radius float64) (*Circle, error) {
	if radius <= 0 {
		return nil, errors.New("radius must be positive")
	}
	return &Circle{Radius: radius}, nil
}

// Area calculates and returns the area of the circle (π × r²).
func (c *Circle) Area() float64 {
	return pi * c.Radius * c.Radius
}

// Perimeter calculates and returns the circumference of the circle (2 × π × r).
func (c *Circle) Perimeter() float64 {
	return 2 * pi * c.Radius
}

// String returns a string representation of the circle.
func (c *Circle) String() string {
	return fmt.Sprintf("Circle(radius=%.2f)", c.Radius)
}

// Triangle represents a triangle with three sides.
type Triangle struct {
	SideA float64
	SideB float64
	SideC float64
}

// NewTriangle creates a new Triangle with validation.
// Returns an error if any side is not positive or if the triangle inequality is violated.
func NewTriangle(a, b, c float64) (*Triangle, error) {
	// Проверка на положительность сторон
	if a <= 0 || b <= 0 || c <= 0 {
		return nil, errors.New("sides must be positive")
	}
	// Проверка неравенства треугольника: сумма любых двух сторон должна быть больше третьей
	if a+b <= c || a+c <= b || b+c <= a {
		return nil, errors.New("invalid triangle: triangle inequality violated")
	}
	return &Triangle{
		SideA: a,
		SideB: b,
		SideC: c,
	}, nil
}

// Area calculates the area of the triangle using Heron's formula
func (t *Triangle) Area() float64 {
	// Формула Герона: A = sqrt(s * (s-a) * (s-b) * (s-c)), где s - полупериметр
	s := t.Perimeter() / 2
	return math.Sqrt(s * (s - t.SideA) * (s - t.SideB) * (s - t.SideC))
}

// Perimeter calculates the perimeter of the triangle
func (t *Triangle) Perimeter() float64 {
	return t.SideA + t.SideB + t.SideC
}

// String returns a string representation of the triangle
func (t *Triangle) String() string {
	return fmt.Sprintf("Triangle(sides: %.2f, %.2f, %.2f)", t.SideA, t.SideB, t.SideC)
}

// ShapeCalculator provides utility functions for shapes
type ShapeCalculator struct{}

// NewShapeCalculator creates a new ShapeCalculator
func NewShapeCalculator() *ShapeCalculator {
	return &ShapeCalculator{}
}

// PrintProperties prints the properties of a shape
func (sc *ShapeCalculator) PrintProperties(s Shape) {
	fmt.Printf("%s - Area: %.2f, Perimeter: %.2f\n",
		s.String(), s.Area(), s.Perimeter())
}

// TotalArea calculates the sum of areas of all shapes
func (sc *ShapeCalculator) TotalArea(shapes []Shape) float64 {
	total := 0.0
	for _, shape := range shapes {
		total += shape.Area()
	}
	return total
}

// LargestShape finds the shape with the largest area
func (sc *ShapeCalculator) LargestShape(shapes []Shape) Shape {
	if len(shapes) == 0 {
		return nil
	}

	largest := shapes[0]
	maxArea := largest.Area()

	for _, shape := range shapes[1:] {
		if area := shape.Area(); area > maxArea {
			maxArea = area
			largest = shape
		}
	}

	return largest
}

// SortByArea sorts shapes by area in ascending or descending order
func (sc *ShapeCalculator) SortByArea(shapes []Shape, ascending bool) []Shape {
	result := make([]Shape, len(shapes))
	copy(result, shapes)

	sort.Slice(result, func(i, j int) bool {
		areaI := result[i].Area()
		areaJ := result[j].Area()

		if ascending {
			return areaI < areaJ
		}
		return areaI > areaJ
	})

	return result
}
