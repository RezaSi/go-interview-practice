// Package challenge10 contains the solution for Challenge 10.
package challenge10

import (
	"fmt"
	"math"
	"slices"
	// Add any necessary imports here
)

type LineError struct {
	Message string
}

func (e *LineError) Error() string {
	return e.Message
}

type ShapeError struct {
	Message string
}

func (s *ShapeError) Error() string {
	return s.Message
}

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

// NewRectangle creates a new Rectangle with validation
func NewRectangle(width, height float64) (*Rectangle, error) {
	if width <= 0 {
		return nil, &LineError{"width must be > 0"}
	}
	if height <= 0 {
		return nil, &LineError{"height must be > 0"}
	}
	return &Rectangle{
		Width:  width,
		Height: height,
	}, nil
}

// Area calculates the area of the rectangle
func (r *Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Perimeter calculates the perimeter of the rectangle
func (r *Rectangle) Perimeter() float64 {
	return (r.Width + r.Height) * 2
}

// String returns a string representation of the rectangle
func (r *Rectangle) String() string {
	return fmt.Sprintf("Rectangle: width %f, height %f", r.Width, r.Height)
}

// Circle represents a perfectly round shape
type Circle struct {
	Radius float64
}

// NewCircle creates a new Circle with validation
func NewCircle(radius float64) (*Circle, error) {
	if radius <= 0 {
		return nil, &LineError{"radius must be > 0"}
	}
	return &Circle{
		Radius: radius,
	}, nil
}

// Area calculates the area of the circle
func (c *Circle) Area() float64 {
	return math.Pi * math.Pow(c.Radius, 2)
}

// Perimeter calculates the circumference of the circle
func (c *Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// String returns a string representation of the circle
func (c *Circle) String() string {
	return fmt.Sprintf("Circle: radius %f", c.Radius)
}

// Triangle represents a three-sided polygon
type Triangle struct {
	SideA float64
	SideB float64
	SideC float64
}

// NewTriangle creates a new Triangle with validation
func NewTriangle(a, b, c float64) (*Triangle, error) {
	if a <= 0 {
		return nil, &LineError{"side A must be > 0"}
	}
	if b <= 0 {
		return nil, &LineError{"side B must be > 0"}
	}
	if c <= 0 {
		return nil, &LineError{"side C must be > 0"}
	}
	if a+b <= c || a+c <= b || c+b <= a {
		return nil, &ShapeError{"triangle inequality violated"}
	}

	return &Triangle{
		SideA: a,
		SideB: b,
		SideC: c,
	}, nil
}

// Area calculates the area of the triangle using Heron's formula
func (t *Triangle) Area() float64 {
	p := t.Perimeter() / 2.0
	return math.Sqrt(p * (p - t.SideA) * (p - t.SideB) * (p - t.SideC))
}

// Perimeter calculates the perimeter of the triangle
func (t *Triangle) Perimeter() float64 {
	return t.SideA + t.SideB + t.SideC
}

// String returns a string representation of the triangle
func (t *Triangle) String() string {
	return fmt.Sprintf("Triangle: sides A %f B %f C %f", t.SideA, t.SideB, t.SideC)
}

// ShapeCalculator provides utility functions for shapes
type ShapeCalculator struct{}

// NewShapeCalculator creates a new ShapeCalculator
func NewShapeCalculator() *ShapeCalculator {
	return &ShapeCalculator{}
}

// PrintProperties prints the properties of a shape
func (sc *ShapeCalculator) PrintProperties(s Shape) {
	fmt.Println(s)
}

// TotalArea calculates the sum of areas of all shapes
func (sc *ShapeCalculator) TotalArea(shapes []Shape) float64 {
	totalArea := 0.0
	for _, shape := range shapes {
		totalArea += shape.Area()
	}
	return totalArea
}

// LargestShape finds the shape with the largest area
func (sc *ShapeCalculator) LargestShape(shapes []Shape) Shape {
	return slices.MaxFunc(shapes, sc.compareShape)
}

// SortByArea sorts shapes by area in ascending or descending order
func (sc *ShapeCalculator) SortByArea(shapes []Shape, ascending bool) []Shape {
	slices.SortFunc(shapes, func(a, b Shape) int {
		if ascending {
			return sc.compareShape(a, b)
		}
		return sc.compareShape(b, a)
	})
	return shapes
}

func (sc *ShapeCalculator) compareShape(a, b Shape) int {
	switch {
	case a.Area() < b.Area():
		return -1
	case a.Area() > b.Area():
		return 1
	default:
		return 0
	}
}
