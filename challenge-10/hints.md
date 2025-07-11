# Hints for Polymorphic Shape Calculator

## Hint 1: Interface Definition
Define the Shape interface with required methods:
```go
type Shape interface {
    Area() float64
    Perimeter() float64
    fmt.Stringer // Embeds String() string
}
```

## Hint 2: Rectangle Implementation
Implement the Rectangle struct and its methods:
```go
type Rectangle struct {
    Width, Height float64
}

func (r *Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r *Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}
```

## Hint 3: Circle Implementation
For circle calculations, use `math.Pi`:
```go
func (c *Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

func (c *Circle) Perimeter() float64 {
    return 2 * math.Pi * c.Radius
}
```

## Hint 4: Triangle Area with Heron's Formula
Implement triangle area using Heron's formula:
```go
func (t *Triangle) Area() float64 {
    s := (t.SideA + t.SideB + t.SideC) / 2 // semi-perimeter
    return math.Sqrt(s * (s - t.SideA) * (s - t.SideB) * (s - t.SideC))
}
```

## Hint 5: Constructor Validation
Validate inputs in constructor functions:
```go
func NewTriangle(a, b, c float64) (*Triangle, error) {
    if a <= 0 || b <= 0 || c <= 0 {
        return nil, errors.New("sides must be positive")
    }
    
    // Triangle inequality: sum of any two sides > third side
    if a+b <= c || a+c <= b || b+c <= a {
        return nil, errors.New("sides do not form a valid triangle")
    }
    
    return &Triangle{SideA: a, SideB: b, SideC: c}, nil
}
```

## Hint 6: String Method Implementation
Implement the String method for each shape:
```go
func (r *Rectangle) String() string {
    return fmt.Sprintf("Rectangle(width=%.2f, height=%.2f)", r.Width, r.Height)
}

func (c *Circle) String() string {
    return fmt.Sprintf("Circle(radius=%.2f)", c.Radius)
}
```

## Hint 7: Total Area Calculation
Iterate through shapes and sum their areas:
```go
func (sc *ShapeCalculator) TotalArea(shapes []Shape) float64 {
    var total float64
    for _, shape := range shapes {
        total += shape.Area()
    }
    return total
}
```

## Hint 8: Finding Largest Shape
Compare areas to find the largest:
```go
func (sc *ShapeCalculator) LargestShape(shapes []Shape) Shape {
    if len(shapes) == 0 {
        return nil
    }
    
    largest := shapes[0]
    for _, shape := range shapes[1:] {
        if shape.Area() > largest.Area() {
            largest = shape
        }
    }
    return largest
}
```

## Hint 9: Sorting by Area
Use `sort.Slice` to sort shapes by area:
```go
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
```

## Hint 10: PrintProperties Method
Print shape information using the interface:
```go
func (sc *ShapeCalculator) PrintProperties(s Shape) {
    fmt.Printf("Shape: %s\n", s)
    fmt.Printf("Area: %.2f\n", s.Area())
    fmt.Printf("Perimeter: %.2f\n", s.Perimeter())
}
``` 