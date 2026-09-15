# Challenge 1: Optional Fields Without Temporaries

## The situation

`Config` uses pointers for its fields so that **unset** and **set to the zero value**
are different things. `Timeout` being `nil` means nobody chose a timeout;
`*Timeout == 0` would mean somebody chose zero.

Populating a struct like that used to need a named variable per field:

```go
retries := 3
debug := true
label := "primary"

cfg := Config{Retries: &retries, Debug: &debug, Label: &label}
```

Go 1.26 lets `new` take a value, so those three variables can go away.

## The task

### `NewConfig`

```go
func NewConfig(retries int, debug bool, label string) Config
```

Set `Retries`, `Debug` and `Label` from the arguments using `new(expr)`. Leave
`Timeout` as `nil`, because this constructor does not take one.

### `Bump`

```go
func Bump(n int) *int
```

Return a pointer to `n + 1`. One expression.

Leave `Describe` alone.

## What the tests are checking

Beyond the obvious, three of them are about semantics rather than values:

**`TestNewConfigCarriesZeroValues`** calls `NewConfig(0, false, "")`. Every pointer must
still be non-nil. Zero is a value somebody chose, and the whole reason these fields are
pointers is to keep that distinguishable from unset.

**`TestBumpDoesNotAliasTheArgument`** writes through the returned pointer and checks the
caller's variable did not move. `new(expr)` allocates a **new variable** holding the
value, so it cannot alias anything.

**`TestEachCallGetsItsOwnMemory`** builds two configs and writes through one. Each call
allocates independently.

## Rules

- Use `new(expr)`. No temporary variables, no `ptr` helper.
- Leave `Timeout` nil.
- Do not change `Describe`.

## Run it

```bash
go test -v ./...
```
