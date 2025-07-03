// Package challenge10 contains the solution for Challenge 10.
package challenge10

import (
	"fmt"
	"math"
	"slices"
	// Add any necessary imports here
)

const (
	pi = math.Pi // Small value for floating point comparisons
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
	// validate positive dims
	if width <= 0 {
		return nil, fmt.Errorf("invalid width, [%.3f] must be > 0", width)
	}
	if height <= 0 {
		return nil, fmt.Errorf("invalid height, [%.3f] must be > 0", height)
	}

	// return new rectangle
	return &Rectangle{
		Width:  width,
		Height: height,
	}, nil
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
	return fmt.Sprintf("Rectangle - width=%.3f height=%.3f", r.Width, r.Height)
}

// Circle represents a perfectly round shape
type Circle struct {
	Radius float64
}

// NewCircle creates a new Circle with validation
func NewCircle(radius float64) (*Circle, error) {
	// validate positive dims
	if radius <= 0 {
		return nil, fmt.Errorf("invalid radius, [%.3f] must be > 0", radius)
	}

	// return new circle
	return &Circle{
		Radius: radius,
	}, nil
}

// Area calculates the area of the circle
func (c *Circle) Area() float64 {
	return pi * c.Radius * c.Radius
}

// Perimeter calculates the circumference of the circle
func (c *Circle) Perimeter() float64 {
	return 2 * pi * c.Radius
}

// String returns a string representation of the circle
func (c *Circle) String() string {
	return fmt.Sprintf("Circle - radius=%.3f ", c.Radius)
}

// Triangle represents a three-sided polygon
type Triangle struct {
	SideA float64
	SideB float64
	SideC float64
}

// NewTriangle creates a new Triangle with validation
func NewTriangle(a, b, c float64) (*Triangle, error) {
	// validate positive dims
	if a <= 0 || b <= 0 || c <= 0 {
		return nil, fmt.Errorf("sides must be positive, [%.3f] [%.3f] [%.3f]", a, b, c)
	}

	// validate inequality theorem
	if a+b <= c || a+c <= b || b+c <= a {
		return nil, fmt.Errorf("2 sides must be > 3rd side, [%.3f] [%.3f] [%.3f]", a, b, c)
	}

	// return new rectangle
	return &Triangle{
		SideA: a,
		SideB: b,
		SideC: c,
	}, nil
}

// Area calculates the area of the triangle using Heron's formula
func (t *Triangle) Area() float64 {
	s := t.Perimeter() / 2
	return math.Sqrt(s * (s - t.SideA) * (s - t.SideB) * (s - t.SideC))
}

// Perimeter calculates the perimeter of the triangle
func (t *Triangle) Perimeter() float64 {
	return t.SideA + t.SideB + t.SideC
}

// String returns a string representation of the triangle
func (t *Triangle) String() string {
	return fmt.Sprintf("Triangle - sides=%.3f, %.3f, %.3f", t.SideA, t.SideB, t.SideC)
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
	// init area
	var area float64

	// accumulate areas
	for _, sh := range shapes {
		area += sh.Area()
	}

	return area
}

// LargestShape finds the shape with the largest area
func (sc *ShapeCalculator) LargestShape(shapes []Shape) Shape {
	// init
	var (
		shape   Shape
		maxArea float64
	)

	// loop thru shapes tracking one with biggest area
	for _, sh := range shapes {
		if sh.Area() > maxArea {
			maxArea = sh.Area()
			shape = sh
		}
	}

	return shape
}

// SortByArea sorts shapes by area in ascending or descending order
func (sc *ShapeCalculator) SortByArea(shapes []Shape, ascending bool) []Shape {
	slices.SortFunc(shapes, func(a, b Shape) int {
		switch {
		case a.Area() < b.Area():
			if ascending {
				return -1
			}
			return 1
		case a.Area() > b.Area():
			if ascending {
				return 1
			}
			return -1
		default:
			return 0
		}
	})
	
	return shapes
}
