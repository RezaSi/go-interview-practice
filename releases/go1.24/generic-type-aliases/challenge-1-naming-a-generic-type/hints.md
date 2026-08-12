# Hints for Naming a Generic Type

## Hint 1: A Set Is Just a Map

Because `Set[T]` is an alias, every map operation works directly:

```go
s := make(Set[T], len(items))
s[item] = struct{}{}
_, ok := s[item]
delete(s, item)
```

`struct{}{}` is the empty value: it takes no memory, which is why this is the idiomatic
set in Go.

## Hint 2: NewSet Is make Plus a Loop

```go
s := make(Set[T], len(items))
for _, it := range items {
    s[it] = struct{}{}
}
return s
```

Using `make` rather than declaring a nil map is what makes the no-items case return an
empty usable set instead of `nil`.

## Hint 3: Contains Is the Comma-ok Form

```go
_, ok := s[v]
return ok
```

## Hint 4: Union Must Allocate

Copy both inputs into a fresh set:

```go
u := make(Set[T], len(a)+len(b))
for k := range a { u[k] = struct{}{} }
for k := range b { u[k] = struct{}{} }
return u
```

Writing into `a` and returning it would pass the union test and fail the one that checks
the inputs are untouched. Maps are reference types, so the caller sees the change.

## Hint 5: Zip Needs the Shorter Length First

```go
n := len(as)
if len(bs) < n {
    n = len(bs)
}
```

Then build `Pair[A, B]{First: as[i], Second: bs[i]}` for each index below `n`.

## Hint 6: See the Alias Restriction Yourself

Add this and run the tests:

```go
func (s Set[T]) Add(v T) { s[v] = struct{}{} }
```

The compiler refuses: methods need a defined type, and an alias is not one. That single
error is the clearest statement of the difference between the two forms.
