# Hints for Greedy Coin Change

## Hint 1: Greedy Algorithm Strategy
Start with the largest denomination and use as many as possible, then move to the next largest:
```go
func makeChange(amount int, denominations []int) map[int]int {
    result := make(map[int]int)
    
    // Sort denominations in descending order
    sort.Slice(denominations, func(i, j int) bool {
        return denominations[i] > denominations[j]
    })
    
    // Use greedy approach
    for _, denom := range denominations {
        if amount >= denom {
            count := amount / denom
            result[denom] = count
            amount -= count * denom
        }
    }
    
    return result
}
```

## Hint 2: Sorting Denominations
Always sort denominations in descending order first:
```go
sort.Slice(denominations, func(i, j int) bool {
    return denominations[i] > denominations[j]
})
```

## Hint 3: Calculate Coin Count
Use integer division to find how many coins of each denomination:
```go
count := amount / denom
if count > 0 {
    result[denom] = count
    amount -= count * denom
}
```

## Hint 4: Handle Impossible Cases
Check if exact change can be made:
```go
func canMakeChange(amount int, denominations []int) bool {
    // After greedy algorithm, check if amount becomes 0
    remaining := amount
    for _, denom := range sortedDenominations {
        remaining %= denom
    }
    return remaining == 0
}
```

## Hint 5: Alternative Return Format
If you need to return total number of coins:
```go
func minCoins(amount int, denominations []int) int {
    totalCoins := 0
    for _, denom := range sortedDenominations {
        coins := amount / denom
        totalCoins += coins
        amount -= coins * denom
    }
    return totalCoins
}
```

## Hint 6: Edge Cases
Handle special cases:
```go
if amount == 0 {
    return make(map[int]int) // Empty result
}

if len(denominations) == 0 {
    return nil // Cannot make change
}
``` 