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
	String() string
	//fmt.Stringer // Includes String() string method
}

// Rectangle represents a four-sided shape with perpendicular sides
type Rectangle struct {
	Width  float64
	Height float64
}

// NewRectangle creates a new Rectangle with validation
func NewRectangle(width, height float64) (*Rectangle, error) {
	if width <= 0 || height <= 0 {
		return nil, errors.New("Width and Height must be positive numbers")
	}

	return &Rectangle{Width: width, Height: height}, nil
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
	return fmt.Sprintf("Rectangle(width=%.2f, height=%.2f)", r.Width, r.Height)
}

// Circle represents a perfectly round shape
type Circle struct {
	Radius float64
}

// NewCircle creates a new Circle with validation
func NewCircle(radius float64) (*Circle, error) {
	if radius <= 0 {
		return nil, errors.New("Radius must be positive numbers")
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
	return fmt.Sprintf("Circle(radius=%.2f)", c.Radius)
}

// Triangle represents a three-sided polygon
type Triangle struct {
	SideA float64
	SideB float64
	SideC float64
}

// NewTriangle creates a new Triangle with validation
func NewTriangle(a, b, c float64) (*Triangle, error) {
	if a < 0 || b < 0 || c < 0 {
		return nil, errors.New("All values must be positive numbers")
	}
	if a+b <= c || a+c <= b || b+c <= a {
		return nil, errors.New("It's not a form of valid triangle")
	}
	return &Triangle{SideA: a, SideB: b, SideC: c}, nil
}

// Area calculates the area of the triangle using Heron's formula
func (t *Triangle) Area() float64 {

	perimeter := t.Perimeter() / 2

	x := perimeter * (perimeter - t.SideA) * (perimeter - t.SideB) * (perimeter - t.SideC)

	return math.Sqrt(x)
}

// Perimeter calculates the perimeter of the triangle
func (t *Triangle) Perimeter() float64 {
	return t.SideA + t.SideB + t.SideC
}

// String returns a string representation of the triangle
func (t *Triangle) String() string {
	return fmt.Sprintf("Triangle sides(a=%.2f,b=%.2f,c=%.2f)", t.SideA, t.SideB, t.SideC)
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
	var total float64 = 0.0

	for _, v := range shapes {
		total += v.Area()
	}

	return total
}

// LargestShape finds the shape with the largest area
func (sc *ShapeCalculator) LargestShape(shapes []Shape) Shape {
	var maxShape Shape = nil

	for _, v := range shapes {
		if maxShape == nil || v.Area() > maxShape.Area() {
			maxShape = v
		}
	}
	return maxShape
}

// SortByArea sorts shapes by area in ascending or descending order
func (sc *ShapeCalculator) SortByArea(shapes []Shape, ascending bool) []Shape {
	var newShapes []Shape

	for len(shapes) > 0 {
		var maxShape Shape = nil
		maxKey := 0

		for k, v := range shapes {

			if !ascending {
				if maxShape == nil || v.Area() > maxShape.Area() {
					maxShape = v
					maxKey = k
				}
			} else {
				if maxShape == nil || v.Area() < maxShape.Area() {
					maxShape = v
					maxKey = k
				}
			}

		}
		newShapes = append(newShapes, maxShape)
		shapes = append(shapes[0:maxKey], shapes[maxKey+1:len(shapes)]...)
	}
	return newShapes
}
