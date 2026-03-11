// Package challenge10 contains the solution for Challenge 10.
package challenge10

import (
	"cmp"
	"errors"
	"fmt"
	"math"
	"slices"
)

// Constants
const PI = math.Pi

// Custom errors
var (
	ErrNegativeParam = errors.New("Negative parameter")
	ErrInvalidParam  = errors.New("Sum of two sides less than third")
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
		return nil, ErrNegativeParam
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
	return fmt.Sprintf("Shape type: Rectangle\nProperties:\nWidth: %.2f\nHeight: %.2f\n", r.Width, r.Height)
}

// Circle represents a perfectly round shape
type Circle struct {
	Radius float64
}

// NewCircle creates a new Circle with validation
func NewCircle(radius float64) (*Circle, error) {
	if radius <= 0 {
		return nil, ErrNegativeParam
	}
	return &Circle{
		Radius: radius,
	}, nil
}

// Area calculates the area of the circle
func (c *Circle) Area() float64 {
	return PI * c.Radius * c.Radius
}

// Perimeter calculates the circumference of the circle
func (c *Circle) Perimeter() float64 {
	return float64(2) * PI * c.Radius
}

// String returns a string representation of the circle
func (c *Circle) String() string {
	return fmt.Sprintf("Shape type: Circle\nProperties:\nRadius: %.2f\n", c.Radius)
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
		return nil, ErrNegativeParam
	}
	if a+b <= c || b+c <= a || c+a <= b {
		return nil, ErrInvalidParam
	}
	return &Triangle{
		SideA: a,
		SideB: b,
		SideC: c,
	}, nil
}

// Area calculates the area of the triangle using Heron's formula
func (t *Triangle) Area() float64 {
	sp := (t.SideA + t.SideB + t.SideC) / float64(2)
	return math.Sqrt(sp * (sp - t.SideA) * (sp - t.SideB) * (sp - t.SideC))
}

// Perimeter calculates the perimeter of the triangle
func (t *Triangle) Perimeter() float64 {
	return t.SideA + t.SideB + t.SideC
}

// String returns a string representation of the triangle
func (t *Triangle) String() string {
	return fmt.Sprintf("Shape type: Triangle\nSides:\nSide A: %.2f\nSide B: %.2f\nSide C: %.2f\n", t.SideA, t.SideB, t.SideC)
}

// ShapeCalculator provides utility functions for shapes
type ShapeCalculator struct{}

// NewShapeCalculator creates a new ShapeCalculator
func NewShapeCalculator() *ShapeCalculator {
	return &ShapeCalculator{}
}

// PrintProperties prints the properties of a shape
func (sc *ShapeCalculator) PrintProperties(s Shape) {
	if s == nil {
		fmt.Println("No shape - No properties")
		return
	}
	fmt.Println(s.String())
}

// TotalArea calculates the sum of areas of all shapes
func (sc *ShapeCalculator) TotalArea(shapes []Shape) float64 {
	if len(shapes) == 0 {
		return 0
	}
	var res float64
	for _, s := range shapes {
		if s == nil {
			continue
		}
		res += s.Area()
	}
	return res
}

// LargestShape finds the shape with the largest area
// it returns nil for empty slice. Check resul for nil return
func (sc *ShapeCalculator) LargestShape(shapes []Shape) Shape {
	var (
		a   float64
		res Shape
	)
	for _, s := range shapes {
		if s == nil {
			continue
		}
		if a < s.Area() {
			a = s.Area()
			res = s
		}
	}
	return res
}

// SortByArea sorts shapes by area in ascending or descending order
func (sc *ShapeCalculator) SortByArea(shapes []Shape, ascending bool) []Shape {
	res := make([]Shape, len(shapes))
	copy(res, shapes)
	if ascending {
		slices.SortFunc(res, func(a, b Shape) int { return cmp.Compare(a.Area(), b.Area()) })
	} else {
		slices.SortFunc(res, func(a, b Shape) int { return cmp.Compare(b.Area(), a.Area()) })
	}
	return res
}

