// Package challenge10 contains the solution for Challenge 10.
package challenge10

import (
	"fmt"
	"math"
	"sort"
	"errors"
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
		return nil, errors.New("height or width should be positive")
	}

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
	return 2 * (r.Height + r.Width)
}

// String returns a string representation of the rectangle
func (r *Rectangle) String() string {
    return fmt.Sprintf("Rectangle height: %g, Rectangle width: %g", r.Height, r.Width)
}

// Circle represents a perfectly round shape
type Circle struct {
	Radius float64
}

// NewCircle creates a new Circle with validation
func NewCircle(radius float64) (*Circle, error) {
    if radius <= 0 {
		return nil, errors.New("radius shoud be positive")
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
    return fmt.Sprintf("Circle radius: %g", c.Radius)
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
		return nil, errors.New("all sides must be positive")
	}
	if a+b <= c || a+c <= b || b+c <= a {
		return nil, errors.New("does not fulfill triangle inequality")
	}

	return &Triangle{SideA: a, SideB: b, SideC: c}, nil

}

// Area calculates the area of the triangle using Heron's formula
func (t *Triangle) Area() float64 {
    s := (t.SideA + t.SideB + t.SideC) / 2

	area := math.Sqrt(s * (s - t.SideA) * (s - t.SideB) * (s - t.SideC))

	return area
}

// Perimeter calculates the perimeter of the triangle
func (t *Triangle) Perimeter() float64 {
    return t.SideA + t.SideB + t.SideC
}

// String returns a string representation of the triangle
func (t *Triangle) String() string {
    return fmt.Sprintf("Triangle with sides %g, %g, %g", t.SideA, t.SideB, t.SideC)

}

// ShapeCalculator provides utility functions for shapes
type ShapeCalculator struct{}

// NewShapeCalculator creates a new ShapeCalculator
func NewShapeCalculator() *ShapeCalculator {
	return &ShapeCalculator{}
}

// PrintProperties prints the properties of a shape
func (sc *ShapeCalculator) PrintProperties(s Shape) {
	fmt.Printf("%s\n - Area: %.2f\n - Perimeter: %.2f\n", s, s.Area(), s.Perimeter())

}

// TotalArea calculates the sum of areas of all shapes
func (sc *ShapeCalculator) TotalArea(shapes []Shape) float64 {
	totalArea := 0.0

	for _, s := range shapes {
		totalArea += s.Area()
	}

	return totalArea

}

// LargestShape finds the shape with the largest area
func (sc *ShapeCalculator) LargestShape(shapes []Shape) Shape {
	if len(shapes) == 0 {
		return nil
	}

	largest := shapes[0]

	for _, s := range shapes {
		if s.Area() > largest.Area() {
			largest = s
		}
	}

	return largest
}

// SortByArea sorts shapes by area in ascending or descending order
func (sc *ShapeCalculator) SortByArea(shapes []Shape, ascending bool) []Shape {
    sort.Slice(shapes, func(i, j int) bool {
		if ascending {
			return shapes[i].Area() < shapes[j].Area()
		}
		return shapes[i].Area() > shapes[j].Area()
	})
	return shapes
} 