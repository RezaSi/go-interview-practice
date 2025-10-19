// Package challenge10 contains the solution for Challenge 10.
package challenge10

import (
	"fmt"
	"errors"
	"strconv"
	"math"
	"sort"
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

var ErrInvalidSideValue = errors.New("Invalid side value")
var ErrInvalidRadiusValue = errors.New("Invalid radius value")

// NewRectangle creates a new Rectangle with validation
func NewRectangle(width, height float64) (*Rectangle, error) {
    if width <= 0 || height <= 0 {
        return nil, ErrInvalidSideValue
    }
	return &Rectangle{
        Width: width,
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
    width := strconv.FormatFloat(r.Width, 'f', -1, 64)
    height := strconv.FormatFloat(r.Height, 'f', -1, 64)
	return fmt.Sprintf("Rectangle: Width = %s, Height = %s", width, height)
}

// Circle represents a perfectly round shape
type Circle struct {
	Radius float64
}

// NewCircle creates a new Circle with validation
func NewCircle(radius float64) (*Circle, error) {
    if radius <= 0 {
        return nil, ErrInvalidRadiusValue
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
	radius := strconv.FormatFloat(c.Radius, 'f', -1, 64)
	return fmt.Sprintf("Circle: Radius = %s", radius)
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
        return nil, ErrInvalidSideValue
    }
    if a >= b + c || b >= a + c || c >= a + b {
        return nil, ErrInvalidSideValue
    }
	return &Triangle{
	    SideA: a,
	    SideB: b,
	    SideC: c,
	}, nil
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
	sideA := strconv.FormatFloat(t.SideA, 'f', -1, 64)
	sideB := strconv.FormatFloat(t.SideB, 'f', -1, 64)
	sideC := strconv.FormatFloat(t.SideC, 'f', -1, 64)
	return fmt.Sprintf("Triangle: sides=%s,%s,%s", sideA, sideB, sideC)
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
	var total float64 = 0
	for _, shape := range shapes {
	    total += shape.Area()
	}
	return total
}

// LargestShape finds the shape with the largest area
func (sc *ShapeCalculator) LargestShape(shapes []Shape) Shape {
	var mmax float64 = 0
	var result Shape
	for _, shape := range shapes {
	    if shape.Area() > mmax {
	        mmax = shape.Area()
	        result = shape
	    }
	}
	return result
}

// SortByArea sorts shapes by area in ascending or descending order
func (sc *ShapeCalculator) SortByArea(shapes []Shape, ascending bool) []Shape {
	sort.Slice(shapes, func(i, j int) bool {
	    if ascending {
	        return shapes[i].Area() < shapes[j].Area()
	    } else {
	        return shapes[i].Area() > shapes[j].Area()
	    }
	})
	return shapes
} 