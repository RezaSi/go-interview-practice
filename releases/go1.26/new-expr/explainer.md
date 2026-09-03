## The two line dance

You want a pointer to a value. Go says no.

```go
p := &yearsSince(born)   // cannot take the address of yearsSince(born)
q := &42                 // cannot take the address of 42
```

So you invent a variable whose only job is to be addressable:

```go
age := yearsSince(born)
p := &age
```

Anyone who has built a struct with optional fields has written a small pile of these:

```go
retries := 3
debug := true
label := "primary"

cfg := Config{
	Retries: &retries,
	Debug:   &debug,
	Label:   &label,
}
```

Three variables that exist only to have addresses, sitting in scope for the rest of the
function where anything can reassign them. This is why half the Go codebases in the
world contain a `ptr` helper:

```go
func ptr[T any](v T) *T { return &v }
```

## What changed in Go 1.26

`new` now accepts an expression as well as a type.

```go
p := new(yearsSince(born))   // *int, pointing at a variable holding the result
q := new(42)                 // *int
s := new("hello")            // *string
```

The struct above becomes:

```go
cfg := Config{
	Retries: new(3),
	Debug:   new(true),
	Label:   new("primary"),
}
```

No temporaries, nothing left in scope, and no generic helper to import.

## What it actually does

The rule is small and worth stating exactly:

- `new(T)` where `T` is a **type** allocates a zero value of `T` and returns `*T`. This
  is the behaviour Go has always had.
- `new(expr)` where `expr` is a **value** allocates a new variable initialised to that
  value and returns a pointer to it. The type is the type of the expression.

The key phrase is *a new variable*. It is not a pointer into something that already
exists, so this is not a way to alias a parameter:

```go
func Bump(n int) *int {
	return new(n + 1)
}
```

`n + 1` is computed, copied into fresh memory, and you get the address of that memory.
The caller's variable is untouched. Same as the old two-liner, because the old
two-liner also copied.

## The ambiguity you might be worrying about

If `new` takes both types and values, what happens when a name is both? It cannot be:
Go has one namespace, so an identifier is either a type or a value, never both in the
same scope. The compiler already knows which one you meant.

There is one case worth knowing about. If you shadow a type name with a variable, `new`
follows the same scoping rules as everything else:

```go
type Celsius float64

func f() {
	Celsius := 20.0     // shadows the type
	p := new(Celsius)   // *float64, holding 20.0, not *Celsius holding 0
	_ = p
}
```

That is not a new hazard. Shadowing a type with a variable was already confusing and
`go vet` will tell you about it. It is just now slightly more visible.

## Where it helps most

**Optional struct fields.** The `*int` and `*bool` pattern used to signal "unset", most
often in JSON and API request types.

**Test fixtures.** Building pointer-heavy structs inline instead of ten lines of setup.

**Function results.** `new(compute())` where you needed the value once and only as a
pointer.

**Deleting your `ptr` helper.** Almost every codebase has one, usually more than one,
usually with different names in different packages.

## When not to use it

The old form is still the right one when you actually need the variable:

```go
buf := make([]byte, 0, 1024)
// ... use buf ...
p := &buf
```

And `&Config{...}` remains the idiomatic way to build a pointer to a struct literal.
`new(Config{...})` is legal but says the same thing with more ceremony.

The rule of thumb: reach for `new(expr)` when the variable would have had **no name
worth reading** and **no second use**. If you were about to call it `tmp`, this is the
feature.

> **Try it:** the challenge below builds a config struct full of optional pointer
> fields, and includes a function whose test checks that the pointer you return does
> not alias the caller's value.
