package main

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
)

// FetchAll calls fetch for every id concurrently and returns the results in the
// same order as ids, whatever order the goroutines happen to finish in.
func FetchAll(ids []int, fetch func(int) string) []string {
	// TODO: start one goroutine per id with wg.Go, write each result to its own
	// index, and wait for them all
	return nil
}

// CountMatching calls match for every input concurrently and returns how many
// returned true.
func CountMatching(inputs []string, match func(string) bool) int {
	// TODO: start one goroutine per input with wg.Go and count the hits safely
	return 0
}

func main() {
	ids := []int{1, 2, 3}
	got := FetchAll(ids, func(id int) string { return fmt.Sprintf("item-%d", id) })
	fmt.Println(got)

	words := []string{"go", "rust", "gopher"}
	n := CountMatching(words, func(s string) bool { return strings.HasPrefix(s, "go") })
	fmt.Println("matches:", n)

	// unused in the template, kept so the imports stay honest
	var _ sync.WaitGroup
	var _ atomic.Int64
}
