// Package challenge10 contains the solution for Challenge 10.
package challenge10

import (
	"fmt"
	"math"
	"errors"
	"sort"
	// Add any necessary imports here
)

const PI = math.Pi

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
	if width <= 0 || height <= 0{
	    return nil, errors.New("invalid dimentions")
	}
	var rect Rectangle = Rectangle{Width: width, Height: height}
	return &rect, nil
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
	return fmt.Sprintf("Rectangle, %v, %v width, height", r.Width ,r.Height)
}

// Circle represents a perfectly round shape
type Circle struct {
	Radius float64
}

// NewCircle creates a new Circle with validation
func NewCircle(radius float64) (*Circle, error) {
	if radius <= 0{
	    return nil, errors.New("invalid dimentions")
	}
	cir := Circle{Radius: radius}
	return &cir, nil
}

// Area calculates the area of the circle
func (c *Circle) Area() float64 {
	return PI * c.Radius * c.Radius
}

// Perimeter calculates the circumference of the circle
func (c *Circle) Perimeter() float64 {
	return 2 * PI * c.Radius
}

// String returns a string representation of the circle
func (c *Circle) String() string {
	return fmt.Sprintf("Circle, %v, radius" ,c.Radius)
}

// Triangle represents a three-sided polygon
type Triangle struct {
	SideA float64
	SideB float64
	SideC float64
}

// NewTriangle creates a new Triangle with validation
func NewTriangle(a, b, c float64) (*Triangle, error) {
	if a <= 0 || b <= 0 || c <= 0 || a + b <= c || b + c <= a || c + a <= b{
	    return nil, errors.New("invalid dimentions")
	}
	tri := Triangle{SideA: a, SideB: b, SideC: c}
	return &tri, nil
}

// Area calculates the area of the triangle using Heron's formula
func (t *Triangle) Area() float64 {
	s := (t.SideA + t.SideB + t.SideC)/2
	
	return math.Sqrt(s*(s-t.SideA)*(s-t.SideB)*(s-t.SideC))
}

// Perimeter calculates the perimeter of the triangle
func (t *Triangle) Perimeter() float64 {
	return t.SideA + t.SideB + t.SideC
}

// String returns a string representation of the triangle
func (t *Triangle) String() string {
	return fmt.Sprintf("Triangle %v, %v, %v sides", t.SideA ,t.SideB, t.SideC)
}

// ShapeCalculator provides utility functions for shapes
type ShapeCalculator struct{}

// NewShapeCalculator creates a new ShapeCalculator
func NewShapeCalculator() *ShapeCalculator {
	// TODO: Implement constructor
	return nil
}

// PrintProperties prints the properties of a shape
func (sc *ShapeCalculator) PrintProperties(s Shape) {
	fmt.Println(s.String())
}

// TotalArea calculates the sum of areas of all shapes
func (sc *ShapeCalculator) TotalArea(shapes []Shape) float64 {
    var totalArea float64 = 0.0
    for _, shape := range shapes{
        totalArea += shape.Area()
    }
	return totalArea
}

// LargestShape finds the shape with the largest area
func (sc *ShapeCalculator) LargestShape(shapes []Shape) Shape {
	var maxShape Shape
	var maxArea float64 = -1.0
    for _, shape := range shapes{
        if maxArea < shape.Area(){
            maxShape = shape
            maxArea = shape.Area()
        }
    }
	return maxShape
}

// SortByArea sorts shapes by area in ascending or descending order
func (sc *ShapeCalculator) SortByArea(shapes []Shape, ascending bool) []Shape {
    if ascending{
        sort.Slice(shapes, func(i, j int) bool{
            return (shapes[i].Area() < shapes[j].Area())
        })
    }else{
        sort.Slice(shapes, func(i, j int) bool{
            return shapes[i].Area() > shapes[j].Area()
        })
    }
	return shapes
} 