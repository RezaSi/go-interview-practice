# Learning Materials for Inspecting a Wrapped Error

## Wrapping, and Why Assertions Stop Working

`%w` in `fmt.Errorf` produces an error that holds the original and exposes it through
`Unwrap`:

```go
err := fmt.Errorf("saving user: %w", inner)
```

The result is a `*fmt.wrapError`, not your type. So:

```go
ve, ok := err.(*ValidationError)   // false, even though one is in there
```

A type assertion looks at exactly one error. Wrapping means the one you care about is
usually not the one you are holding.

## The Chain

```
fmt.Errorf("saving user: %w", …)      ← what you have
    └─ fmt.Errorf("validating: %w", …)
           └─ &ValidationError{…}     ← what you want
```

`errors.Is`, `errors.As` and `errors.AsType` all walk this chain, checking each error in
turn and stopping at the first match. An error can also wrap several others by
implementing `Unwrap() []error`, and the search covers all branches.

## Three Functions, Three Questions

```go
errors.Is(err, ErrNotFound)                  // is this exact value in the chain?
errors.AsType[*ValidationError](err)         // is there one of this type? give it to me
errors.As(err, &ve)                          // same, into a variable I already declared
```

**`Is` is for sentinels.** Package-level values like `ErrNotFound`, `io.EOF`,
`sql.ErrNoRows`. There is no type to match, only identity, and `==` breaks the moment
someone wraps it.

**`AsType` is for typed errors carrying data.** When you need the field, the status
code, the retry delay.

The mistake worth avoiding: using `Is` when you need a value off the error. It answers
a yes/no question and cannot give you the field.

## Why AsType Exists

`errors.As` predates generics, so it reports its answer through a pointer:

```go
var ve *ValidationError
if errors.As(err, &ve) {
    use(ve.Field)
}
```

Three consequences:

- a declaration line whose only purpose is to hold the answer
- `ve` outlives the `if`, usually holding `nil`
- the target is an `any`, so `errors.As(err, ve)` compiles and panics at run time

`AsType` returns the value instead, so the variable is scoped to the `if` and the type
is checked when you build:

```go
if ve, ok := errors.AsType[*ValidationError](err); ok {
    use(ve.Field)
}
```

Same search, same result, different calling convention.

## Pointer Receivers and the Type Argument

```go
func (e *ValidationError) Error() string
```

Because `Error` has a pointer receiver, `*ValidationError` implements `error` and
`ValidationError` does not. So the type argument is `errors.AsType[*ValidationError]`.
Writing it without the pointer is a compile error, which is one of the mistakes `AsType`
catches earlier than `As` did.

## When to Keep Using errors.As

`errors.As` is not deprecated and is still the right choice when:

- you need the value after the `if` block
- the target type is decided at run time
- you are supporting Go versions before 1.26

Convert opportunistically, when you are editing the code for another reason.

## Further Reading

- [Go 1.26 release notes](https://go.dev/doc/go1.26)
- [Proposal #51945](https://github.com/golang/go/issues/51945)
- [Go blog: working with errors in Go 1.13](https://go.dev/blog/go1.13-errors)
