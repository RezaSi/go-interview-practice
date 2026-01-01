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

type Rectangle struct {
	Width  float64
	Height float64
}

func NewRectangle(width, height float64) (*Rectangle, error) {
	if width > 0 && height > 0 {
		rect := Rectangle{width, height}
		return &rect, nil
	}
	return nil, errors.New("Width and Height must be positive")
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
	if radius > 0 {
		circle := Circle{radius}
		return &circle, nil
	}
	return nil, errors.New("Radius must be positive")
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

func NewTriangle(a, b, c float64) (*Triangle, error) {
	if a > 0 && b > 0 && c > 0 && a+b > c && a+c > b && b+c > a {
		traingle := Triangle{a, b, c}
		return &traingle, nil
	}
	return nil, errors.New("All sides must be positive")
}

func (t *Triangle) Area() float64 {
	s := t.Perimeter() / 2
	return math.Sqrt(s * (s - t.SideA) * (s - t.SideB) * (s - t.SideC))
}

func (t *Triangle) Perimeter() float64 {
	return t.SideA + t.SideB + t.SideC
}

func (t *Triangle) String() string {
	return fmt.Sprintf("Triangle(sides=%.2f, %.2f, %.2f)", t.SideA, t.SideB, t.SideC)
}

type ShapeCalculator struct {
	shapes []Shape
}

func NewShapeCalculator() *ShapeCalculator {
	return &ShapeCalculator{shapes: make([]Shape, 0)}
}

func (sc *ShapeCalculator) PrintProperties(s Shape) {
	fmt.Println(s.String())
}
func (c *ShapeCalculator) TotalArea(shapes []Shape) float64 {
	total := 0.0
	for _, s := range shapes {
		total += s.Area()
	}
	return total
}
func (sc *ShapeCalculator) LargestShape(shapes []Shape) Shape {
	if len(shapes) == 0 {
		return nil
	}
	maxShapeIndex, maxArea := 0, shapes[0].Area()
	for i, s := range shapes {
		area := s.Area()
		if area > maxArea {
			maxShapeIndex, maxArea = i, area
		}
	}
	return shapes[maxShapeIndex]
}
func (sc *ShapeCalculator) SortByArea(shapes []Shape, ascending bool) []Shape {
	sortedShape := make([]Shape, len(shapes))
	copy(sortedShape, shapes)
	if ascending {
		sort.Slice(sortedShape, func(i, j int) bool {
			return sortedShape[i].Area() <= sortedShape[j].Area()
		})
	} else {
		sort.Slice(sortedShape, func(i, j int) bool {
			return sortedShape[i].Area() >= sortedShape[j].Area()
		})
	}
	return sortedShape
}

func main() {
	rect, _ := NewRectangle(5, 3)
	circle, _ := NewCircle(4)
	triangle, _ := NewTriangle(3, 4, 5)

	calculator := NewShapeCalculator()
	shapes := []Shape{rect, circle, triangle}

	totalArea := calculator.TotalArea(shapes)
	fmt.Printf("Total area: %.2f\n", totalArea)

	sortedShapes := calculator.SortByArea(shapes, true)
	for _, s := range sortedShapes {
		calculator.PrintProperties(s)
	}

	largest := calculator.LargestShape(shapes)
	fmt.Printf("Largest shape: %s with area %.2f\n", largest, largest.Area())
}

