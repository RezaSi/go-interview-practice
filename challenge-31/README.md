[View the Scoreboard](SCOREBOARD.md)

# Challenge 31: Knapsack Problem

## Problem Statement

Given a set of items with weights and values, find the subset of items that maximizes the total value while not exceeding a given weight limit.

## Function Signature

```go
func Knapsack(items []Item, capacity int) int
```

Where `Item` is defined as:
```go
type Item struct {
    Weight int
    Value  int
}
```

## Input Format

The input consists of:
1. A list of items, where each item has a weight and value
2. A capacity (maximum weight limit)

The input format is:

- First line: `n` (number of items)
- Next `n` lines: `weight value` (space-separated)
- Last line: `capacity`

## Output Format


- An integer representing the maximum value that can be achieved without exceeding the weight capacity.

## Constraints


- `1 <= n <= 100` (number of items)
- `1 <= weight, value <= 1000` for each item
- `1 <= capacity <= 10000`

## Sample Input and Output

### Sample Input 1

```
4
2 3
3 4
4 5
5 6
10
```

### Sample Output 1

```
13
```

**Explanation**: We can take items with weights [2, 3, 5] and values [3, 4, 6] for a total value of 13.

### Sample Input 2

```
3
1 1
2 4
3 5
4
```

### Sample Output 2

```
6
```

**Explanation**: We can take items with weights [1, 3] and values [1, 5] for a total value of 6.

## Instructions


- **Fork** the repository.
- **Clone** your fork to your local machine.
- **Create** a directory named after your GitHub username inside `challenge-31/submissions/`.
- **Copy** the `solution-template.go` file into your submission directory.
- **Implement** the `Knapsack` function.
- **Test** your solution locally by running the test file.
- **Commit** and **push** your code to your fork.
- **Create** a pull request to submit your solution.

## Testing Your Solution Locally

Run the following command in the `challenge-31/` directory:

```bash
go test -v
``` 
