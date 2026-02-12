// Package challenge10 contains the solution for Challenge 10.
package challenge10

import (
    "fmt"
	"cmp"
	"errors"
	"math"
	"slices"
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
	switch {
	case width <= 0:
		return nil, errors.New("width must be > 0")
	case height <= 0:
		return nil, errors.New("height must be > 0")
	}

	return &Rectangle{
		Width:  width,
		Height: height}, nil
}

// Area calculates the area of the rectangle
func (r *Rectangle) Area() float64 {
	// TODO: Implement area calculation
	return r.Height * r.Width
}

// Perimeter calculates the perimeter of the rectangle
func (r *Rectangle) Perimeter() float64 {
	// TODO: Implement perimeter calculation
	return (r.Height + r.Width) * 2
}

// String returns a string representation of the rectangle
func (r *Rectangle) String() string {
	// TODO: Implement string representation
	return fmt.Sprintf("A rectangle consisting of %f widths and %f heights.", r.Width, r.Height)
}

// Circle represents a perfectly round shape
type Circle struct {
	Radius float64
}

// NewCircle creates a new Circle with validation
func NewCircle(radius float64) (*Circle, error) {
	// TODO: Implement validation and construction
	if radius <= 0 {
		return nil, errors.New("radius must be > 0")
	}
	return &Circle{
		Radius: radius,
	}, nil
}

// Area calculates the area of the circle
func (c *Circle) Area() float64 {
	// TODO: Implement area calculation
	return c.Radius * c.Radius * math.Pi
}

// Perimeter calculates the circumference of the circle
func (c *Circle) Perimeter() float64 {
	// TODO: Implement perimeter calculation
	return c.Radius * 2 * math.Pi
}

// String returns a string representation of the circle
func (c *Circle) String() string {
	// TODO: Implement string representation
	return fmt.Sprintf("A circle consisting of %f radius", c.Radius)
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
	if a <= 0 || b <= 0 || c <= 0 {
		return nil, errors.New("side must be > 0")
	}
	
	if a+b <= c || a+c <= b || b+c <= a {
	    return nil, errors.New("invalid triangle side")
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
	return fmt.Sprintf("A triangle consisting of three sides: %f, %f, %f", t.SideA, t.SideB, t.SideC)
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
func (sc *ShapeCalculator) TotalArea(shapes []Shape) float64 {
	// TODO: Implement total area calculation
	sum := 0.0
	ch := make(chan float64, len(shapes))

	for _, s := range shapes {
		go func(sh Shape) {
			ch <- sh.Area()
		}(s)
	}

	for i := 0; i < len(shapes); i++ {
		sum += <-ch
	}

	return sum
}

// LargestShape finds the shape with the largest area
func (sc *ShapeCalculator) LargestShape(shapes []Shape) Shape {
	// TODO: Implement finding largest shape

	type result struct {
		shape Shape
		area  float64
	}

	ch := make(chan result, len(shapes))
	sl := []result{}

	for _, s := range shapes {
		go func(sh Shape) {
			r := result{
				shape: sh,
				area:  sh.Area(),
			}
			ch <- r
		}(s)
	}

	for i := 0; i < len(shapes); i++ {
		sl = append(sl, <-ch)
	}

	slices.SortFunc(sl, func(a, b result) int {
		return cmp.Compare(a.area, b.area)
	})

	return sl[len(sl)-1].shape
}

// SortByArea sorts shapes by area in ascending or descending order
func (sc *ShapeCalculator) SortByArea(shapes []Shape, ascending bool) []Shape {
	// TODO: Implement sorting shapes by area
	type result struct {
		shape Shape
		area  float64
	}

	ch := make(chan result, len(shapes))
	sl := []result{}
	rv := []Shape{}

	for _, s := range shapes {
		go func(sh Shape) {
			r := result{
				shape: sh,
				area:  sh.Area(),
			}
			ch <- r
		}(s)
	}

	for i := 0; i < len(shapes); i++ {
		sl = append(sl, <-ch)
	}

	slices.SortFunc(sl, func(a, b result) int {
		return cmp.Compare(a.area, b.area)
	})

	for _, s := range sl {
		rv = append(rv, s.shape)
	}

	if !ascending {
		slices.Reverse(rv)
	}

	return rv
}