// Package challenge10 contains the solution for Challenge 10.
package challenge10

import (
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
		return nil, fmt.Errorf("zero input")
	}
	rectangle := Rectangle{
		Width: width,
		Height: height,
	}
	return &rectangle, nil
}

// Area calculates the area of the rectangle
func (r *Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Perimeter calculates the perimeter of the rectangle
func (r *Rectangle) Perimeter() float64 {
	return 2 * (r.Height + r.Width)
}

// String returns a string representation of the rectangle
func (r *Rectangle) String() string {
	return fmt.Sprintf("rectangle with width %.2f and height %.2f", r.Width, r.Height)

}

// Circle represents a perfectly round shape
type Circle struct {
	Radius float64
}

// NewCircle creates a new Circle with validation
func NewCircle(radius float64) (*Circle, error) {
	if radius <=0 {
		return nil, fmt.Errorf("zero input")
	}
	circle := Circle {
		Radius: radius,
	}

	return &circle, nil
}

// Area calculates the area of the circle
func (c *Circle) Area() float64 {
	return c.Radius * c.Radius * math.Pi
}

// Perimeter calculates the circumference of the circle
func (c *Circle) Perimeter() float64 {
	return c.Radius * 2 * math.Pi
}

// String returns a string representation of the circle
func (c *Circle) String() string {
	// TODO: Implement string representation
    return fmt.Sprintf("circle with radius %.2f", c.Radius)

}

// Triangle represents a three-sided polygon
type Triangle struct {
	SideA float64
	SideB float64
	SideC float64
}

// NewTriangle creates a new Triangle with validation
func NewTriangle(a, b, c float64) (*Triangle, error) {
	
	if a >= b + c || b >= a + c || c >= a + b {
		return nil, fmt.Errorf("inequality theorem")
	}
	if a <= 0 || b <= 0 || c<=0 {
		return nil, fmt.Errorf("zero input")
	}
	triangle := Triangle {
		SideA: a,
		SideB: b,
		SideC: c,
	}
	return &triangle, nil
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
	return fmt.Sprintf("triangle with sides %.2f, %.2f, and %.2f", t.SideA, t.SideB, t.SideC)

}

// ShapeCalculator provides utility functions for shapes
type ShapeCalculator struct{}

// NewShapeCalculator creates a new ShapeCalculator
func NewShapeCalculator() *ShapeCalculator {
	shapeCalculator := ShapeCalculator{}
	return &shapeCalculator
}

// PrintProperties prints the properties of a shape
func (sc *ShapeCalculator) PrintProperties(s Shape) {
	s.String()
}

// TotalArea calculates the sum of areas of all shapes
func (sc *ShapeCalculator) TotalArea(shapes []Shape) float64 {
	totalArea := 0.0
	for _, shape := range shapes {
		totalArea += shape.Area()
	}
	return totalArea
}

// LargestShape finds the shape with the largest area
func (sc *ShapeCalculator) LargestShape(shapes []Shape) Shape {
	latgestShape := 0.0
	var selectedShape Shape

	for _, shape := range shapes {
		if latgestShape < shape.Area(){
			latgestShape = shape.Area()
			selectedShape = shape
		}
	}
	return selectedShape
}

func (sc *ShapeCalculator) SortByArea(shapes []Shape, ascending bool) []Shape {
    sortedShapes := make([]Shape, len(shapes))
    copy(sortedShapes, shapes)
    
    n := len(sortedShapes)
    for i := 0; i < n-1; i++ {
        for j := 0; j < n-i-1; j++ {
            shouldSwap := false
            if ascending {
                shouldSwap = sortedShapes[j].Area() > sortedShapes[j+1].Area()
            } else {
                shouldSwap = sortedShapes[j].Area() < sortedShapes[j+1].Area()
            }
            
            if shouldSwap {
                sortedShapes[j], sortedShapes[j+1] = sortedShapes[j+1], sortedShapes[j]
            }
        }
    }
    
    return sortedShapes
}