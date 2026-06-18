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
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("width and height must be > 0")
	}
	return &Rectangle{
		Width:  width,
		Height: height,
	}, nil
}

// Area calculates the area of the rectangle
func (r *Rectangle) Area() float64 {
	// TODO: Implement area calculation
	if r == nil {
		return 0
	}
	return r.Height * r.Width
}

// Perimeter calculates the perimeter of the rectangle
func (r *Rectangle) Perimeter() float64 {
	// TODO: Implement perimeter calculation
	if r == nil {
		return 0
	}
	return 2 * (r.Height + r.Width)
}

// String returns a string representation of the rectangle
func (r *Rectangle) String() string {
	// TODO: Implement string representation
	if r == nil {
		return ""
	}
	return fmt.Sprintf("Rectangle, width=%f, height=%f", r.Width, r.Height)
}

// Circle represents a perfectly round shape
type Circle struct {
	Radius float64
}

// NewCircle creates a new Circle with validation
func NewCircle(radius float64) (*Circle, error) {
	// TODO: Implement validation and construction
	if radius <= 0 {
		return nil, fmt.Errorf("radius must be > 0")
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
	return fmt.Sprintf("Circle, radius=%f", c.Radius)
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
	if a+b <= c || a+c <= b || b+c <= a {
		return nil, fmt.Errorf("invalid parameter")
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
	if t == nil {
		return 0
	}
	s := (t.SideA + t.SideB + t.SideC) / 2
	A := math.Sqrt(s * (s - t.SideA) * (s - t.SideB) * (s - t.SideC))
	return A
}

// Perimeter calculates the perimeter of the triangle
func (t *Triangle) Perimeter() float64 {
	// TODO: Implement perimeter calculation
	if t == nil {
		return 0
	}
	return t.SideA + t.SideB + t.SideC
}

// String returns a string representation of the triangle
func (t *Triangle) String() string {
	if t == nil {
		return ""
	}
	// TODO: Implement string representation
	return fmt.Sprintf("Triangle, sides:%f, %f, %f", t.SideA, t.SideB, t.SideC)
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
	fmt.Println(s.String())
}

// TotalArea calculates the sum of areas of all shapes
func (sc *ShapeCalculator) TotalArea(shapes []Shape) (sum float64) {
	// TODO: Implement total area calculation
	if len(shapes) == 0 {
		return 0
	}
	for _, shape := range shapes {
		sum += shape.Area()
	}
	return
}

// LargestShape finds the shape with the largest area
func (sc *ShapeCalculator) LargestShape(shapes []Shape) Shape {
	// TODO: Implement finding largest shape
	if len(shapes) == 0 {
		return nil
	}
	largestShape := shapes[0]
	for i := 1; i < len(shapes); i++ {
		if shapes[i].Area() > largestShape.Area() {
			largestShape = shapes[i]
		}
	}
	return largestShape
}

// SortByArea sorts shapes by area in ascending or descending order
func (sc *ShapeCalculator) SortByArea(shapes []Shape, ascending bool) []Shape {
	// TODO: Implement sorting shapes by area
	if ascending {
		sort.Slice(shapes, func(i, j int) bool {
			return shapes[i].Area() < shapes[j].Area()
		})
	} else {
		sort.Slice(shapes, func(i, j int) bool {
			return shapes[i].Area() > shapes[j].Area()
		})
	}
	return shapes
}
