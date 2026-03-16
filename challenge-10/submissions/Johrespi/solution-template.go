// Package challenge10 contains the solution for Challenge 10.
package challenge10

import (
	"errors"
	"fmt"
	"math"
	"sort"
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

func isPositive(number float64) bool {
	return number > 0.0
}

// NewRectangle creates a new Rectangle with validation
func NewRectangle(width, height float64) (*Rectangle, error) {
	// TODO: Implement validation and construction
	//
	if !isPositive(width) || !isPositive(height) {
		return nil, errors.New("not positive values")
	}

	rectangle := Rectangle{Width: width, Height: height}

	return &rectangle, nil
}

// Area calculates the area of the rectangle
func (r *Rectangle) Area() float64 {
	// TODO: Implement area calculation
	return r.Width * r.Height
}

// Perimeter calculates the perimeter of the rectangle
func (r *Rectangle) Perimeter() float64 {
	// TODO: Implement perimeter calculation
	return (2 * r.Height) + (2 * r.Width)
}

// String returns a string representation of the rectangle
func (r *Rectangle) String() string {
	// TODO: Implement string representation
	return fmt.Sprintf("Rectangle width: %.2f height: %.2f, Area: %.2f, perimeter: %.2f", r.Width, r.Height, r.Area(), r.Perimeter())
}

// Circle represents a perfectly round shape
type Circle struct {
	Radius float64
}

// NewCircle creates a new Circle with validation
func NewCircle(radius float64) (*Circle, error) {
	// TODO: Implement validation and construction
	if !isPositive(radius) {
		return nil, errors.New("Not positive")
	}
	circle := Circle{Radius: radius}
	return &circle, nil
}

// Area calculates the area of the circle
func (c *Circle) Area() float64 {
	// TODO: Implement area calculation
	return math.Pi * math.Pow(c.Radius, 2)
}

// Perimeter calculates the circumference of the circle
func (c *Circle) Perimeter() float64 {
	// TODO: Implement perimeter calculation
	return 2 * math.Pi * c.Radius
}

// String returns a string representation of the circle
func (c *Circle) String() string {
	// TODO: Implement string representation
	return fmt.Sprintf("circle radius: %.2f area: %.2f, perimeter: %.2f", c.Radius, c.Area(), c.Perimeter())
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
	if !isPositive(a) || !isPositive(b) || !isPositive(c) {
		return nil, errors.New("not positive values")
	}
	
	if a+b <= c || a+c <= b || b+c <= a {
		return nil, errors.New("invalid triangle: sum of any two sides must be greater than the third side")
	}
	
	triangle := Triangle{SideA: a, SideB: b, SideC: c}
	return &triangle, nil
}

// Area calculates the area of the triangle using Heron's formula
func (t *Triangle) Area() float64 {
	s := (t.SideA + t.SideB + t.SideC) / 2.0
	area := math.Sqrt(s * (s - t.SideA) * (s - t.SideB) * (s - t.SideC))

	return area
}

// Perimeter calculates the perimeter of the triangle
func (t *Triangle) Perimeter() float64 {
	// TODO: Implement perimeter calculation
	return t.SideA + t.SideB + t.SideC
}

// String returns a string representation of the triangle
func (t *Triangle) String() string {
	// TODO: Implement string representation
	return fmt.Sprintf("triangle sides sideA: %.2f sideB: %.2f sideC: %.2f area: %.2f, perimeter: %.2f", t.SideA, t.SideB, t.SideC, t.Area(), t.Perimeter())
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
	s.String()
}

// TotalArea calculates the sum of areas of all shapes
func (sc *ShapeCalculator) TotalArea(shapes []Shape) float64 {
	var total float64
	for _, e := range shapes {
		total += e.Area()
	}
	return total
}

// LargestShape finds the shape with the largest area
func (sc *ShapeCalculator) LargestShape(shapes []Shape) Shape {
	// TODO: Implement finding largest shape
	if len(shapes) == 0 {
		return nil
	}

	largest := shapes[0]
	for i := 1; i < len(shapes); i++ {
		if shapes[i].Area() > largest.Area() {
			largest = shapes[i]
		}
	}
	return largest
}

// SortByArea sorts shapes by area in ascending or descending order
func (sc *ShapeCalculator) SortByArea(shapes []Shape, ascending bool) []Shape {
	if len(shapes) == 0 {
		return nil
	}

	sortedShapes := make([]Shape, len(shapes))
	copy(sortedShapes, shapes)

	sort.Slice(sortedShapes, func(i, j int) bool {
		areaI := sortedShapes[i].Area()
		areaJ := sortedShapes[j].Area()

		if ascending {
			return areaI < areaJ
		}

		return areaI > areaJ
	})

	return sortedShapes
}
