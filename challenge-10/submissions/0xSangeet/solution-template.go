package challenge10

import (
	"fmt"
	"math"
	"errors"
	"sort"
)

type Shape interface {
	Area() float64
	Perimeter() float64
	fmt.Stringer
}

type Rectangle struct {
	Width  float64
	Height float64
}

func NewRectangle(width, height float64) (*Rectangle, error) {
    if width <= 0 || height <= 0 {
        return nil, errors.New("cant have zero or negative values")
    }
	return &Rectangle{width, height}, nil
}

func (r *Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r *Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

func (r *Rectangle) String() string {
	return fmt.Sprintf("Rectangle(width=%.2f, height=%.2f)", r.Width, r.Height)
}

type Circle struct {
	Radius float64
}

func NewCircle(radius float64) (*Circle, error) {
	if radius <= 0 {
	    return nil, errors.New("cant have zero or negative values")
	}
	return &Circle{radius}, nil
}

func (c *Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c *Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

func (c *Circle) String() string {
	return fmt.Sprintf("Circle(radius=%.2f)", c.Radius)
}

type Triangle struct {
	SideA float64
	SideB float64
	SideC float64
}

// NewTriangle creates a new Triangle with validation
func NewTriangle(a, b, c float64) (*Triangle, error) {
	
	if (a + b > c && a + c > b && b + c > a) && a > 0 && b > 0 && c > 0 {
	    return &Triangle{a, b, c}, nil
	}
	return nil, errors.New("invalid triangle")
}

func (t *Triangle) Area() float64 {
	s := (t.SideA + t.SideB + t.SideC) / 2
	area := math.Sqrt(s * (s - t.SideA) * (s - t.SideB) * (s - t.SideC))
	return area
}

func (t *Triangle) Perimeter() float64 {
	return t.SideA + t.SideB + t.SideC
}

func (t *Triangle) String() string {
    return fmt.Sprintf("Triangle(sides=%.2f, %.2f, %.2f)", t.SideA, t.SideB, t.SideC)
}

type ShapeCalculator struct{}

func NewShapeCalculator() *ShapeCalculator {
	return &ShapeCalculator{}
}

// PrintProperties prints the properties of a shape
func (sc *ShapeCalculator) PrintProperties(s Shape) {
	fmt.Println(s.String())
}

// TotalArea calculates the sum of areas of all shapes
func (sc *ShapeCalculator) TotalArea(shapes []Shape) float64 {
	var totalArea float64
	
	for _, shape := range shapes {
	    totalArea += shape.Area()
	}
	
	return totalArea
}

func (sc *ShapeCalculator) LargestShape(shapes []Shape) Shape {
    if len(shapes) == 0 {
        return nil
    }
    
	maxS := shapes[0]
	maxA := shapes[0].Area()
	
	for _, shape := range shapes {
	    if currA := shape.Area(); currA > maxA {
	        maxA = currA
	        maxS = shape
	    }
	}
	
	return maxS
}

// SortByArea sorts shapes by area in ascending or descending order
func (sc *ShapeCalculator) SortByArea(shapes []Shape, ascending bool) []Shape {
	sorted := make([]Shape, len(shapes))
    copy(sorted, shapes)
    sort.Slice(sorted, func(i, j int) bool {
        if ascending {
            return sorted[i].Area() < sorted[j].Area()
        }
        return sorted[i].Area() > sorted[j].Area()
    })
    return sorted
} 