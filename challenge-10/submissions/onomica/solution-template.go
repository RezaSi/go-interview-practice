// Package challenge10 contains the solution for Challenge 10.
package challenge10

import (
	"errors"
	"fmt"
	"math"
	"sort"
)

type Shape interface {
	Area() float64
	Perimeter() float64
	fmt.Stringer
}

// Rectangle represents a four-sided shape with perpendicular sides
type Rectangle struct {
	Width  float64
	Height float64
}

// NewRectangle creates a new Rectangle with validation
func NewRectangle(width, height float64) (*Rectangle, error) {
	if width <= 0 || height <= 0 {
		return nil, errors.New("rectangle width and height must be > 0")
	}
	return &Rectangle{Width: width, Height: height}, nil
}

func (r *Rectangle) Area() float64      { return r.Width * r.Height }
func (r *Rectangle) Perimeter() float64 { return 2 * (r.Width + r.Height) }
func (r *Rectangle) String() string {
	return fmt.Sprintf("Rectangle(Width: %.2f, Height: %.2f)", r.Width, r.Height)
}

// Circle represents a perfectly round shape
type Circle struct {
	Radius float64
}

// NewCircle creates a new Circle with validation
func NewCircle(radius float64) (*Circle, error) {
	if radius <= 0 {
		return nil, errors.New("circle radius must be > 0")
	}
	return &Circle{Radius: radius}, nil
}

func (c *Circle) Area() float64      { return math.Pi * c.Radius * c.Radius }
func (c *Circle) Perimeter() float64 { return 2 * math.Pi * c.Radius }
func (c *Circle) String() string {
	return fmt.Sprintf("Circle(Radius: %.2f)", c.Radius)
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
		return nil, errors.New("triangle sides must be > 0")
	}
	// Triangle inequality: sum of any two sides must be greater than the third
	if a+b <= c || a+c <= b || b+c <= a {
		return nil, errors.New("invalid triangle: triangle inequality violated")
	}
	return &Triangle{SideA: a, SideB: b, SideC: c}, nil
}

func (t *Triangle) Area() float64 {
	s := (t.SideA + t.SideB + t.SideC) / 2
	return math.Sqrt(s * (s - t.SideA) * (s - t.SideB) * (s - t.SideC))
}
func (t *Triangle) Perimeter() float64 { return t.SideA + t.SideB + t.SideC }
func (t *Triangle) String() string {
	return fmt.Sprintf("Triangle(sides A: %.2f, B: %.2f, C: %.2f)", t.SideA, t.SideB, t.SideC)
}

// ShapeCalculator provides utility functions for shapes
type ShapeCalculator struct {
	shapes []Shape
}

// NewShapeCalculator creates a new ShapeCalculator
func NewShapeCalculator() *ShapeCalculator {
	sc := &ShapeCalculator{shapes: make([]Shape, 0, 3)}

	tr, _ := NewTriangle(3, 4, 5)
	ci, _ := NewCircle(20)
	re, _ := NewRectangle(10, 20)

	sc.shapes = append(sc.shapes, tr, ci, re)
	return sc
}

// PrintProperties prints the properties of a shape
func (sc *ShapeCalculator) PrintProperties(s Shape) {
	fmt.Printf("%s | Area: %.2f | Perimeter: %.2f\n", s.String(), s.Area(), s.Perimeter())
}

// TotalArea calculates the sum of areas of all shapes
func (sc *ShapeCalculator) TotalArea(shapes []Shape) float64 {
	var total float64
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
	largest := shapes[0]
	maxArea := largest.Area()

	for i := 1; i < len(shapes); i++ {
		a := shapes[i].Area()
		if a > maxArea {
			maxArea = a
			largest = shapes[i]
		}
	}
	return largest
}

// SortByArea sorts shapes by area in ascending or descending order
func (sc *ShapeCalculator) SortByArea(shapes []Shape, ascending bool) []Shape {
	out := make([]Shape, len(shapes))
	copy(out, shapes)

	sort.Slice(out, func(i, j int) bool {
		if ascending {
			return out[i].Area() < out[j].Area()
		}
		return out[i].Area() > out[j].Area()
	})

	return out
}
