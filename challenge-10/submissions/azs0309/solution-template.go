// Package challenge10 contains the solution for Challenge 10.
package challenge10

import (
	"fmt"
	"math"
	"slices"
	"cmp"
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

// type NegativeInt float64 

// // Implement the interface
// func (n *NegativeInt) Error() string {
//     return fmt.Sprintf("Negative integer")
// }


// NewRectangle creates a new Rectangle with validation
func NewRectangle(width, height float64) (*Rectangle, error) {
	var r Rectangle
	if width <=0 || height <=0 {
	    return nil, fmt.Errorf("width %f or %f height is negative", width, height)
	}
	r.Width = width
	r.Height = height
	return &r, nil
}

// Area calculates the area of the rectangle
func (r Rectangle) Area() float64 {
	
	return r.Width*r.Height
}

// Perimeter calculates the perimeter of the rectangle
func (r Rectangle) Perimeter() float64 {
	// TODO: Implement perimeter calculation
	return 2*(r.Width+r.Height)
}

// String returns a string representation of the rectangle
func (r Rectangle) String() string {
	
	return fmt.Sprintf("Rectangle(Width: %v, Height: %v)",r.Width,r.Height)
}

// Circle represents a perfectly round shape
type Circle struct {
	Radius float64
}

// NewCircle creates a new Circle with validation
func NewCircle(radius float64) (*Circle, error) {
		var r Circle
	if radius <=0  {
	    return nil, fmt.Errorf("%f radius is negative", radius)
	}
	r.Radius = radius

	return &r, nil
}


// Area calculates the area of the circle
func (c Circle) Area() float64 {
	return math.Pi*c.Radius*c.Radius

}

// Perimeter calculates the circumference of the circle
func (c Circle) Perimeter() float64 {
	
	return 2*math.Pi*c.Radius
}

// String returns a string representation of the circle
func (c Circle) String() string {
	
	return fmt.Sprintf("Circle(Radius: %v)",c.Radius)
}

// Triangle represents a three-sided polygon
type Triangle struct {
	SideA float64
	SideB float64
	SideC float64
}

// NewTriangle creates a new Triangle with validation
func NewTriangle(a, b, c float64) (*Triangle, error) {
	var t Triangle
	if a <=0 || b <=0 || c<=0 {
    return nil, fmt.Errorf("one of the sides is negative")
}
	if a +b <= c {
    return nil, fmt.Errorf("one of the sides is negative")
}

    t.SideA = a
    t.SideB = b
    t.SideC = c
	return &t, nil
	
}

// Area calculates the area of the triangle using Heron's formula
func (t Triangle) Area() float64 {
	s := t.Perimeter()/2
 
	
	return math.Sqrt(s*(s-t.SideA)*(s-t.SideB)*(s-t.SideC))
}

// Perimeter calculates the perimeter of the triangle
func (t Triangle) Perimeter() float64 {
	// TODO: Implement perimeter calculation
	return t.SideA + t.SideB + t.SideC
}

// String returns a string representation of the triangle
func (t Triangle) String() string {
	
	 return fmt.Sprintf("Triangle(sides=%.2f, %.2f, %.2f)", t.SideA, t.SideB, t.SideC)
}

// ShapeCalculator provides utility functions for shapes
type ShapeCalculator struct{
    shapes []Shape
}

// NewShapeCalculator creates a new ShapeCalculator
func NewShapeCalculator() *ShapeCalculator {
    return &ShapeCalculator{shapes: make([]Shape,0)}
}

// PrintProperties prints the properties of a shape
func (sc *ShapeCalculator) PrintProperties(s Shape) {
 fmt.Printf("Type: %T | Properties: %+v | Area: %.2f\n", s, s, s.Area())

}

// TotalArea calculates the sum of areas of all shapes
func (sc *ShapeCalculator) TotalArea(shapes []Shape) float64 {
	if len(sc.shapes) < 0{
	    return 0
	}
	var total float64
	for _, shape := range shapes {
	    total += shape.Area()
	}
	return total
}

// LargestShape finds the shape with the largest area
func (sc *ShapeCalculator) LargestShape(shapes []Shape) Shape {
    if len(shapes) <= 0{
	    return nil
	}
	maxshape := shapes[0]
	for _, shape := range shapes[1:] {
	    if maxshape.Area() < shape.Area(){
	        maxshape = shape
	    }
	}
	return maxshape
}

// SortByArea sorts shapes by area in ascending or descending order
func (sc *ShapeCalculator) SortByArea(shapes []Shape, ascending bool) []Shape {
    sortedShapes := slices.Clone(shapes)

    slices.SortFunc(sortedShapes, func(a, b Shape) int {
        areaA := a.Area()
        areaB := b.Area()

        if ascending {
            // cmp.Compare returns -1 if a < b, 0 if equal, 1 if a > b
            return cmp.Compare(areaA, areaB)
        }
        // Reverse the comparison for descending order
        return cmp.Compare(areaB, areaA)
    })

    return sortedShapes
} 