// Package challenge10 contains the solution for Challenge 10.
package challenge10

import (
	"fmt"
	"errors"
	"math"
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
	if width <= 0 || height <= 0 {
	    return nil, errors.New("Cannot create rectangle with negative width or height")
	}
	return &Rectangle{Width: width, Height: height}, nil
}

// Area calculates the area of the rectangle
func (r *Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Perimeter calculates the perimeter of the rectangle
func (r *Rectangle) Perimeter() float64 {
	return 2*(r.Width + r.Height)
}

// String returns a string representation of the rectangle
func (r *Rectangle) String() string {
	return fmt.Sprintf("Rectangle with width: %v and height: %v", r.Width, r.Height)
}

// Circle represents a perfectly round shape
type Circle struct {
	Radius float64
}

// NewCircle creates a new Circle with validation
func NewCircle(radius float64) (*Circle, error) {
	if radius <= 0 {
	    return nil, errors.New("Cannot create circle with negative radius")
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
	return fmt.Sprintf("Circle with radius: %v", c.Radius)
}

// Triangle represents a three-sided polygon
type Triangle struct {
	SideA float64
	SideB float64
	SideC float64
}

// NewTriangle creates a new Triangle with validation
func NewTriangle(a, b, c float64) (*Triangle, error) {
	if (a+b) <= c || (a+c) <= b || (b+c) <= a {
	    return nil, errors.New("Cannot create triangle with given side lengths")
	}
	return &Triangle{SideA: a, SideB: b, SideC: c}, nil
}

// Area calculates the area of the triangle using Heron's formula
func (t *Triangle) Area() float64 {
	s := (t.SideA + t.SideB + t.SideC)/ 2
	return math.Sqrt(s * (s-t.SideA) * (s-t.SideB) * (s-t.SideC))
}

// Perimeter calculates the perimeter of the triangle
func (t *Triangle) Perimeter() float64 {
	return t.SideA + t.SideB + t.SideC
}

// String returns a string representation of the triangle
func (t *Triangle) String() string {
	return fmt.Sprintf("Triangle with sides A:%v, B:%v, and C:%v", t.SideA, t.SideB, t.SideC)
}

// ShapeCalculator provides utility functions for shapes
type ShapeCalculator struct{}

// NewShapeCalculator creates a new ShapeCalculator
func NewShapeCalculator() *ShapeCalculator {
	return &ShapeCalculator{}
}

// PrintProperties prints the properties of a shape
func (sc *ShapeCalculator) PrintProperties(s Shape) {
    fmt.Printf("%v\n",s)
}

// TotalArea calculates the sum of areas of all shapes
func (sc *ShapeCalculator) TotalArea(shapes []Shape) float64 {
    var accum float64
	for _, s := range shapes {
	    accum += s.Area()
	}
	return accum
}

// LargestShape finds the shape with the largest area
func (sc *ShapeCalculator) LargestShape(shapes []Shape) Shape {
	var large Shape
	var largeArea float64
	for _, s := range shapes {
	    if s.Area() > largeArea {
	        large = s
	        largeArea = s.Area()
	    }
	}
	return large
}

// SortByArea sorts shapes by area in ascending or descending order
func (sc *ShapeCalculator) SortByArea(shapes []Shape, ascending bool) []Shape {
	length := len(shapes)
	for i := length/2 - 1; i >= 0; i--{
	    heapify(shapes, length, i)
	}
	
	for i := length - 1; i > 0; i-- {
		shapes[0], shapes[i] = shapes[i], shapes[0]
		heapify(shapes, i, 0)
	}
	if ascending{
	    return shapes
	}else{
	    slices.Reverse(shapes)
	    return shapes
	}
} 

func heapify(shapes []Shape, length int, idx int){
    largest := idx
    left, right := 2*idx + 1, 2*idx + 2
    if left < length && shapes[left].Area() > shapes[largest].Area(){
        largest = left
    }
    if right < length && shapes[right].Area() > shapes[largest].Area(){
        largest = right
    }
    if largest != idx {
        shapes[largest], shapes[idx] = shapes[idx], shapes[largest]
        heapify(shapes, length, largest)
    }
}
