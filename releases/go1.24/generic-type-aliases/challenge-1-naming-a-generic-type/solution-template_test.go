package main

import (
	"reflect"
	"sort"
	"testing"
)

func keys(s Set[string]) []string {
	out := make([]string, 0, len(s))
	for k := range s {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func TestNewSetHoldsEveryItem(t *testing.T) {
	s := NewSet("go", "rust", "zig")

	if got, want := keys(s), []string{"go", "rust", "zig"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("NewSet = %v, want %v", got, want)
	}
}

func TestNewSetDeduplicates(t *testing.T) {
	s := NewSet("go", "go", "go")

	if len(s) != 1 {
		t.Fatalf("NewSet with three identical items has %d entries, want 1", len(s))
	}
}

func TestNewSetWithNoItems(t *testing.T) {
	s := NewSet[string]()

	if s == nil {
		t.Fatal("NewSet with no items returned nil; return an empty set so it is usable")
	}
	if len(s) != 0 {
		t.Fatalf("NewSet() has %d entries, want 0", len(s))
	}
}

// TestAliasIsTheSameType is the property that makes this an alias rather than a
// defined type: no conversion is needed anywhere.
func TestAliasIsTheSameType(t *testing.T) {
	s := NewSet("a", "b")

	// Passing a Set[string] straight into a func that takes a plain map only
	// compiles because Set is an alias.
	if got := CountKeys(s); got != 2 {
		t.Fatalf("CountKeys = %d, want 2", got)
	}

	// And the reverse direction works too.
	var raw map[string]struct{} = s
	if len(raw) != 2 {
		t.Fatalf("assigning a Set to a map gave %d entries, want 2", len(raw))
	}

	if reflect.TypeOf(s) != reflect.TypeOf(map[string]struct{}{}) {
		t.Fatalf("Set[string] is %v, want it to be exactly map[string]struct{}", reflect.TypeOf(s))
	}
}

func TestSetSupportsBuiltins(t *testing.T) {
	// An alias for a map is a map, so every map operation is available.
	s := NewSet("a", "b", "c")

	delete(s, "b")
	if len(s) != 2 {
		t.Fatalf("after delete, len = %d, want 2", len(s))
	}

	s["d"] = struct{}{}
	if _, ok := s["d"]; !ok {
		t.Fatal("could not add to the set with plain map syntax")
	}
}

func TestContains(t *testing.T) {
	s := NewSet(1, 2, 3)

	if !Contains(s, 2) {
		t.Error("Contains(2) = false, want true")
	}
	if Contains(s, 9) {
		t.Error("Contains(9) = true, want false")
	}
}

func TestUnionCombinesBoth(t *testing.T) {
	a := NewSet("a", "b")
	b := NewSet("b", "c")

	u := Union(a, b)

	if got, want := keys(u), []string{"a", "b", "c"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Union = %v, want %v", got, want)
	}
}

func TestUnionDoesNotModifyItsInputs(t *testing.T) {
	a := NewSet("a")
	b := NewSet("b")

	_ = Union(a, b)

	if len(a) != 1 || len(b) != 1 {
		t.Fatalf("Union modified an input: len(a)=%d len(b)=%d, want 1 and 1", len(a), len(b))
	}
}

func TestZipPairsUpValues(t *testing.T) {
	got := Zip([]int{1, 2, 3}, []string{"one", "two", "three"})

	want := []Pair[int, string]{
		{First: 1, Second: "one"},
		{First: 2, Second: "two"},
		{First: 3, Second: "three"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Zip = %v, want %v", got, want)
	}
}

func TestZipStopsAtTheShorterSlice(t *testing.T) {
	got := Zip([]int{1, 2, 3}, []string{"one"})

	if len(got) != 1 {
		t.Fatalf("Zip returned %d pairs, want 1", len(got))
	}
}

// TestPairIsAnAnonymousStruct shows the same alias property for a struct type.
func TestPairIsAnAnonymousStruct(t *testing.T) {
	p := Pair[int, string]{First: 1, Second: "one"}

	var raw struct {
		First  int
		Second string
	} = p

	if raw.First != 1 || raw.Second != "one" {
		t.Fatalf("assigning a Pair to the equivalent anonymous struct lost data: %+v", raw)
	}
}
