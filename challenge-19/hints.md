# Hints for Challenge 19: Slice Operations

## Hint 1: FindMax - Iterating Through Slice
Start by handling the empty slice case and iterate to find maximum:
```go
func FindMax(numbers []int) int {
    if len(numbers) == 0 {
        return 0
    }
    
    max := numbers[0] // Initialize with first element
    for _, num := range numbers[1:] {
        if num > max {
            max = num
        }
    }
    return max
}
```

## Hint 2: RemoveDuplicates - Using Map for Tracking
Use a map to track seen values while preserving order:
```go
func RemoveDuplicates(numbers []int) []int {
    seen := make(map[int]bool)
    result := make([]int, 0, len(numbers))
    
    for _, num := range numbers {
        if !seen[num] {
            seen[num] = true
            result = append(result, num)
        }
    }
    return result
}
```

## Hint 3: ReverseSlice - Two Approaches
You can reverse by creating a new slice or in-place swapping:
```go
// Approach 1: Create new slice
func ReverseSlice(slice []int) []int {
    result := make([]int, len(slice))
    for i, val := range slice {
        result[len(slice)-1-i] = val
    }
    return result
}

// Approach 2: Copy and reverse in-place
func ReverseSlice(slice []int) []int {
    result := make([]int, len(slice))
    copy(result, slice)
    
    for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
        result[i], result[j] = result[j], result[i]
    }
    return result
}
```

## Hint 4: FilterEven - Using Modulo Operator
Filter even numbers using the modulo operator:
```go
func FilterEven(numbers []int) []int {
    var result []int
    for _, num := range numbers {
        if num%2 == 0 {
            result = append(result, num)
        }
    }
    return result
}
```

## Hint 5: Optimizing with Pre-allocation
For better performance, pre-allocate slices when you know approximate size:
```go
func FilterEven(numbers []int) []int {
    // Pre-allocate with estimated capacity
    result := make([]int, 0, len(numbers)/2)
    for _, num := range numbers {
        if num%2 == 0 {
            result = append(result, num)
        }
    }
    return result
}
```

## Hint 6: Alternative FindMax Using Built-in Functions
You can also use Go's sort package for comparison:
```go
import "sort"

func FindMax(numbers []int) int {
    if len(numbers) == 0 {
        return 0
    }
    
    // Make a copy to avoid modifying original
    temp := make([]int, len(numbers))
    copy(temp, numbers)
    sort.Ints(temp)
    
    return temp[len(temp)-1]
}
```

## Hint 7: Edge Cases to Consider
Always handle edge cases in your functions:
```go
func FindMax(numbers []int) int {
    if len(numbers) == 0 {
        return 0 // or handle as error
    }
    // Handle negative numbers correctly
    max := numbers[0] // Don't assume max is 0
    for _, num := range numbers[1:] {
        if num > max {
            max = num
        }
    }
    return max
}
```

## Key Slice Concepts:
- **Range Loops**: Use `for _, val := range slice` for iteration
- **Slice Creation**: Use `make([]int, length, capacity)` for allocation
- **Map Lookup**: Use `map[key]bool` for efficient duplicate checking
- **Slice Copying**: Use `copy(dest, src)` to avoid modifying original
- **Modulo Operator**: Use `%` to check for even/odd numbers
- **Edge Cases**: Always handle empty slices and negative numbers 