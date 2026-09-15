# Learning Materials for Naming a Generic Type

## Defined Types and Aliases

Go has two ways to name a type, and the `=` is the entire difference:

```go
type Celsius float64     // defined type: new and distinct
type Temp    = float64   // alias: another name for float64
```

A **defined type** is a new type. It has its own identity, its own method set, and needs
an explicit conversion to move between it and its underlying type.

An **alias** is a second spelling. Everything that was true of the original stays true,
because there is only one type involved.

Go 1.24 allowed type parameters on the alias form for the first time:

```go
type Set[T comparable] = map[T]struct{}
```

## The Property That Matters

```go
func CountKeys(m map[string]struct{}) int

s := NewSet("a", "b")
CountKeys(s)                       // compiles: it is the same type
var raw map[string]struct{} = s    // assignable, both directions
```

With a defined type, both of those need a conversion, and every boundary between your
code and anything speaking plain maps grows a `map[string]struct{}(s)`.

You can confirm identity at run time:

```go
reflect.TypeOf(Set[string]{}) == reflect.TypeOf(map[string]struct{}{})   // true
```

## No Methods on an Alias

```go
type Set[T comparable] = map[T]struct{}

func (s Set[T]) Add(v T) { }   // invalid
```

Methods attach to defined types. An alias does not define a type, so there is nothing to
attach to. This is not specific to generics: you have never been able to write methods on
`map[string]struct{}` from your own package either.

There is a related rule clarified during the same work: you cannot declare methods on an
*instantiated* generic type, so `func (s Set[string]) Add(...)` is not valid regardless.

## Choosing Between Them

Reach for a **defined type** when the type is its own abstraction:

- you want methods
- you want the compiler to reject a raw map where a `Set` is meant
- you may change the representation later without touching callers

Reach for an **alias** when you only want a name:

- long generic signatures repeated throughout a file
- a type you do not own, so methods were never an option
- something deliberately interchangeable with its underlying form
- moving a type between packages without breaking callers

That last one is the classic:

```go
package oldpkg

// Deprecated: moved to newpkg.
type Cache[K comparable, V any] = newpkg.Cache[K, V]
```

Existing code keeps compiling because it is not a different type. With a defined type,
every caller breaks.

## The Empty Struct Set

```go
map[T]struct{}
```

`struct{}` occupies zero bytes, so the map stores only keys. `map[T]bool` works too and
costs a byte per entry plus the question of what `false` means. The `struct{}` form is
the convention.

## Maps Are Reference Types

```go
func Union[T comparable](a, b Set[T]) Set[T]
```

`a` and `b` are copies of the map headers, not of the data. Writing `a[k] = ...` inside
the function is visible to the caller. That is why `Union` has to allocate a new map, and
why one of the tests checks for it.

## Further Reading

- [Go 1.24 release notes: language changes](https://go.dev/doc/go1.24#language)
- [Proposal #46477: generic type aliases](https://github.com/golang/go/issues/46477)
- [Go spec: alias declarations](https://go.dev/ref/spec#Alias_declarations)
