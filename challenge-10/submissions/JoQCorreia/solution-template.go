// Package challenge10 contains the solution for Challenge 10.
package challenge10

import (
	"fmt"
	"math"
	"slices"
	"cmp"
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

	if width > 0 && height > 0 {
		r := Rectangle{Width: width, Height: height}
		return &r, nil
	}
	return nil, fmt.Errorf("Not a valid rectangle")
}

// Area calculates the area of the rectangle
func (r *Rectangle) Area() float64 {
	a := r.Height * r.Width
	return a
}

// Perimeter calculates the perimeter of the rectangle
func (r *Rectangle) Perimeter() float64 {
	p := (2 * r.Height) + (2 * r.Width)
	return p
}

// String returns a string representation of the rectangle
func (r *Rectangle) String() string {

	return fmt.Sprintf("This shape is a rectangle with %v height and %v width", r.Height, r.Width)
}

// Circle represents a perfectly round shape
type Circle struct {
	Radius float64
}

// NewCircle creates a new Circle with validation
func NewCircle(radius float64) (*Circle, error) {

	if radius > 0 {
		c := Circle{Radius: radius}
		return &c, nil
	}
	return nil, fmt.Errorf("Not a valid circle")
}

// Area calculates the area of the circle
func (c *Circle) Area() float64 {

	a := math.Pi * math.Pow(c.Radius, 2)
	return a
}

// Perimeter calculates the circumference of the circle
func (c *Circle) Perimeter() float64 {
	p := 2 * math.Pi * c.Radius
	return p
}

// String returns a string representation of the circle
func (c *Circle) String() string {

	return fmt.Sprintf("This shape is a circle with %v radius", c.Radius)
}

// Triangle represents a three-sided polygon
type Triangle struct {
	SideA float64
	SideB float64
	SideC float64
}

// NewTriangle creates a new Triangle with validation
func NewTriangle(a, b, c float64) (*Triangle, error) {

	if a > 0 && b > 0 && c > 0 {
		if (a + b) > c && (a + c) > b && (b + c) > a {
		t := &Triangle{SideA: a, SideB: b, SideC: c}
		return t, nil}
	}

	return nil, fmt.Errorf("Not a valid triangle")
}

// Area calculates the area of the triangle using Heron's formula
func (t *Triangle) Area() float64 {
	s := (t.SideA + t.SideB + t.SideC) / 2

	a := math.Sqrt(s * (s - t.SideA) * (s - t.SideB) * (s - t.SideC))

	return a
}

// Perimeter calculates the perimeter of the triangle
func (t *Triangle) Perimeter() float64 {
	p := t.SideA + t.SideB + t.SideC
	return p
}

// String returns a string representation of the triangle
func (t *Triangle) String() string {
	return fmt.Sprintf("This shape is a triangle and has sides of size %v, %v and %v", t.SideA, t.SideB, t.SideC)
}

// ShapeCalculator provides utility functions for shapes
type ShapeCalculator struct {
	shapes []Shape
}

// NewShapeCalculator creates a new ShapeCalculator
func NewShapeCalculator() *ShapeCalculator {

	return &ShapeCalculator{make([]Shape, 0)}
}

// PrintProperties prints the properties of a shape
func (sc *ShapeCalculator) PrintProperties(s Shape) {
	fmt.Print(s.Area(), s.Perimeter(), s.String())
}

// TotalArea calculates the sum of areas of all shapes
func (sc *ShapeCalculator) TotalArea(shapes []Shape) float64 {
	var tA float64
	for _, a := range shapes {
		tA = tA + a.Area()
	}
	return tA
}

// LargestShape finds the shape with the largest area
func (sc *ShapeCalculator) LargestShape(shapes []Shape) Shape {
	if len(shapes) == 0 {
		return nil
	}
	large := shapes[0]

	for i, a := range shapes {
		if large.Area() < a.Area() {
			large = shapes[i]
		}
	}

	return large
}

// SortByArea sorts shapes by area in ascending or descending order
func (sc *ShapeCalculator) SortByArea(shapes []Shape, ascending bool) []Shape {

	order := slices.Clone(shapes)

	if !ascending {
		slices.SortStableFunc(order, func(a, b Shape) int {
			return cmp.Compare(b.Area(), a.Area())

		})
		return order
	}

	slices.SortStableFunc(order, func(a, b Shape) int {
		return cmp.Compare(a.Area(), b.Area())

	})
	return order
} 