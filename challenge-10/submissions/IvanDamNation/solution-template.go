// Package challenge10 contains the solution for Challenge 10.
package challenge10

import (
    "cmp"
    "errors"
	"fmt"
	"math"
	"slices"
)

var (
    ErrLengthVal = errors.New("negative or zero value")
    ErrNanOrInfValue = errors.New("NaN or Inf value")
    ErrTriangleZero = errors.New("triangle inequality violation")
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
    if isInvalid(width) || isInvalid(height) {
        return nil, fmt.Errorf("invalid side value: %w", ErrNanOrInfValue)
    }
    
    if width <= 0 || height <= 0 {
        return nil, fmt.Errorf("side must be positive: %w", ErrLengthVal)
    }
    
	rect := &Rectangle{
	    Width: width,
	    Height: height,
	}
	
	return rect, nil
}

// Area calculates the area of the rectangle
func (r *Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Perimeter calculates the perimeter of the rectangle
func (r *Rectangle) Perimeter() float64 {
	return r.Width*2 + r.Height*2
}

// String returns a string representation of the rectangle
func (r *Rectangle) String() string {
	return fmt.Sprintf("Rectangle (width=%.2f, height=%.2f)", r.Width, r.Height)
}

// Circle represents a perfectly round shape
type Circle struct {
	Radius float64
}

// NewCircle creates a new Circle with validation
func NewCircle(radius float64) (*Circle, error) {
	if isInvalid(radius) {
        return nil, fmt.Errorf("invalid radius value: %w", ErrNanOrInfValue)
    }
    
    if radius <= 0 {
        return nil, fmt.Errorf("radius must be positive: %w", ErrLengthVal)
    }
    
    circle := &Circle{
        Radius: radius,
    }
    
	return circle, nil
}

// Area calculates the area of the circle
func (c *Circle) Area() float64 {
	return c.Radius * c.Radius * math.Pi
}

// Perimeter calculates the circumference of the circle
func (c *Circle) Perimeter() float64 {
	return c.Radius * math.Pi * 2
}

// String returns a string representation of the circle
func (c *Circle) String() string {
	return fmt.Sprintf("Circle (radius=%.2f)", c.Radius)
}

// Triangle represents a three-sided polygon
type Triangle struct {
	SideA float64
	SideB float64
	SideC float64
}

// NewTriangle creates a new Triangle with validation
func NewTriangle(a, b, c float64) (*Triangle, error) {
	if isInvalid(a) || isInvalid(b) || isInvalid(c) {
        return nil, fmt.Errorf("invalid side value: %w", ErrNanOrInfValue)
    }
    
    if a <= 0 || b <= 0 || c <= 0 {
        return nil, fmt.Errorf("side must be positive: %w", ErrLengthVal)
    }
    
    if a <= b-c || a <= c-b || b <= a-c || b <= c-a || c <= b-a || c <= a-b  {
        return nil, fmt.Errorf("invalid geometric sides %.2f, %.2f, %.2f: %w", a, b, c, ErrTriangleZero)
    }
    
	tri := &Triangle{
	    SideA: a,
	    SideB: b,
	    SideC: c,
	}
	
	return tri, nil
}

// Area calculates the area of the triangle using Heron's formula
func (t *Triangle) Area() float64 {
	s := t.Perimeter() / 2
	return math.Sqrt(s * (s - t.SideA) * (s - t.SideB) * (s - t.SideC))
}

// Perimeter calculates the perimeter of the triangle
func (t *Triangle) Perimeter() float64 {
	return t.SideA + t.SideB + t.SideC
}

// String returns a string representation of the triangle
func (t *Triangle) String() string {
	return fmt.Sprintf("Triangle (Sides=%.2f, %.2f, %.2f)", t.SideA, t.SideB, t.SideC)
}

// ShapeCalculator provides utility functions for shapes
type ShapeCalculator struct{}

// NewShapeCalculator creates a new ShapeCalculator
func NewShapeCalculator() *ShapeCalculator {
	return &ShapeCalculator{}
}

// PrintProperties prints the properties of a shape
func (sc *ShapeCalculator) PrintProperties(s Shape) {
	fmt.Println(s)
}

// TotalArea calculates the sum of areas of all shapes
func (sc *ShapeCalculator) TotalArea(shapes []Shape) float64 {
    total := 0.0
	for _, v := range shapes {
	    total += v.Area()
	}
	return total
}

// LargestShape finds the shape with the largest area
func (sc *ShapeCalculator) LargestShape(shapes []Shape) Shape {
	if len(shapes) == 0 {
	    return nil
	}
	
	if len(shapes) == 1 {
	    return shapes[0]
	}
	
	maxShape := shapes[0]
	for _, v := range shapes[1:] {
	    if v.Area() > maxShape.Area() {
	        maxShape = v
	    }
	}
	
	return maxShape
}

// SortByArea sorts shapes by area in ascending or descending order
func (sc *ShapeCalculator) SortByArea(shapes []Shape, ascending bool) []Shape {
	sorted := slices.Clone(shapes)
	slices.SortFunc(sorted, func(i, j Shape) int {
	    if ascending {
	        return cmp.Compare(i.Area(), j.Area())
	    }
	    return cmp.Compare(j.Area(), i.Area())
	})
	
	return sorted
}

func isInvalid(f float64) bool {
    return math.IsNaN(f) || math.IsInf(f, 0)
}
