package main

import "fmt"

// Set is an ALIAS for a map used as a set, so Set[T] and map[T]struct{} are the
// same type and can be used interchangeably.
type Set[T comparable] = map[T]struct{}

// Pair is an alias for an anonymous struct, which saves writing it out everywhere.
type Pair[A, B any] = struct {
	First  A
	Second B
}

// NewSet builds a set from the given items
func NewSet[T comparable](items ...T) Set[T] {
	// TODO: build and return the set
	return nil
}

// Union returns a new set containing every element of a and b
func Union[T comparable](a, b Set[T]) Set[T] {
	// TODO: combine both into a new set without modifying either input
	return nil
}

// Contains reports whether the set holds v
func Contains[T comparable](s Set[T], v T) bool {
	// TODO
	return false
}

// Zip pairs up two slices, stopping at the shorter one
func Zip[A, B any](as []A, bs []B) []Pair[A, B] {
	// TODO: build the pairs
	return nil
}

// CountKeys deliberately takes a plain map. Because Set is an alias, a Set[string]
// can be passed here with no conversion. Do not change this signature.
func CountKeys(m map[string]struct{}) int {
	return len(m)
}

func main() {
	s := NewSet("go", "rust", "zig")
	fmt.Println("size:", CountKeys(s), "has go:", Contains(s, "go"))

	pairs := Zip([]int{1, 2}, []string{"one", "two"})
	fmt.Println(pairs)
}
