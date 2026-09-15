# Hints for Optional Fields Without Temporaries

## Hint 1: new Takes the Value Directly

The parameter is already a value, so pass it:

```go
Retries: new(retries),
```

That allocates an `int` holding whatever `retries` is and gives you the `*int`.

## Hint 2: The Whole Constructor Is One Literal

```go
return Config{
    Retries: new(retries),
    Debug:   new(debug),
    Label:   new(label),
}
```

`Timeout` is not mentioned, so it stays `nil`, which is what the test wants.

## Hint 3: Bump Is One Line

```go
return new(n + 1)
```

The expression is evaluated first, then stored in a fresh variable, then you get its
address.

## Hint 4: Why the Aliasing Test Passes for Free

`n` is a parameter, so it is already a copy of the caller's variable. `new(n + 1)`
copies again into new memory. There is no path from the returned pointer back to
anything the caller holds, so writing through it cannot be observed outside.

## Hint 5: If the Zero Value Test Fails

You probably guarded the assignment, something like `if retries != 0 { ... }`. Do not.
Zero is a legitimate value here; only the *absence* of a value is `nil`, and this
constructor always receives one.

## Hint 6: If Two Configs Share a Pointer

That happens if you hoisted a single variable out and reused its address across calls.
Every `new(...)` call allocates separately, so keeping the calls inside the literal is
enough.
