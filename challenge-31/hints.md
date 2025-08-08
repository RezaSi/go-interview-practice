# Hints for Knapsack Problem

## Hint 1: Understanding the Problem
The Knapsack Problem is a classic optimization problem. You need to select items to maximize value while staying within the weight limit. Each item can only be used once (0/1 knapsack).

## Hint 2: Dynamic Programming Approach
This problem is typically solved using dynamic programming. Consider using a 2D array where `dp[i][w]` represents the maximum value achievable using the first `i` items with capacity `w`.

## Hint 3: Base Cases

- If capacity is 0, the maximum value is 0
- If no items are available, the maximum value is 0

## Hint 4: Recursive Formula
For each item, you have two choices:
1. **Include the item**: `dp[i][w] = dp[i-1][w-weight[i]] + value[i]` (if weight[i] <= w)
2. **Exclude the item**: `dp[i][w] = dp[i-1][w]`

Choose the maximum of these two options.

## Hint 5: Implementation Strategy
1. Create a 2D slice: `dp := make([][]int, n+1)`
2. Initialize the first row and column to 0
3. Fill the dp table using the recursive formula
4. Return `dp[n][capacity]`

## Hint 6: Space Optimization
You can optimize space by using only a 1D array and filling it backwards to avoid overwriting values you still need.

## Hint 7: Edge Cases

- Handle empty items slice
- Handle zero capacity
- Ensure all weights and values are positive 
