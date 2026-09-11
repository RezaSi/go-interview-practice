// Package challenge10 contains the solution for Challenge 10.
package challenge10

import (
	"fmt"
	"math"
	"sort"
	// Add any necessary imports here
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

// NewRectangle creates a new Rectangle with validation
func NewRectangle(width, height float64) (*Rectangle, error) {
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("width and height must be positive")
	}
	return &Rectangle{width, height}, nil
}

// Area calculates the area of the rectangle
func (r *Rectangle) Area() float64 {
	return r.Height * r.Width
}

// Perimeter calculates the perimeter of the rectangle
func (r *Rectangle) Perimeter() float64 {
	return (r.Height + r.Width) * 2
}

// String returns a string representation of the rectangle
func (r *Rectangle) String() string {
	return fmt.Sprintf("Rectangle(width=%v, height=%v)", r.Width, r.Height)
}

// Circle represents a perfectly round shape
type Circle struct {
	Radius float64
}

// NewCircle creates a new Circle with validation
func NewCircle(radius float64) (*Circle, error) {
	if radius <= 0 {
		return nil, fmt.Errorf("radius must be positive")
	}
	return &Circle{radius}, nil
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
	return fmt.Sprintf("Circle(radius=%v)", c.Radius)
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
		return nil, fmt.Errorf("a, b, c must be positive")
	}
	if a+b <= c || a+c <= b || b+c <= a {
		return nil, fmt.Errorf("not a valid triangle")
	}
	return &Triangle{a, b, c}, nil
}

// Area calculates the area of the triangle using Heron's formula
func (t *Triangle) Area() float64 {
	return math.Sqrt(
		((t.SideA + t.SideB + t.SideC) / 2.0) *
			(((t.SideA + t.SideB + t.SideC) / 2.0) - t.SideA) *
			(((t.SideA + t.SideB + t.SideC) / 2.0) - t.SideB) *
			(((t.SideA + t.SideB + t.SideC) / 2.0) - t.SideC),
	)
}

// Perimeter calculates the perimeter of the triangle
func (t *Triangle) Perimeter() float64 {
	return t.SideA + t.SideB + t.SideC
}

// String returns a string representation of the triangle
func (t *Triangle) String() string {
	return fmt.Sprintf("Triangle(sides=%v,%v,%v)", t.SideA, t.SideB, t.SideC)
}

// ShapeCalculator provides utility functions for shapes
type ShapeCalculator struct{}

// NewShapeCalculator creates a new ShapeCalculator
func NewShapeCalculator() *ShapeCalculator {
	return &ShapeCalculator{}
}

// PrintProperties prints the properties of a shape
func (sc *ShapeCalculator) PrintProperties(s Shape) {
	fmt.Println(s.String())
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
	if len(shapes) == 0 {
		return nil
	}
	largestArea := 0.0
	var lastShape Shape
	for _, shape := range shapes {
		if shape.Area() > largestArea {
			largestArea = shape.Area()
			lastShape = shape
		}
	}
	return lastShape
}

// SortByArea sorts shapes by area in ascending or descending order
func (sc *ShapeCalculator) SortByArea(shapes []Shape, ascending bool) []Shape {
	result := make([]Shape, len(shapes))
	copy(result, shapes)
	if ascending {
		sort.Slice(result, func(i, j int) bool { return result[i].Area() < result[j].Area() })
	} else {
		sort.Slice(result, func(i, j int) bool { return result[i].Area() > result[j].Area() })
	}
	return result
}
