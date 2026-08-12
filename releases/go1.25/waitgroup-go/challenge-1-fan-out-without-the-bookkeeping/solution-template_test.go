package main

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestFetchAllReturnsEveryResult(t *testing.T) {
	ids := []int{1, 2, 3, 4, 5}

	got := FetchAll(ids, func(id int) string { return fmt.Sprintf("item-%d", id) })

	want := []string{"item-1", "item-2", "item-3", "item-4", "item-5"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("FetchAll = %#v, want %#v", got, want)
	}
}

// TestFetchAllKeepsInputOrder fails if results are appended as they finish.
func TestFetchAllKeepsInputOrder(t *testing.T) {
	ids := []int{1, 2, 3, 4, 5}

	// The first id takes the longest, so an append-as-you-go implementation
	// returns it last.
	got := FetchAll(ids, func(id int) string {
		time.Sleep(time.Duration(6-id) * 10 * time.Millisecond)
		return fmt.Sprintf("item-%d", id)
	})

	want := []string{"item-1", "item-2", "item-3", "item-4", "item-5"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("FetchAll = %#v, want %#v\n"+
			"Results must land at their own index, not be appended in completion order.", got, want)
	}
}

func TestFetchAllHandlesAnEmptySlice(t *testing.T) {
	got := FetchAll(nil, func(int) string { return "never" })

	if len(got) != 0 {
		t.Fatalf("FetchAll(nil) = %#v, want an empty result", got)
	}
}

// TestFetchAllActuallyRunsConcurrently uses a barrier: every call blocks until all
// of them have arrived. A sequential implementation can never get past the first
// one, so it fails on the timeout instead of quietly passing.
func TestFetchAllActuallyRunsConcurrently(t *testing.T) {
	const n = 4

	var mu sync.Mutex
	arrived := 0
	all := make(chan struct{})

	done := make(chan []string, 1)
	go func() {
		done <- FetchAll([]int{1, 2, 3, 4}, func(id int) string {
			mu.Lock()
			arrived++
			if arrived == n {
				close(all)
			}
			mu.Unlock()

			<-all // wait until every goroutine has arrived
			return fmt.Sprintf("item-%d", id)
		})
	}()

	select {
	case got := <-done:
		want := []string{"item-1", "item-2", "item-3", "item-4"}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("FetchAll = %#v, want %#v", got, want)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("FetchAll did not run the fetches concurrently: all four must be in " +
			"flight at once, so start a goroutine per id rather than calling fetch in a loop")
	}
}

func TestCountMatchingCountsHits(t *testing.T) {
	words := []string{"go", "rust", "gopher", "zig", "golang"}

	got := CountMatching(words, func(s string) bool { return strings.HasPrefix(s, "go") })

	if got != 3 {
		t.Fatalf("CountMatching = %d, want 3", got)
	}
}

func TestCountMatchingOnAnEmptySlice(t *testing.T) {
	if got := CountMatching(nil, func(string) bool { return true }); got != 0 {
		t.Fatalf("CountMatching(nil) = %d, want 0", got)
	}
}

// TestCountMatchingUnderContention runs enough goroutines that an unsynchronised
// counter will lose increments.
func TestCountMatchingUnderContention(t *testing.T) {
	const n = 500
	inputs := make([]string, n)
	for i := range inputs {
		inputs[i] = "go"
	}

	got := CountMatching(inputs, func(string) bool { return true })

	if got != n {
		t.Fatalf("CountMatching = %d, want %d\n"+
			"Increments from separate goroutines need atomic or a mutex; a plain count++ loses some.", got, n)
	}
}

func TestCountMatchingWaitsForEveryGoroutine(t *testing.T) {
	inputs := []string{"a", "b", "c"}

	got := CountMatching(inputs, func(string) bool {
		time.Sleep(20 * time.Millisecond)
		return true
	})

	if got != 3 {
		t.Fatalf("CountMatching = %d, want 3; the function returned before the goroutines finished", got)
	}
}
