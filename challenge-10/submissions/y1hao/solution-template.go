// Package challenge10 contains the solution for Challenge 10.
package challenge10

import (
	"errors"
	"fmt"
	"math"
	"slices"
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
		return nil, errors.New("width and height must be greater than 0")
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
	return fmt.Sprintf("Rectangle[ Width = %.2f, Height = %.2f ]", r.Width, r.Height)
}

// Circle represents a perfectly round shape
type Circle struct {
	Radius float64
}

// NewCircle creates a new Circle with validation
func NewCircle(radius float64) (*Circle, error) {
	if radius <= 0 {
		return nil, errors.New("radius must be creater than 0")
	}
	return &Circle{
		Radius: radius,
	}, nil
}

// Area calculates the area of the circle
func (c *Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// Perimeter calculates the circumference of the circle
func (c *Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// String returns a string representation of the circle
func (c *Circle) String() string {
	return fmt.Sprintf("Circle[ Radius = %.2f ]", c.Radius)
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
		return nil, errors.New("sides must be greater than 0")
	}
	if a+b <= c || a-b > c || b-a > c {
		return nil, errors.New("sides cannot form a triangle")
	}
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
	return fmt.Sprintf("Triangle[ Sides = %.2f, %.2f, %.2f ]", t.SideA, t.SideB, t.SideC)
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
	total := 0.0
	for _, s := range shapes {
		total += s.Area()
	}
	return total
}

// LargestShape finds the shape with the largest area
func (sc *ShapeCalculator) LargestShape(shapes []Shape) Shape {
	largest := shapes[0]
	largestArea := largest.Area()
	for i := 1; i < len(shapes); i++ {
		area := shapes[i].Area()
		if area > largestArea {
			largest = shapes[i]
			largestArea = area
		}
	}
	return largest
}

// SortByArea sorts shapes by area in ascending or descending order
func (sc *ShapeCalculator) SortByArea(shapes []Shape, ascending bool) []Shape {
	slices.SortFunc(shapes, func(a, b Shape) int {
		if a.Area() == b.Area() {
			return 0
		}
		if a.Area() > b.Area() {
			return 1
		}
		return -1
	})
	if !ascending {
		slices.Reverse(shapes)
	}
	return shapes
}
