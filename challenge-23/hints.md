# Hints for Challenge 23: String Pattern Matching

## Hint 1: Naive Pattern Matching - Brute Force Approach
Start with the simplest approach that checks every position:
```go
func NaivePatternMatch(text, pattern string) []int {
    var matches []int
    
    if len(pattern) == 0 || len(pattern) > len(text) {
        return matches
    }
    
    for i := 0; i <= len(text)-len(pattern); i++ {
        match := true
        for j := 0; j < len(pattern); j++ {
            if text[i+j] != pattern[j] {
                match = false
                break
            }
        }
        if match {
            matches = append(matches, i)
        }
    }
    
    return matches
}
```

## Hint 2: KMP Algorithm - Preprocessing the Pattern
First, build the prefix function (failure function) for the pattern:
```go
func computeLPS(pattern string) []int {
    m := len(pattern)
    lps := make([]int, m)
    length := 0
    i := 1
    
    for i < m {
        if pattern[i] == pattern[length] {
            length++
            lps[i] = length
            i++
        } else {
            if length != 0 {
                length = lps[length-1]
            } else {
                lps[i] = 0
                i++
            }
        }
    }
    return lps
}
```

## Hint 3: KMP Algorithm - Main Search Function
Use the LPS array to avoid unnecessary character comparisons:
```go
func KMPSearch(text, pattern string) []int {
    var matches []int
    
    if len(pattern) == 0 || len(pattern) > len(text) {
        return matches
    }
    
    lps := computeLPS(pattern)
    i := 0 // index for text
    j := 0 // index for pattern
    
    for i < len(text) {
        if pattern[j] == text[i] {
            i++
            j++
        }
        
        if j == len(pattern) {
            matches = append(matches, i-j)
            j = lps[j-1]
        } else if i < len(text) && pattern[j] != text[i] {
            if j != 0 {
                j = lps[j-1]
            } else {
                i++
            }
        }
    }
    
    return matches
}
```

## Hint 4: Rabin-Karp Algorithm - Hash Function Setup
Use rolling hash to efficiently compare pattern and text windows:
```go
const (
    prime = 101 // A prime number for hashing
    base  = 256 // Number of characters in ASCII
)

func RabinKarpSearch(text, pattern string) []int {
    var matches []int
    
    if len(pattern) == 0 || len(pattern) > len(text) {
        return matches
    }
    
    m := len(pattern)
    n := len(text)
    
    // Calculate hash values
    patternHash := 0
    textHash := 0
    h := 1
    
    // h = base^(m-1) % prime
    for i := 0; i < m-1; i++ {
        h = (h * base) % prime
    }
    
    // Calculate initial hash values
    for i := 0; i < m; i++ {
        patternHash = (base*patternHash + int(pattern[i])) % prime
        textHash = (base*textHash + int(text[i])) % prime
    }
    
    // Rest of implementation...
}
```

## Hint 5: Rabin-Karp Algorithm - Rolling Hash Implementation
Implement the rolling hash technique for efficient window sliding:
```go
func RabinKarpSearch(text, pattern string) []int {
    var matches []int
    
    if len(pattern) == 0 || len(pattern) > len(text) {
        return matches
    }
    
    m := len(pattern)
    n := len(text)
    
    patternHash := 0
    textHash := 0
    h := 1
    
    // Calculate h = base^(m-1) % prime
    for i := 0; i < m-1; i++ {
        h = (h * base) % prime
    }
    
    // Calculate initial hash values
    for i := 0; i < m; i++ {
        patternHash = (base*patternHash + int(pattern[i])) % prime
        textHash = (base*textHash + int(text[i])) % prime
    }
    
    // Slide the pattern over text one by one
    for i := 0; i <= n-m; i++ {
        // Check if hash values match
        if patternHash == textHash {
            // Check characters one by one to avoid false positives
            match := true
            for j := 0; j < m; j++ {
                if text[i+j] != pattern[j] {
                    match = false
                    break
                }
            }
            if match {
                matches = append(matches, i)
            }
        }
        
        // Calculate hash for next window
        if i < n-m {
            textHash = (base*(textHash-int(text[i])*h) + int(text[i+m])) % prime
            
            // Handle negative hash values
            if textHash < 0 {
                textHash = textHash + prime
            }
        }
    }
    
    return matches
}
```

## Hint 6: Edge Cases and Input Validation
Handle edge cases properly in all algorithms:
```go
func handleEdgeCases(text, pattern string) ([]int, bool) {
    // Empty pattern
    if len(pattern) == 0 {
        return []int{}, true
    }
    
    // Pattern longer than text
    if len(pattern) > len(text) {
        return []int{}, true
    }
    
    // Empty text but non-empty pattern
    if len(text) == 0 {
        return []int{}, true
    }
    
    // Continue with normal processing
    return nil, false
}
```

## Hint 7: Algorithm Optimization Tips
Consider these optimizations for better performance:
```go
// For KMP: Optimize the LPS computation
func computeLPSOptimized(pattern string) []int {
    m := len(pattern)
    lps := make([]int, m)
    
    for i, length := 1, 0; i < m; {
        if pattern[i] == pattern[length] {
            length++
            lps[i] = length
            i++
        } else if length != 0 {
            length = lps[length-1]
        } else {
            lps[i] = 0
            i++
        }
    }
    return lps
}

// For Rabin-Karp: Use better hash function to reduce collisions
func betterHash(s string) uint64 {
    var hash uint64 = 0
    for i := 0; i < len(s); i++ {
        hash = hash*31 + uint64(s[i])
    }
    return hash
}
```

## Hint 8: Performance Comparison and Testing
Create benchmarks to compare algorithm performance:
```go
import "testing"

func BenchmarkNaive(b *testing.B) {
    text := strings.Repeat("ABCDEFGH", 1000)
    pattern := "DEFG"
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        NaivePatternMatch(text, pattern)
    }
}

func BenchmarkKMP(b *testing.B) {
    text := strings.Repeat("ABCDEFGH", 1000)
    pattern := "DEFG"
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        KMPSearch(text, pattern)
    }
}
```

## Key Pattern Matching Concepts:
- **Naive Approach**: Simple but O(n*m) time complexity
- **KMP Algorithm**: Uses prefix function to avoid redundant comparisons
- **Rolling Hash**: Rabin-Karp uses hash values for quick comparisons
- **LPS Array**: Longest Proper Prefix which is also Suffix
- **Hash Collisions**: Rabin-Karp needs character-by-character verification
- **Edge Cases**: Handle empty strings and boundary conditions
- **Time Complexity**: Understanding when each algorithm performs best 