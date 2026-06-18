package challenge6

import (
    "strings"
    "unicode"
)

func CountWordFrequency(text string) map[string]int {
    results := make(map[string]int)
    words := SplitWords(text)

    for _, w := range words {
        results[w]++
    }
    return results
}

func SplitWords(text string) []string {
    var b strings.Builder

    for _, r := range text {
        switch {
        case unicode.IsLetter(r) || unicode.IsDigit(r):
            b.WriteRune(unicode.ToLower(r))

        case r == '\'' || r == '’':
            continue

        default:
            b.WriteByte(' ')
        }
    }

    return strings.Fields(b.String())
}
