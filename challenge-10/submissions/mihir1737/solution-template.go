// Package challenge10 contains the solution for Challenge 10.
package challenge10

import (
	"errors"
	"fmt"
	"math"
	"sort"
)

// Shape interface defines methods that all shapes must implement
type Shape interface {
	Area() float64
	Perimeter() float64
	fmt.Stringer // Includes String() string method
}

// Rectangle represents a four-sided shape with perpendicular sides
type Rectangle struct {
	Width  float64
	Height float64
}

var InvalidInputError = errors.New("Invalid value")

// NewRectangle creates a new Rectangle with validation
func NewRectangle(width, height float64) (*Rectangle, error) {
	if width <= 0 || height <= 0 {
		return nil, InvalidInputError
	}
	rec := Rectangle{Width: width, Height: height}
	return &rec, nil
}

// Area calculates the area of the rectangle
func (r *Rectangle) Area() float64 {
	return r.Height * r.Width
}

// Perimeter calculates the perimeter of the rectangle
func (r *Rectangle) Perimeter() float64 {
	return 2 * (r.Height + r.Width)
}

// String returns a string representation of the rectangle
func (r *Rectangle) String() string {
	return fmt.Sprintf("Rectangle with width %.2f and height %.2f", r.Width, r.Height)
}

// Circle represents a perfectly round shape
type Circle struct {
	Radius float64
}

// NewCircle creates a new Circle with validation
func NewCircle(radius float64) (*Circle, error) {
	if radius <= 0 {
		return nil, InvalidInputError
	}
	c := Circle{Radius: radius}
	return &c, nil
}

// Area calculates the area of the circle
func (c *Circle) Area() float64 {
	return c.Radius * c.Radius * math.Pi
}

// Perimeter calculates the circumference of the circle
func (c *Circle) Perimeter() float64 {
	return c.Radius * math.Pi * 2
}

// String returns a string representation of the circle
func (c *Circle) String() string {
	return fmt.Sprintf("Circle with radius: %2f", c.Radius)
}

// Triangle represents a three-sided polygon
type Triangle struct {
	SideA float64
	SideB float64
	SideC float64
}

// NewTriangle creates a new Triangle with validation
func NewTriangle(a, b, c float64) (*Triangle, error) {
	if a <= 0 || b <= 0 || c <= 0 {
		return nil, InvalidInputError
	}

	if a+b <= c || b+c <= a || a+c <= b {
		return nil, InvalidInputError
	}

	t := Triangle{SideA: a, SideB: b, SideC: c}
	return &t, nil
}

// Area calculates the area of the triangle using Heron's formula
func (t *Triangle) Area() float64 {
	s := (t.SideA + t.SideB + t.SideC) / 2
	return math.Sqrt(s * (s - t.SideA) * (s - t.SideB) * (s - t.SideC))
}

// Perimeter calculates the perimeter of the triangle
func (t *Triangle) Perimeter() float64 {
	return t.SideA + t.SideB + t.SideC
}

// String returns a string representation of the triangle
func (t *Triangle) String() string {
	return fmt.Sprintf("Triangle %2f, %2f, %2f sides", t.SideA, t.SideB, t.SideC)
}

// ShapeCalculator provides utility functions for shapes
type ShapeCalculator struct {
}

// NewShapeCalculator creates a new ShapeCalculator
func NewShapeCalculator() *ShapeCalculator {
	return &ShapeCalculator{}
}

// PrintProperties prints the properties of a shape
func (sc *ShapeCalculator) PrintProperties(s Shape) {
	fmt.Printf("Perimeter : %2f", s.Perimeter())
	fmt.Printf("Area : %2f", s.Area())
}

// TotalArea calculates the sum of areas of all shapes
func (sc *ShapeCalculator) TotalArea(shapes []Shape) float64 {
	total := 0.00

	for _, s := range shapes {
		total += s.Area()
	}
	return total
}

// LargestShape finds the shape with the largest area
func (sc *ShapeCalculator) LargestShape(shapes []Shape) Shape {
	var ans Shape
	var maxArea float64 = 0.00

	for _, s := range shapes {
		if s.Area() > maxArea {
			ans = s
			maxArea = s.Area()
		}
	}

	return ans
}

// SortByArea sorts shapes by area in ascending or descending order
func (sc *ShapeCalculator) SortByArea(shapes []Shape, ascending bool) []Shape {
	sort.Slice(shapes, func(i, j int) bool {
		if ascending {
			return shapes[i].Area() <= shapes[j].Area()
		}
		return shapes[i].Area() >= shapes[j].Area()
	})
	return shapes
}
