// Package challenge10 contains the solution for Challenge 10.
package challenge10

import (
	"fmt"
	"math"
	"errors"
	"sort"
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
	// parameters validation
	if width <= 0 || height <= 0 {
	    return nil, errors.New("width and height should be positive")
	}
	
	return &Rectangle{
		Width: width,
		Height:  height,
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
	return fmt.Sprintf("Rectangle with width %.2f and height %.2f", r.Width, r.Height)
}

// Circle represents a perfectly round shape
type Circle struct {
	Radius float64
}

// NewCircle creates a new Circle with validation
func NewCircle(radius float64) (*Circle, error) {
	// parameters validation
	if radius <= 0 {
	    return nil, errors.New("radius should be positive")
	}
	return &Circle {
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
	return fmt.Sprintf("Circle with radius %.2f", c.Radius)
}

// Triangle represents a three-sided polygon
type Triangle struct {
	SideA float64
	SideB float64
	SideC float64
}

// NewTriangle creates a new Triangle with validation
func NewTriangle(a, b, c float64) (*Triangle, error) {
	// parameters validation
	if a <= 0 || b <= 0 || c <= 0 {
	    return nil, errors.New("all sides should be positive")
	}
	// Put the values into a slice
	sides := []float64{a, b, c}
	// Sort the slice in increasing order
	sort.Float64s(sides)
	if (sides[2] >= sides[0] + sides[1]) {
	    return nil, errors.New("not a triangle")
	}
	return &Triangle{
	    SideA: a,
	    SideB: b,
	    SideC: c,
	}, nil
}

// Area calculates the area of the triangle using Heron's formula
func (t *Triangle) Area() float64 {
	// Heron's formula
    s := t.Perimeter() / 2
    return math.Sqrt(s * (s - t.SideA) * (s - t.SideB) * (s - t.SideC))
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
	fmt.Printf("Area: %.2f\n", s.Area())
    fmt.Printf("Perimeter: %.2f\n", s.Perimeter())
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
	if len(shapes) == 0 {
		return nil
	}
	
	sla := shapes[0]
	largestArea := sla.Area()
	for _, s := range shapes {
	    currentArea := s.Area()
	    if currentArea > largestArea {
	        largestArea = currentArea
	        sla = s
	    }
	}
	return sla
	
}

// SortByArea sorts shapes by area in ascending or descending order
func (sc *ShapeCalculator) SortByArea(shapes []Shape, ascending bool) []Shape {
	if len(shapes) == 0 {
		return nil
	}
	
	sortedShapes := make([]Shape, len(shapes))
	copy(sortedShapes, shapes)
	
	slices.SortFunc(sortedShapes, func(a, b Shape) int {
		areaA := a.Area()
		areaB := b.Area()

		if ascending {
			// Ascending order: lowest area first
			if areaA < areaB { return -1 }
			if areaA > areaB { return 1 }
		} else {
			// Descending order: highest area first
			if areaA > areaB { return -1 }
			if areaA < areaB { return 1 }
		}
		return 0
	})

	return sortedShapes
} 