package challenge10

import (
	"cmp"
	"errors"
	"fmt"
	"math"
	"slices"
)

/*
    Challenge 10: Polymorphic Shape Calculator
	Problem Statement
	Implement a system to calculate properties of various geometric shapes using Go interfaces. This challenge focuses on understanding and correctly implementing Go's interface system to enable polymorphism.

	Requirements
	Implement a Shape interface with the following methods:

	Area() float64: Calculates the area of the shape

	Perimeter() float64: Calculates the perimeter (or circumference) of the shape

	String() string: Returns a string representation of the shape (implementing fmt.Stringer)

	Implement the following concrete shapes:

	Rectangle: Defined by Width and Height
	Circle: Defined by Radius
	Triangle: Defined by three sides (use Heron's formula for area)
	Implement a ShapeCalculator that can:

	Take any shape and return its properties

	Calculate the total area of multiple shapes
	Find the shape with the largest area from a collection
	Sort shapes by area in ascending or descending order

	Function Signatures
	// Shape interface
	type Shape interface {
		Area() float64
		Perimeter() float64
		fmt.Stringer // Includes String() string method
	}

	// Concrete types
	type Rectangle struct {
		Width, Height float64
	}

	type Circle struct {
		Radius float64
	}

	type Triangle struct {
		SideA, SideB, SideC float64
	}

	// Constructor functions
	func NewRectangle(Width, Height float64) (*Rectangle, error)
	func NewCircle(Radius float64) (*Circle, error)
	func NewTriangle(a, b, c float64) (*Triangle, error)

	// ShapeCalculator
	type ShapeCalculator struct{}

	func NewShapeCalculator() *ShapeCalculator
	func (sc *ShapeCalculator) PrintProperties(s Shape)
	func (sc *ShapeCalculator) TotalArea(shapes []Shape) float64
	func (sc *ShapeCalculator) LargestShape(shapes []Shape) Shape
	func (sc *ShapeCalculator) SortByArea(shapes []Shape, ascending bool) []Shape
	Constraints
	All measurements must be positive values
	Triangle sides must satisfy the triangle inequality theorem (sum of lengths of any two sides must exceed the length of the remaining side)
	Implement proper validation in constructors and return appropriate errors
	Use constants for π (pi) when calculating circle properties
	The String() method should return a formatted string with shape type and dimensions
	Sample Usage
	// Create shapes
	rect, _ := NewRectangle(5, 3)
	circle, _ := NewCircle(4)
	triangle, _ := NewTriangle(3, 4, 5)

	// Use shapes polymorphically
	calculator := NewShapeCalculator()
	shapes := []Shape{rect, circle, triangle}

	// Calculate total area
	totalArea := calculator.TotalArea(shapes)
	fmt.Printf("Total area: %.2f\n", totalArea)

	// Sort shapes by area
	sortedShapes := calculator.SortByArea(shapes, true)
	for _, s := range sortedShapes {
		calculator.PrintProperties(s)
	}

	// Find largest shape
	largest := calculator.LargestShape(shapes)
	fmt.Printf("Largest shape: %s with area %.2f\n", largest, largest.Area())
*/

type Shape interface {
	Area() float64
	Perimeter() float64
	String() string
}

type Rectangle struct {
	Height, Width float64
}
type Circle struct {
	Radius float64
}
type Triangle struct {
	SideA, SideB, SideC float64
}

func NewRectangle(width, height float64) (*Rectangle, error) {
	if width <= 0 {
		return nil, errors.New("width of rectangle can't be negative")
	}
	if height <= 0 {
		return nil, errors.New("height of rectangle can't be negative")
	}
	return &Rectangle{
		Width:  width,
		Height: height,
	}, nil
}

func NewCircle(radius float64) (*Circle, error) {
	if radius <= 0 {
		return nil, errors.New("radius fo circle cant be negative")
	}
	return &Circle{
		Radius: radius,
	}, nil
}

func NewTriangle(a, b, c float64) (*Triangle, error) {
	if a <= 0 {
		return nil, errors.New("side length 1 can't be negative")
	}
	if b <= 0 {
		return nil, errors.New("side length 2 can't be negative")
	}
	if c <= 0 {
		return nil, errors.New("side length 3 can't be negative")
	}
	if a+b <= c {
		return nil, errors.New("sum of two sides should be greater than third side")
	}
	return &Triangle{
		SideA: a,
		SideB: b,
		SideC: c,
	}, nil
}

func (r *Rectangle) Area() float64 {
	return r.Height * r.Width
}

func (r *Rectangle) Perimeter() float64 {
	return 2 * (r.Height + r.Width)
}
func (r *Rectangle) String() string {
	return fmt.Sprintf("Rectangle(Width: %.2f, Height: %.2f)", r.Width, r.Height)
}

func (c *Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}
func (c *Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}
func (c *Circle) String() string {
	return fmt.Sprintf("Circle(Radius: %.2f)", c.Radius)
}

func (t *Triangle) Area() float64 {
	s := (t.SideA + t.SideB + t.SideC) / 2
	ans := math.Sqrt(s * (s - t.SideA) * (s - t.SideB) * (s - t.SideC))
	return ans
}

func (t *Triangle) Perimeter() float64 {
	return t.SideA + t.SideB + t.SideC
}

func (t *Triangle) String() string {
	return fmt.Sprintf("Triangle Sides(SideA: %.2f, SideB: %.2f, SideC: %.2f)", t.SideA, t.SideB, t.SideC)
}

type ShapeCalculator struct{}

func NewShapeCalculator() *ShapeCalculator {
	return &ShapeCalculator{}
}
func (sc *ShapeCalculator) PrintProperties(s Shape) {
	fmt.Println(s.String())
}
func (sc *ShapeCalculator) TotalArea(shapes []Shape) float64 {
	ans := 0.0
	for _, val := range shapes {
		ans += val.Area()
	}
	return ans
}
func (sc *ShapeCalculator) LargestShape(shapes []Shape) Shape {
	area := 0.0
	var ans Shape
	for _, val := range shapes {
		if val.Area() > area {
			area = val.Area()
			ans = val
		}
	}
	return ans
}

func (sc *ShapeCalculator) SortByArea(shapes []Shape, ascending bool) []Shape {
	type res struct {
		val   float64
		index int
	}
	ans := make([]res, 0, len(shapes))
	for i := 0; i < len(shapes); i++ {
		ans = append(ans, res{val: shapes[i].Area(),
			index: i,
		})
	}
	slices.SortFunc(ans, func(a, b res) int {
		if ascending {
			if n := cmp.Compare(a.val, b.val); n != 0 {
				return n
			}
			return cmp.Compare(a.index, b.index)
		} else {
			if n := cmp.Compare(b.val, a.val); n != 0 {
				return n
			}
			return cmp.Compare(a.index, b.index)
		}
	})
	result := make([]Shape, 0, len(shapes))
	for i := 0; i < len(ans); i++ {
		result = append(result, shapes[ans[i].index])
	}
	return result
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	rect, _ := NewRectangle(5, 3)
	circle, _ := NewCircle(4)
	triangle, _ := NewTriangle(3, 4, 5)

	// Use shapes polymorphically
	calculator := NewShapeCalculator()
	shapes := []Shape{rect, circle, triangle}

	// Calculate total area
	totalArea := calculator.TotalArea(shapes)
	fmt.Printf("Total area: %.2f\n", totalArea)

	// Sort shapes by area
	sortedShapes := calculator.SortByArea(shapes, true)
	for _, s := range sortedShapes {
		calculator.PrintProperties(s)
	}

	// Find largest shape
	largest := calculator.LargestShape(shapes)
	fmt.Printf("Largest shape: %s with area %.2f\n", largest, largest.Area())
}
