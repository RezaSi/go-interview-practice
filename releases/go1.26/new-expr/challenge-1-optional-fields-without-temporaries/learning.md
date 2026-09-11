# Learning Materials for Optional Fields Without Temporaries

## Addressability, and Why It Was in the Way

Go's `&` operator needs an **addressable** operand. Variables, slice elements, struct
fields and pointer dereferences are addressable. Function results, arithmetic and
literals are not, because there is no variable to point at:

```go
p := &f()      // invalid
q := &(a + b)  // invalid
r := &42       // invalid
```

Composite literals are the special case that makes `&Config{}` work, and their being
special is exactly why the rest felt inconsistent.

The workaround created a variable so there was something to address:

```go
v := f()
p := &v
```

## What new(expr) Does

```go
new(T)      // allocate a zero value of type T, return *T
new(expr)   // allocate a new variable initialised to expr, return a pointer to it
```

The type of `new(expr)` is the type of the expression, so `new(3)` is `*int` and
`new("hi")` is `*string`.

The important word is **new variable**. It allocates; it does not alias. This is the
same behaviour the two-line version had, which is why converting old code is
mechanical.

## Pointers as Optional Fields

Why the struct in this challenge uses pointers at all:

```go
type Config struct {
    Retries *int
}
```

With a plain `int`, `0` is ambiguous: it might be "the user asked for zero retries" or
"nobody said". With `*int`, `nil` means unset and `new(0)` means the user chose zero.

This shows up constantly in JSON request bodies, where a field being absent and a field
being `0` mean different things, and in configuration merging, where you need to know
which values to overwrite.

The cost is that every read needs a nil check, which is why `Describe` in this challenge
looks the way it does. Use the pattern when the distinction matters, not by default.

## Where the Value Lives

`new` says nothing about the heap. Go decides via escape analysis: if the pointer
outlives the function it is heap-allocated, otherwise it can stay on the stack. That was
already true of `&v` on a local variable, so `new(expr)` changes nothing about
performance. You can see the decision with:

```bash
go build -gcflags='-m' ./...
```

## When Not to Reach for It

- **You need the variable again.** If it has a meaningful name and a second use, write
  the variable.
- **Struct literals.** `&Config{Retries: new(3)}` is idiomatic. `new(Config{...})` is
  legal but says the same thing less directly.
- **You are shadowing a type name.** `new(Celsius)` where `Celsius` is a local variable
  gives you a pointer to that value, not to the type's zero value. Not a new hazard, but
  worth knowing it now applies to `new` too.

The rule of thumb: if the variable would have been called `tmp`, `v` or `p`, this is the
feature.

## Further Reading

- [Go 1.26 release notes](https://go.dev/doc/go1.26)
- [Go spec: allocation](https://go.dev/ref/spec#Allocation)
- [Go spec: address operators](https://go.dev/ref/spec#Address_operators)
