# Hints for Inspecting a Wrapped Error

## Hint 1: The Shape of AsType

The type you are looking for is a type argument, and the result comes back as a value:

```go
if ve, ok := errors.AsType[*ValidationError](err); ok {
    // ve is a *ValidationError
}
```

Note the pointer in the type argument. `Error()` is defined on `*ValidationError`, so
that is the type that implements `error`.

## Hint 2: FailedField Is That Pattern Plus a Return

```go
if ve, ok := errors.AsType[*ValidationError](err); ok {
    return ve.Field, true
}
return "", false
```

## Hint 3: You Do Not Always Need the Value

`Classify` only cares whether a match exists, so discard it:

```go
if _, ok := errors.AsType[*ValidationError](err); ok {
    return "validation"
}
```

## Hint 4: The Sentinel Is Different

`ErrNotFound` is a value created by `errors.New`. There is no type to match on, and
comparing with `==` fails as soon as it is wrapped. Use:

```go
if errors.Is(err, ErrNotFound) {
    return "notfound"
}
```

## Hint 5: nil Needs No Special Case

`errors.AsType(nil)` returns `false` and `errors.Is(nil, target)` returns `false` for a
non-nil target. The nil tests pass without an explicit guard, as long as you do not
dereference before checking `ok`.

## Hint 6: Order Matters in Classify

Check validation, then rate limit, then the sentinel. The tests use errors that only
match one thing each, but the required order is part of the spec and costs nothing to
follow.
