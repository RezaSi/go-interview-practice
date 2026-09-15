package main

import "fmt"

// Config uses pointers so that "unset" and "set to the zero value" are different
type Config struct {
	Retries *int
	Debug   *bool
	Label   *string
	Timeout *int
}

// NewConfig builds a Config from plain values, with no temporary variables
func NewConfig(retries int, debug bool, label string) Config {
	// TODO: fill Retries, Debug and Label using new(expr). Leave Timeout nil.
	return Config{}
}

// Bump returns a pointer to n+1. The caller's n must not be affected by writes
// through the returned pointer.
func Bump(n int) *int {
	// TODO: return a pointer to a new variable holding n+1
	return nil
}

// Describe renders a config, showing "unset" for nil fields
func Describe(c Config) string {
	out := "retries="
	if c.Retries == nil {
		out += "unset"
	} else {
		out += fmt.Sprint(*c.Retries)
	}
	out += " timeout="
	if c.Timeout == nil {
		out += "unset"
	} else {
		out += fmt.Sprint(*c.Timeout)
	}
	return out
}

func main() {
	c := NewConfig(3, true, "primary")
	fmt.Println(Describe(c))

	n := 41
	p := Bump(n)
	fmt.Println(*p, n)
}
