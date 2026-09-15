# Challenge 1: Naming a Generic Type

## The situation

Two aliases are declared for you:

```go
type Set[T comparable] = map[T]struct{}

type Pair[A, B any] = struct {
	First  A
	Second B
}
```

Note the `=`. These are **aliases**, not defined types, which was a syntax error before
Go 1.24. `Set[string]` is not merely *like* `map[string]struct{}`, it **is**
`map[string]struct{}`.

## The task

Implement four generic functions.

```go
func NewSet[T comparable](items ...T) Set[T]
func Union[T comparable](a, b Set[T]) Set[T]
func Contains[T comparable](s Set[T], v T) bool
func Zip[A, B any](as []A, bs []B) []Pair[A, B]
```

- `NewSet` deduplicates, and with no items returns an empty, usable set rather than
  `nil`.
- `Union` returns a **new** set. Neither input may be modified.
- `Zip` stops at the shorter slice.

Leave `CountKeys` exactly as it is. It takes a plain `map[string]struct{}` on purpose.

## Why these are functions and not methods

You might expect `s.Add(v)` and `a.Union(b)`. You cannot have them:

```go
func (s Set[T]) Add(v T) { }   // invalid: Set is an alias
```

Methods belong to a **defined** type. An alias does not define one, so `Set` here is
just another way of spelling `map[T]struct{}`, and you have never been able to attach
methods to that from your own package.

If you want methods, you want `type Set[T comparable] map[T]struct{}` without the `=`,
and you accept that it becomes a distinct type needing conversions at every boundary.
That trade-off is the whole point of this feature existing.

## What the tests check

`TestAliasIsTheSameType` is the one that matters:

```go
CountKeys(s)                            // no conversion
var raw map[string]struct{} = s         // assignable both ways
reflect.TypeOf(s) == reflect.TypeOf(map[string]struct{}{})
```

All three fail immediately if `Set` were a defined type. `TestSetSupportsBuiltins` makes
the same point from the other side: `delete`, indexing and `len` all work because it is
a map.

## Run it

```bash
go test -v ./...
```
