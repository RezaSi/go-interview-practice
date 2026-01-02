// Package challenge10 contains the solution for Challenge 10.
package challenge10

import (
	"errors"
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
		return nil, errors.New("width and height should be positive numbers")
	}
	return &Rectangle{
		Width:  width,
		Height: height,
	}, nil
}

// Area calculates the area of the rectangle
func (r *Rectangle) Area() float64 {
	// DONE: Implement area calculation
	return r.Width * r.Height
}

// Perimeter calculates the perimeter of the rectangle
func (r *Rectangle) Perimeter() float64 {
	// DONE: Implement perimeter calculation
	return (r.Width + r.Height) * 2
}

// String returns a string representation of the rectangle
func (r *Rectangle) String() string {
	// DONE: Implement string representation
	return fmt.Sprintf("Rectangle(width=%.2f, height=%.2f)", r.Width, r.Height)
}

// Circle represents a perfectly round shape
type Circle struct {
	Radius float64
}

// NewCircle creates a new Circle with validation
func NewCircle(radius float64) (*Circle, error) {
	// DONE: Implement validation and construction
	if radius <= 0 {
		return nil, errors.New("radius should be a positive number")
	}
	return &Circle{
		Radius: radius,
	}, nil
}

// Area calculates the area of the circle
func (c *Circle) Area() float64 {
	// DONE: Implement area calculation
	return math.Pi * c.Radius * c.Radius
}

// Perimeter calculates the circumference of the circle
func (c *Circle) Perimeter() float64 {
	// DONE: Implement perimeter calculation
	return 2 * math.Pi * c.Radius
}

// String returns a string representation of the circle
func (c *Circle) String() string {
	// DONE: Implement string representation
	return fmt.Sprintf("Circle(radius=%.2f)", c.Radius)
}

// Triangle represents a three-sided polygon
type Triangle struct {
	SideA float64
	SideB float64
	SideC float64
}

// NewTriangle creates a new Triangle with validation
func NewTriangle(a, b, c float64) (*Triangle, error) {
	// DONE: Implement validation and construction
	if a <= 0 || b <= 0 || c <= 0 {
		return nil, errors.New("the three sides of a triangle should be positive numbers")
	} else if a+b <= c || a+c <= b || b+c <= a {
		return nil, errors.New("the sum of any two sides of a triangle should be greater than the third side")
	}
	return &Triangle{
		SideA: a,
		SideB: b,
		SideC: c,
	}, nil
}

// Area calculates the area of the triangle using Heron's formula
func (t *Triangle) Area() float64 {
	// DONE: Implement area calculation using Heron's formula
	a := t.SideA
	b := t.SideB
	c := t.SideC
	p := (a + b + c) / 2
	return math.Sqrt(p * (p - a) * (p - b) * (p - c))
}

// Perimeter calculates the perimeter of the triangle
func (t *Triangle) Perimeter() float64 {
	return t.SideA + t.SideB + t.SideC
}

// String returns a string representation of the triangle
func (t *Triangle) String() string {
	return fmt.Sprintf("Triangle(sides=%.2f, %.2f, %.2f)", t.SideA, t.SideB, t.SideC)
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
	var sum float64
	for _, v := range shapes {
		sum += v.Area()
	}
	return sum
}

// LargestShape finds the shape with the largest area
func (sc *ShapeCalculator) LargestShape(shapes []Shape) Shape {
	largest := 0.0
	index := -1
	for i, v := range shapes {
		if v.Area() > largest {
			largest = v.Area()
			index = i
		}
	}
	if index != -1 {
		return shapes[index]
	}
	return nil
}

// SortByArea sorts shapes by area in ascending or descending order
func (sc *ShapeCalculator) SortByArea(shapes []Shape, ascending bool) []Shape {
	sortedShapes := make([]Shape, len(shapes))
	copy(sortedShapes, shapes)

	sort.Slice(sortedShapes, func(i, j int) bool {
		areaI := sortedShapes[i].Area()
		areaJ := sortedShapes[j].Area()
		if ascending {
			return areaI < areaJ
		} else {
			return areaI > areaJ
		}
	})

	return sortedShapes
}
