# Hints for Word Frequency Counter

## Hint 1: Data Structure Choice
Think about what data structure would be best for counting occurrences. You need something that can map words to their counts.

## Hint 2: Go Maps
Use a Go map with string keys (words) and integer values (counts). You can declare it as `make(map[string]int)`.

## Hint 3: Text Preprocessing
Before counting, you'll likely need to:
- Convert the text to lowercase for consistent counting
- Remove or replace punctuation
- Split the text into individual words

## Hint 4: String Package Functions
The `strings` package has useful functions like:
- `strings.ToLower()` for case conversion
- `strings.Fields()` for splitting by whitespace
- `strings.ReplaceAll()` for removing punctuation

## Hint 5: Regular Expressions
For more advanced text cleaning, consider using the `regexp` package to remove non-alphabetic characters.

## Hint 6: Counting Logic
For each word, check if it exists in the map. If it does, increment its count. If not, set its count to 1. You can use the comma ok idiom or simply increment (Go initializes missing map values to zero).

## Hint 7: Iteration Pattern
Use a `for range` loop to iterate through the words and update the frequency map. 