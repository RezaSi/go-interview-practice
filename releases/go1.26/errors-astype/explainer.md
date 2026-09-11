## Three lines to ask one question

You have an error. You want to know whether anything in its chain is a `*fs.PathError`,
and if so, which file it was about.

```go
var pe *fs.PathError
if errors.As(err, &pe) {
	log.Printf("failed on %s", pe.Path)
}
```

That works, and it has been the shape since Go 1.13. Three things about it grate.

**The variable is declared before the question is asked.** You write `var pe
*fs.PathError` on a line that does nothing except make somewhere for the answer to go.

**It leaks out of the `if`.** After the block, `pe` is still in scope, still typed, and
almost certainly `nil`. Any later use is a nil dereference waiting to happen, and
nothing warns you.

**The target is checked at run time.** `errors.As` takes an `any`, so this compiles:

```go
errors.As(err, pe)   // forgot the &
```

and panics when it runs: *errors: target must be a non-nil pointer*. The compiler had
every piece of information needed to reject it and no way to express that.

## What Go 1.26 added

```go
func AsType[T error](err error) (T, bool)
```

The same question, as an expression:

```go
if pe, ok := errors.AsType[*fs.PathError](err); ok {
	log.Printf("failed on %s", pe.Path)
}
```

`pe` is scoped to the `if`. There is no target to get wrong. The type you are asking
about is a type argument, so a mistake is a compile error rather than a panic.

Verified against a real toolchain:

```go
_, err := os.Open("/definitely/missing")
if pe, ok := errors.AsType[*os.PathError](err); ok {
	fmt.Println("path error on:", pe.Path, "op:", pe.Op)
}
```

```text
path error on: /definitely/missing op: open
```

## It is the same search

`AsType` is not a different algorithm. Both walk the chain produced by `Unwrap`,
checking each error in turn for one assignable to the target type, and both stop at the
first match. Anything you know about `errors.As` still applies:

- `%w` in `fmt.Errorf` is what puts an error in the chain.
- An error can wrap several others (`Unwrap() []error`), and the search covers all of
  them.
- A type implementing `As(any) bool` can customise matching, and `AsType` honours it.

If you have a mental model of `errors.As`, you do not need a new one. Only the calling
convention changed.

## The three of them, side by side

Go now has three ways to interrogate an error, and it is worth being able to pick
instantly:

| Question | Use |
|---|---|
| Is this specific sentinel value in the chain? | `errors.Is(err, fs.ErrNotExist)` |
| Is there an error of this type, and what is it? | `errors.AsType[*fs.PathError](err)` |
| Same, but into an existing variable | `errors.As(err, &pe)` |

`Is` compares values, `AsType` and `As` match types. The usual mistake is reaching for
`Is` when you want a field off the error, and only `As`/`AsType` can give you that.

## Should you convert existing code?

There is no rush. `errors.As` is not deprecated and will not be: it is in the Go 1
compatibility promise, and it remains the right call when you need the result outside
the `if`, or when the target type is only known dynamically.

Convert when you touch the code anyway. The conversion is mechanical:

```go
var target *MyError
if errors.As(err, &target) {

// becomes

if target, ok := errors.AsType[*MyError](err); ok {
```

The one thing to watch is a `var` declaration that was being used after the block. If
it was, keep `errors.As` or hoist the value out deliberately.

> **Try it:** the challenge below has a wrapped error chain and asks you to pull
> structured detail out of it with `AsType`, including the case where the chain contains
> no match at all.
