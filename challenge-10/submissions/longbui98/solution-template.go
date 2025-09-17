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
	if width <= 0 {
		return nil, fmt.Errorf("NewRectangle width must be greater than zero")
	}
	if height <= 0 {
		return nil, fmt.Errorf("NewRectangle height must be greater than zero")
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
	// TODO: Implement perimeter calculation
	return (r.Width + r.Height) * 2
}

// String returns a string representation of the rectangle
func (r *Rectangle) String() string {
	return fmt.Sprintf("Rectangle(width=%.2f, height=%.2f, area=%.2f)", r.Width, r.Height, r.Perimeter())
}

// Circle represents a perfectly round shape
type Circle struct {
	Radius float64
}

// NewCircle creates a new Circle with validation
func NewCircle(radius float64) (*Circle, error) {
	// TODO: Implement validation and construction
	if radius <= 0 {
		return nil, fmt.Errorf("NewCircle radius must be greater than zero")
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
	// TODO: Implement perimeter calculation
	return 2 * math.Pi * c.Radius
}

// String returns a string representation of the circle
func (c *Circle) String() string {
	// TODO: Implement string representation
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
	if a <= 0 || b <= 0 || c <= 0 {
		return nil, fmt.Errorf("NewTriangle a: %v, b: %v, c: %v", a, b, c)
	}
	if a+b == c {
		return nil, fmt.Errorf("NewTriangle a: %v, b: %v, c: %v", a, b, c)
	}
	if a+b < c {
		return nil, fmt.Errorf("NewTriangle a: %v, b: %v, c: %v", a, b, c)
	}
	return &Triangle{
		SideA: a,
		SideB: b,
		SideC: c,
	}, nil
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
	return fmt.Sprintf("Triangle(a=%.2f, b=%.2f, c=%.2f) sides", t.SideA, t.SideB, t.SideC)
}

// ShapeCalculator provides utility functions for shapes
type ShapeCalculator struct {
}

// NewShapeCalculator creates a new ShapeCalculator
func NewShapeCalculator() *ShapeCalculator {
	// TODO: Implement constructor
	return &ShapeCalculator{}
}

// PrintProperties prints the properties of a shape
func (sc *ShapeCalculator) PrintProperties(s Shape) {
	// TODO: Implement printing shape properties
	_ = fmt.Sprintf("Shape %d: %T, Area: %.2f, Perimeter: %.2f, String: %s", s, s.String(), s.Area(), s.Perimeter(), s.String())
}

// TotalArea calculates the sum of areas of all shapes
func (sc *ShapeCalculator) TotalArea(shapes []Shape) float64 {
	// TODO: Implement total area calculation
	total := float64(0)
	for _, v := range shapes {
		total += v.Area()
	}
	return total
}

// LargestShape finds the shape with the largest area
func (sc *ShapeCalculator) LargestShape(shapes []Shape) Shape {
	// TODO: Implement finding largest shape
	largestArea := shapes[0]
	for i := 1; i < len(shapes); i++ {
		if shapes[i].Area() > largestArea.Area() {
			largestArea = shapes[i]
		}
	}
	return largestArea
}

// SortByArea sorts shapes by area in ascending or descending order
func (sc *ShapeCalculator) SortByArea(shapes []Shape, ascending bool) []Shape {
	// TODO: Implement sorting shapes by area
	sort.Slice(shapes, func(i, j int) bool {
		if ascending {
			return shapes[i].Area() < shapes[j].Area()
		}
		return shapes[i].Area() > shapes[j].Area()
	})
	return shapes
}
