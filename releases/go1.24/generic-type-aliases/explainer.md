## One character, two very different types

Go has always had two ways to give a type a name.

```go
type Celsius float64     // a defined type: new, distinct from float64
type Temp    = float64   // an alias: another spelling of float64
```

The `=` is the whole difference. A defined type is a **new type** that happens to have
the same structure; an alias is **the same type** under a second name.

Until Go 1.24, only the first form could take type parameters:

```go
type Set[T comparable]   map[T]struct{}   // fine since 1.18
type Set[T comparable] = map[T]struct{}   // syntax error before 1.24
```

That gap had a real consequence: **naming a generic type always created a new type**,
whether you wanted a new type or not.

## Why that hurt

Say you just want a shorter name for a map used as a set.

```go
type Set[T comparable] map[T]struct{}   // pre-1.24, the only option
```

You have now created a type that is *not* `map[T]struct{}`. So this does not compile:

```go
func countKeys(m map[string]struct{}) int { return len(m) }

s := Set[string]{"go": {}}
countKeys(s)                 // cannot use s as map[string]struct{}
countKeys(map[string]struct{}(s))   // works, and now every call site has a conversion
```

Every boundary between your code and anything that speaks plain maps needs a
conversion. For a type that was meant to be a convenience, that is a lot of friction.

Go 1.24 gives you the other option:

```go
type Set[T comparable] = map[T]struct{}

s := Set[string]{"go": {}}
countKeys(s)   // fine: it is the same type
```

## Which one to reach for

The question to ask is: **do I want a new type, or a shorter name?**

Pick the **defined type** when the type is its own abstraction:

- you want methods on it (`func (s Set[T]) Add(v T)`)
- you want the compiler to stop people passing a raw map where a `Set` is meant
- you might change the representation later without changing callers

Pick the **alias** when you only want readability:

- long generic signatures that appear in many places
- a name for a type you do not own, so you could not add methods anyway
- a type that is deliberately interchangeable with its underlying form
- migrating a type between packages while keeping the old name working

That last one is the classic alias use, now available for generic types too:

```go
package oldpkg

// Deprecated: moved to newpkg.
type Cache[K comparable, V any] = newpkg.Cache[K, V]
```

Existing callers keep compiling, because it is not a different type. With a defined
type they would all break.

## The rule that surprises people

**You cannot define methods on an alias.**

```go
type Set[T comparable] = map[T]struct{}

func (s Set[T]) Add(v T) { s[v] = struct{}{} }   // invalid
```

This is not a limitation of generic aliases specifically. Methods belong to a defined
type, and an alias does not define one. `Set` here is just another way of writing
`map[T]struct{}`, and you have never been allowed to attach methods to `map[T]struct{}`
from your package.

So if your API wants methods, you want a defined type. That decision usually settles the
question before any of the other trade-offs matter.

Related, and worth knowing since it was clarified during the same work: you also cannot
declare methods on an *instantiated* generic type. `func (s Set[string]) Add(...)` is
not a thing, alias or not.

## Where aliases genuinely shine

**Taming a long signature.** When a generic type has several parameters, an alias at
the top of a file can make everything below it readable:

```go
type Handler[Req, Resp any] = func(context.Context, Req) (Resp, error)

var create Handler[CreateUserRequest, User]
```

**A constrained shorthand.** Fixing some parameters while leaving others open:

```go
type StringMap[V any] = map[string]V

var counts StringMap[int]   // map[string]int
```

**Package moves.** As above: the alias keeps the old import path compiling while callers
migrate on their own schedule.

> **Try it:** the challenge below builds a set on top of a generic alias, and the tests
> check the property that makes an alias an alias: a `Set[string]` can be handed
> straight to a function expecting a `map[string]struct{}`, with no conversion.
