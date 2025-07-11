# Hints for Palindrome Checker

## Hint 1: Understanding Palindromes
A palindrome reads the same forwards and backwards. Examples: "racecar", "A man a plan a canal Panama" (ignoring spaces and case).

## Hint 2: String Normalization
Convert the string to lowercase and remove non-alphanumeric characters:
```go
func normalize(s string) string {
    var result strings.Builder
    for _, r := range s {
        if unicode.IsLetter(r) || unicode.IsDigit(r) {
            result.WriteRune(unicode.ToLower(r))
        }
    }
    return result.String()
}
```

## Hint 3: Two-Pointer Approach
Use two pointers from both ends moving inward:
```go
func isPalindrome(s string) bool {
    normalized := normalize(s)
    left, right := 0, len(normalized)-1
    
    for left < right {
        if normalized[left] != normalized[right] {
            return false
        }
        left++
        right--
    }
    return true
}
```

## Hint 4: Rune Conversion Alternative
For Unicode safety, convert to runes:
```go
runes := []rune(normalized)
```

## Hint 5: Simple Reverse Comparison
Alternative approach - reverse the string and compare:
```go
func reverse(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}
``` 