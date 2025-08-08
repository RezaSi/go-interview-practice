# Learning Materials for Knapsack Problem

## Dynamic Programming Fundamentals

Dynamic Programming (DP) is a method for solving complex problems by breaking them down into simpler subproblems. It's particularly useful for optimization problems like the Knapsack Problem.

### Key Concepts


1. **Optimal Substructure**: The optimal solution to the problem contains optimal solutions to subproblems.
2. **Overlapping Subproblems**: The same subproblems are solved multiple times.
3. **Memoization**: Storing results of subproblems to avoid redundant calculations.

### DP Table Structure


For the Knapsack Problem, we use a 2D table where:
- Rows represent items (0 to n)
- Columns represent capacities (0 to capacity)
- `dp[i][w]` = maximum value achievable using first i items with capacity w

## Go Slices and 2D Arrays

### Creating 2D Slices

```go
// Method 1: Using make
dp := make([][]int, n+1)
for i := range dp {
    dp[i] = make([]int, capacity+1)
}

// Method 2: Using append
var dp [][]int
for i := 0; i <= n; i++ {
    row := make([]int, capacity+1)
    dp = append(dp, row)
}
```

### Accessing Elements

```go
// Access element at row i, column j
value := dp[i][j]

// Set element at row i, column j
dp[i][j] = newValue
```

## Knapsack Algorithm Implementation

### Basic Implementation

```go
func Knapsack(items []Item, capacity int) int {
    n := len(items)
    if n == 0 || capacity == 0 {
        return 0
    }
    
    // Create DP table
    dp := make([][]int, n+1)
    for i := range dp {
        dp[i] = make([]int, capacity+1)
    }
    
    // Fill the table
    for i := 1; i <= n; i++ {
        for w := 0; w <= capacity; w++ {
            // Don't include current item
            dp[i][w] = dp[i-1][w]
            
            // Include current item if it fits
            if items[i-1].Weight <= w {
                include := dp[i-1][w-items[i-1].Weight] + items[i-1].Value
                if include > dp[i][w] {
                    dp[i][w] = include
                }
            }
        }
    }
    
    return dp[n][capacity]
}
```

### Space-Optimized Implementation

```go
func KnapsackOptimized(items []Item, capacity int) int {
    n := len(items)
    if n == 0 || capacity == 0 {
        return 0
    }
    
    // Use only 1D array
    dp := make([]int, capacity+1)
    
    for i := 0; i < n; i++ {
        // Fill backwards to avoid overwriting
        for w := capacity; w >= items[i].Weight; w-- {
            include := dp[w-items[i].Weight] + items[i].Value
            if include > dp[w] {
                dp[w] = include
            }
        }
    }
    
    return dp[capacity]
}
```

## Time and Space Complexity

### Time Complexity

- **Basic DP**: O(n × capacity)
- **Space-optimized DP**: O(n × capacity)

### Space Complexity

- **Basic DP**: O(n × capacity)
- **Space-optimized DP**: O(capacity)

## Common Pitfalls

1. **Array Bounds**: Always check array bounds when accessing elements
2. **Initialization**: Properly initialize the DP table
3. **Index Offsets**: Be careful with 0-based vs 1-based indexing
4. **Memory**: For large capacities, consider space optimization

## Testing Your Solution

### Manual Testing

```go
func main() {
    items := []Item{
        {Weight: 2, Value: 3},
        {Weight: 3, Value: 4},
        {Weight: 4, Value: 5},
        {Weight: 5, Value: 6},
    }
    capacity := 10
    
    result := Knapsack(items, capacity)
    fmt.Printf("Maximum value: %d\n", result)
}
```

### Debugging Tips


1. **Print the DP table** to visualize the algorithm
2. **Use small test cases** first
3. **Check edge cases**: empty items, zero capacity
4. **Verify with known solutions**

## Advanced Topics

### Fractional Knapsack
If items can be divided (fractional knapsack), use a greedy approach by sorting by value/weight ratio.

### Multiple Constraints
For problems with multiple constraints, extend the DP table to additional dimensions.

### Backtracking
To find the actual items selected, maintain a separate table or use backtracking.

## Further Reading

- [Dynamic Programming - GeeksforGeeks](https://www.geeksforgeeks.org/dynamic-programming/)
- [Knapsack Problem - Wikipedia](https://en.wikipedia.org/wiki/Knapsack_problem)
- [Go Slices - Tour of Go](https://tour.golang.org/moretypes/8)
- [Effective Go - Arrays and Slices](https://golang.org/doc/effective_go#arrays) 
