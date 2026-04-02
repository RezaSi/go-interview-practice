// Package challenge10 contains the solution for Challenge 10.
package challenge10

import (
    "cmp"
	"fmt"
	"math"
	"errors"
	"slices"
)

var (
    ErrNegativeDimension = errors.New("one or more dimensions is negative")
    ErrInavlidDemension = errors.New("dimensions do not make a valid desired shape")
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
        return nil, ErrNegativeDimension
    }
    
	return &Rectangle{Width: width, Height: height}, nil
}

// Area calculates the area of the rectangle
func (r *Rectangle) Area() float64 {
	return r.Width*r.Height
}

// Perimeter calculates the perimeter of the rectangle
func (r *Rectangle) Perimeter() float64 {
	return 2 * (r.Width+r.Height)
}

// String returns a string representation of the rectangle
func (r *Rectangle) String() string {
	return fmt.Sprintf("Rectangle{width: %v, height: %v", r.Width, r.Height)
}

// Circle represents a perfectly round shape
type Circle struct {
	Radius float64
}

// NewCircle creates a new Circle with validation
func NewCircle(radius float64) (*Circle, error) {
    if radius <= 0 {
        return nil, ErrNegativeDimension
    }
    
	return &Circle{Radius: radius}, nil
}

// Area calculates the area of the circle
func (c *Circle) Area() float64 {
	return math.Pi * math.Pow(c.Radius, 2)
}

// Perimeter calculates the circumference of the circle
func (c *Circle) Perimeter() float64 {
	return 2*math.Pi*c.Radius
}

// String returns a string representation of the circle
func (c *Circle) String() string {
	return fmt.Sprintf("Circle{radius: %v}", c.Radius)
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
        return nil, ErrNegativeDimension
    }
    
    if a + b <= c {
        return nil, ErrInavlidDemension
    }
    
	return &Triangle{SideA: a, SideB: b, SideC: c}, nil
}

// Area calculates the area of the triangle using Heron's formula
func (t *Triangle) Area() float64 {
    s := t.Perimeter()/2
	return math.Sqrt(s*(s-t.SideA)*(s-t.SideB)*(s-t.SideC))
}

// Perimeter calculates the perimeter of the triangle
func (t *Triangle) Perimeter() float64 {
	return t.SideA+t.SideB+t.SideC
}

// String returns a string representation of the triangle
func (t *Triangle) String() string {
	return fmt.Sprintf("Triangle{sides: %v, %v, %v}", t.SideA, t.SideB, t.SideC)
}

// ShapeCalculator provides utility functions for shapes
type ShapeCalculator struct{}

// NewShapeCalculator creates a new ShapeCalculator
func NewShapeCalculator() *ShapeCalculator {
	return &ShapeCalculator{}
}

// PrintProperties prints the properties of a shape
func (sc *ShapeCalculator) PrintProperties(s Shape) {
    fmt.Printf("%v\n", s)
}

// TotalArea calculates the sum of areas of all shapes
func (sc *ShapeCalculator) TotalArea(shapes []Shape) float64 {
    var sum float64
    
    for _, s := range shapes {
        sum += s.Area()
    }
    
	return sum
}

// LargestShape finds the shape with the largest area
func (sc *ShapeCalculator) LargestShape(shapes []Shape) Shape {
    var max float64
    var maxShape Shape
    
    for _, s := range shapes {
        if max < s.Area() {
            max = s.Area()
            maxShape = s
        }
    }
    
	return maxShape
}

// SortByArea sorts shapes by area in ascending or descending order
func (sc *ShapeCalculator) SortByArea(shapes []Shape, ascending bool) []Shape {
    slices.SortFunc(shapes, func(a, b Shape) int {
        comp := cmp.Compare(a.Area(), b.Area())
        
        if !ascending {
            comp = -comp
        }
        
        return comp
    })
    
	return shapes
} 
