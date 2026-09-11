// Package challenge10 contains the solution for Challenge 10.
package challenge10

import (
	"errors"
	"fmt"
	"math"
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
		return nil, errors.New("wrong dimensions")
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
	return 2 * (r.Width + r.Height)
}

// String returns a string representation of the rectangle
func (r *Rectangle) String() string {
	return fmt.Sprintf("Rectangle with width %0.0f and height %0.0f", r.Width, r.Height)
}

// Circle represents a perfectly round shape
type Circle struct {
	Radius float64
}

// NewCircle creates a new Circle with validation
func NewCircle(radius float64) (*Circle, error) {
	if radius <= 0 {
		return nil, errors.New("wrong dimensions")
	}
	return &Circle{Radius: radius}, nil
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
	return fmt.Sprintf("Circle with radius %0.0f", c.Radius)
}

// Triangle represents a three-sided polygon
type Triangle struct {
	SideA float64
	SideB float64
	SideC float64
}

// NewTriangle creates a new Triangle with validation
func NewTriangle(a, b, c float64) (*Triangle, error) {
	if a <= 0 || b <= 0 || c <= 0 || a >= b+c || b >= a+c || c >= a+b {
		return nil, errors.New("wrong dimensions")
	}
	return &Triangle{SideA: a, SideB: b, SideC: c}, nil
}

// Area calculates the area of the triangle using Heron's formula
func (t *Triangle) Area() float64 {
	p := t.Perimeter() / 2
	return math.Sqrt(p * (p - t.SideA) * (p - t.SideB) * (p - t.SideC))
}

// Perimeter calculates the perimeter of the triangle
func (t *Triangle) Perimeter() float64 {
	return t.SideA + t.SideB + t.SideC
}

// String returns a string representation of the triangle
func (t *Triangle) String() string {
	return fmt.Sprintf("Triangle with width sides A: %0.0f, B: %0.0f and C: %0.0f", t.SideA, t.SideB, t.SideC)
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
	var result float64
	for _, s := range shapes {
		result += s.Area()
	}
	return result
}

// LargestShape finds the shape with the largest area
func (sc *ShapeCalculator) LargestShape(shapes []Shape) Shape {
	if len(shapes) == 0 {
		return nil
	}
	largest := 0
	maxArea := shapes[0].Area()
	for i := 1; i < len(shapes); i++ {
		area := shapes[i].Area()
		if area > maxArea {
			largest = i
			maxArea = area
		}
	}

	return shapes[largest]
}

// SortByArea sorts shapes by area in ascending or descending order
func (sc *ShapeCalculator) SortByArea(shapes []Shape, ascending bool) []Shape {
	if len(shapes) <= 1 {
		return shapes
	}
	// sorted := make([]Shape, 0, len(shapes))
	var sorted []Shape

	sorted = append(sorted, shapes[0])

	for _, sh := range shapes[1:] {
		if ascending {

			inserted := false
			for i, s := range sorted {
				if s.Area() > sh.Area() {
					sorted = append(sorted, nil)
					copy(sorted[i+1:], sorted[i:])
					sorted[i] = sh
					inserted = true
					break
				}
			}
			if inserted {
				continue
			}
			sorted = append(sorted, sh)

		} else {
			inserted := false
			for i, s := range sorted {
				if s.Area() < sh.Area() {
					sorted = append(sorted, nil)
					copy(sorted[i+1:], sorted[i:])
					sorted[i] = sh
					inserted = true
					break
				}
			}
			if inserted {
				continue
			}
			sorted = append(sorted, sh)
		}
	}

	return sorted
}
