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
	// TODO: Implement validation and construction

	if width <= 0.0 || height <= 0.0 {
		return nil, fmt.Errorf("negative parameters: %.2f, %.2f", width, height)
	}

	return &Rectangle{Width: width, Height: height}, nil
}

// Area calculates the area of the rectangle
func (r *Rectangle) Area() float64 {
	// TODO: Implement area calculation
	return r.Height * r.Width
}

// Perimeter calculates the perimeter of the rectangle
func (r *Rectangle) Perimeter() float64 {
	// TODO: Implement perimeter calculation
	return r.Height*2 + r.Width*2
}

// String returns a string representation of the rectangle
func (r *Rectangle) String() string {
	// TODO: Implement string representation
	return fmt.Sprintf("Rectangle: (Width: %.2f, Height: %.2f)", r.Width, r.Height)
}

// Circle represents a perfectly round shape
type Circle struct {
	Radius float64
}

// NewCircle creates a new Circle with validation
func NewCircle(radius float64) (*Circle, error) {
	// TODO: Implement validation and construction
	if radius <= 0.0 {
		return nil, fmt.Errorf("negative parameters: %.2f", radius)
	}
	return &Circle{Radius: radius}, nil
}

// Area calculates the area of the circle
func (c *Circle) Area() float64 {
	// TODO: Implement area calculation
	return math.Pi * c.Radius * c.Radius
}

// Perimeter calculates the circumference of the circle
func (c *Circle) Perimeter() float64 {
	// TODO: Implement perimeter calculation
	return 2 * math.Pi * c.Radius
}

// String returns a string representation of the circle
func (c *Circle) String() string {
	// TODO: Implement string representation
	return fmt.Sprintf("Circle: (Radius: %.2f)", c.Radius)
}

// Triangle represents a three-sided polygon
type Triangle struct {
	SideA float64
	SideB float64
	SideC float64
}

// NewTriangle creates a new Triangle with validation
func NewTriangle(a, b, c float64) (*Triangle, error) {
	// TODO: Implement validation and construction
	if a <= 0.0 || b <= 0.0 || c <= 0.0 {
		return nil, fmt.Errorf("negative parameters: %.2f, %.2f, %.2f", a, b, c)
	}

	if a >= b+c || b >= a+c || c >= a+b {
		return nil, fmt.Errorf("shape is not Triangle: %.2f, %.2f, %.2f", a, b, c)
	}

	return &Triangle{SideA: a, SideB: b, SideC: c}, nil
}

// Area calculates the area of the triangle using Heron's formula
func (t *Triangle) Area() float64 {
	// TODO: Implement area calculation using Heron's formula
	s := (t.SideA + t.SideB + t.SideC) / 2
	return math.Sqrt(s * (s - t.SideA) * (s - t.SideB) * (s - t.SideC))

}

// Perimeter calculates the perimeter of the triangle
func (t *Triangle) Perimeter() float64 {
	// TODO: Implement perimeter calculation
	return t.SideA + t.SideB + t.SideC
}

// String returns a string representation of the triangle
func (t *Triangle) String() string {
	// TODO: Implement string representation
	return fmt.Sprintf("Triangle: (sides a=%.2f, b=%.2f, c=%.2f)", t.SideA, t.SideB, t.SideC)
}

// ShapeCalculator provides utility functions for shapes
type ShapeCalculator struct{}

// NewShapeCalculator creates a new ShapeCalculator
func NewShapeCalculator() *ShapeCalculator {
	// TODO: Implement constructor
	return &ShapeCalculator{}
}

// PrintProperties prints the properties of a shape
func (sc *ShapeCalculator) PrintProperties(s Shape) {
	// TODO: Implement printing shape properties
	s.String()
}

// TotalArea calculates the sum of areas of all shapes
func (sc *ShapeCalculator) TotalArea(shapes []Shape) float64 {

	total := 0.0
	for _, s := range shapes {
		total += s.Area()
	}
	// TODO: Implement total area calculation
	return total
}

// LargestShape finds the shape with the largest area
func (sc *ShapeCalculator) LargestShape(shapes []Shape) Shape {
	// TODO: Implement finding largest shape
	if len(shapes) <= 0 {
		return nil
	}

	shape := shapes[0]
	for _, s := range shapes[1:] {
		if s.Area() > shape.Area() {
			shape = s
		}
	}
	return shape
}

// SortByArea sorts shapes by area in ascending or descending order
func (sc *ShapeCalculator) SortByArea(shapes []Shape, ascending bool) []Shape {
	// TODO: Implement sorting shapes by area
	if len(shapes) <= 0 {
		return nil
	}

	shapesCopy := make([]Shape, len(shapes))
	copy(shapesCopy, shapes)

	sort.Slice(shapesCopy, func(i int, j int) bool {
		if ascending {
			return shapesCopy[i].Area() < shapesCopy[j].Area()
		}
		return shapesCopy[i].Area() > shapesCopy[j].Area()
	})

	return shapesCopy
}
