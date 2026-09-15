package main

import "testing"

func TestNewConfigSetsEveryValue(t *testing.T) {
	c := NewConfig(3, true, "primary")

	if c.Retries == nil || c.Debug == nil || c.Label == nil {
		t.Fatalf("NewConfig left a field nil: %+v", c)
	}
	if *c.Retries != 3 {
		t.Errorf("Retries = %d, want 3", *c.Retries)
	}
	if !*c.Debug {
		t.Errorf("Debug = %v, want true", *c.Debug)
	}
	if *c.Label != "primary" {
		t.Errorf("Label = %q, want %q", *c.Label, "primary")
	}
}

// TestNewConfigLeavesTimeoutUnset is what the pointer fields are for: nil means
// "not specified", which is different from "specified as zero".
func TestNewConfigLeavesTimeoutUnset(t *testing.T) {
	c := NewConfig(3, true, "primary")

	if c.Timeout != nil {
		t.Fatalf("Timeout = %v, want nil so that unset stays distinguishable from zero", *c.Timeout)
	}
}

func TestNewConfigCarriesZeroValues(t *testing.T) {
	// Zero is a real value here, not an absence. A pointer to 0 must not be nil.
	c := NewConfig(0, false, "")

	if c.Retries == nil || c.Debug == nil || c.Label == nil {
		t.Fatalf("zero values must still be set, not nil: %+v", c)
	}
	if *c.Retries != 0 || *c.Debug != false || *c.Label != "" {
		t.Fatalf("zero values were not preserved: %d %v %q", *c.Retries, *c.Debug, *c.Label)
	}
}

func TestEachCallGetsItsOwnMemory(t *testing.T) {
	a := NewConfig(1, true, "a")
	b := NewConfig(2, false, "b")

	if a.Retries == nil || b.Retries == nil {
		t.Fatal("NewConfig left Retries nil")
	}
	if a.Retries == b.Retries {
		t.Fatal("two configs share the same Retries pointer; each call must allocate its own")
	}

	*a.Retries = 99
	if *b.Retries != 2 {
		t.Fatalf("writing through one config changed another: b.Retries = %d, want 2", *b.Retries)
	}
}

func TestBumpReturnsTheIncrementedValue(t *testing.T) {
	p := Bump(41)

	if p == nil {
		t.Fatal("Bump returned nil")
	}
	if *p != 42 {
		t.Fatalf("Bump(41) = %d, want 42", *p)
	}
}

// TestBumpDoesNotAliasTheArgument is the point of the exercise: new(expr) allocates
// a new variable, so the caller's value cannot be reached through the result.
func TestBumpDoesNotAliasTheArgument(t *testing.T) {
	n := 41
	p := Bump(n)
	if p == nil {
		t.Fatal("Bump returned nil")
	}

	*p = 1000

	if n != 41 {
		t.Fatalf("writing through the returned pointer changed the caller's variable: n = %d, want 41", n)
	}
}

func TestBumpAllocatesPerCall(t *testing.T) {
	a := Bump(1)
	b := Bump(1)
	if a == nil || b == nil {
		t.Fatal("Bump returned nil")
	}

	if a == b {
		t.Fatal("two calls returned the same pointer; each must allocate a new variable")
	}

	*a = 50
	if *b != 2 {
		t.Fatalf("the two results share memory: *b = %d, want 2", *b)
	}
}

func TestDescribeStillWorks(t *testing.T) {
	c := NewConfig(5, false, "x")

	if got, want := Describe(c), "retries=5 timeout=unset"; got != want {
		t.Fatalf("Describe = %q, want %q", got, want)
	}
}
