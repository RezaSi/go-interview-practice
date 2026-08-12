# Challenge 1: Inspecting a Wrapped Error

## The situation

Three kinds of failure, wrapped in layers of context by the time they reach you:

```go
err := fmt.Errorf("saving user: %w",
	fmt.Errorf("validating: %w", &ValidationError{Field: "email"}))
```

A type assertion sees only the outermost error, which is a `*fmt.wrapError`. To get at
the `*ValidationError` you have to walk the chain, and that is what `errors.AsType`
does.

## The task

### `FailedField`

```go
func FailedField(err error) (string, bool)
```

Return the `Field` from a `*ValidationError` anywhere in the chain, and `false` when
there is none.

### `RetryAfter`

```go
func RetryAfter(err error) (int, bool)
```

Return `RetryAfterSeconds` from a `*RateLimitError` anywhere in the chain.

### `Classify`

```go
func Classify(err error) string
```

Return the first thing you recognise: `"validation"`, `"ratelimit"`, `"notfound"` or
`"unknown"`. Check in that order.

`ErrNotFound` is a **sentinel**, a specific error value rather than a type. That one
needs a different tool.

## The two questions, and the two tools

| Question | Tool |
|---|---|
| Is this exact error value in the chain? | `errors.Is(err, ErrNotFound)` |
| Is there an error of this type, and what is it? | `errors.AsType[*ValidationError](err)` |

Reaching for `errors.Is` when you need a field off the error is the classic mistake.
`Is` gives you a yes or no about identity; only `AsType` hands you the value.

## What the tests check

- Bare errors and errors wrapped two layers deep both have to be found
- `nil` must be handled without panicking
- `RetryAfter` must not match a `ValidationError`
- `Classify` covers all four outcomes, wrapped and unwrapped
- `TestClassifyUsesIsForTheSentinel` wraps `ErrNotFound` twice: a `==` comparison fails
  this, `errors.Is` passes it

## Run it

```bash
go test -v ./...
```
