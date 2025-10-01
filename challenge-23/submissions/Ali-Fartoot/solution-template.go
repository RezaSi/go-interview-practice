package main

import (
    "fmt"
)

func main() {
    testCases := []struct {
        text    string
        pattern string
    }{
        {"ABABDABACDABABCABAB", "ABABCABAB"},
        {"AABAACAADAABAABA", "AABA"},
        {"GEEKSFORGEEKS", "GEEK"},
        {"AAAAAA", "AA"},
    }

    for i, tc := range testCases {
        fmt.Printf("Test Case %d:\n", i+1)
        fmt.Printf("Text: %s\n", tc.text)
        fmt.Printf("Pattern: %s\n", tc.pattern)

        naiveResults := NaivePatternMatch(tc.text, tc.pattern)
        fmt.Printf("Naive Pattern Match: %v\n", naiveResults)

        kmpResults := KMPSearch(tc.text, tc.pattern)
        fmt.Printf("KMP Search: %v\n", kmpResults)

        rkResults := RabinKarpSearch(tc.text, tc.pattern)
        fmt.Printf("Rabin-Karp Search: %v\n", rkResults)

        fmt.Println("------------------------------")
    }
}

func NaivePatternMatch(text, pattern string) []int {
    result := make([]int, 0)
    n := len(text)
    m := len(pattern)
    
    // Empty pattern should return empty result
    if m == 0 {
        return result
    }
    
    // If text is empty or pattern is longer than text, no matches possible
    if n == 0 || n < m {
        return result
    }
    
    for i := 0; i <= n-m; i++ {
        match := true
        for j := 0; j < m; j++ {
            if text[i+j] != pattern[j] {
                match = false
                break
            }
        }
        if match {
            result = append(result, i)
        }
    }
    
    return result
}

func KMPSearch(text, pattern string) []int {
    result := make([]int, 0)
    n := len(text)
    m := len(pattern)
    
    // Empty pattern should return empty result
    if m == 0 {
        return result
    }
    
    // If text is empty or pattern is longer than text, no matches possible
    if n == 0 || n < m {
        return result
    }
    
    lps := computeLPSArray(pattern)
    
    i := 0
    j := 0
    
    for i < n {
        if text[i] == pattern[j] {
            i++
            j++
            
            if j == m {
                result = append(result, i-j)
                j = lps[j-1]
            }
        } else {
            if j != 0 {
                j = lps[j-1]
            } else {
                i++
            }
        }
    }
    
    return result
}

func computeLPSArray(pattern string) []int {
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

func RabinKarpSearch(text, pattern string) []int {
    result := make([]int, 0)
    n := len(text)
    m := len(pattern)
    
    // Empty pattern should return empty result
    if m == 0 {
        return result
    }
    
    // If text is empty or pattern is longer than text, no matches possible
    if n == 0 || n < m {
        return result
    }
    
    const d = 256
    const q = 101
    
    patternHash := 0
    textHash := 0
    h := 1
    
    for i := 0; i < m-1; i++ {
        h = (h * d) % q
    }
    
    for i := 0; i < m; i++ {
        patternHash = (d*patternHash + int(pattern[i])) % q
        textHash = (d*textHash + int(text[i])) % q
    }
    
    for i := 0; i <= n-m; i++ {
        if patternHash == textHash {
            match := true
            for j := 0; j < m; j++ {
                if text[i+j] != pattern[j] {
                    match = false
                    break
                }
            }
            if match {
                result = append(result, i)
            }
        }
        
        if i < n-m {
            textHash = (d*(textHash-int(text[i])*h) + int(text[i+m])) % q
            if textHash < 0 {
                textHash += q
            }
        }
    }
    
    return result
}