## Four years from idea to release

Proposal [#51945](https://github.com/golang/go/issues/51945) was filed in March 2022,
about six weeks after generics shipped in Go 1.18. The idea is not subtle: `errors.As`
has a pointer-shaped out-parameter because it predates type parameters, and now it does
not have to.

It then sat for nearly four years. That delay is the interesting part, because nobody
disagreed about the ergonomics.

## The objections that had to be answered

**"Two functions doing the same thing."** The strongest argument against. Every Go
programmer would now have to know both, and reviewers would have to decide which is
"correct" in each package. The answer that eventually landed is that they are not
equivalent: `errors.As` can target a type known only at run time, `AsType` cannot,
because a type argument is a compile-time thing. So there is a real, if narrow, case
that only `As` serves.

**The name.** `AsType` is not the first choice anyone reaches for. `As2`, `AsT`,
`Get`, `Find`, `To` and a generic overload of `As` itself were all discussed. A generic
`As` was impossible: Go has no overloading, and changing the existing signature would
break every caller. `AsType` won on being unambiguous at a call site rather than on
elegance, which is generally how standard library names get chosen.

**"Wait for a general pattern."** For a while the plan was to see how generic accessors
looked elsewhere in the standard library before committing to a shape here. That
patience is visible across Go 1.21 to 1.25, where `slices`, `maps` and `cmp` worked out
the conventions for generic helpers. `AsType` arrives after that groundwork, not before
it.

## The change itself

Commit [`a846bb0`](https://github.com/golang/go/commit/a846bb0aa523), reviewed as
[CL 707235](https://go-review.googlesource.com/c/go/+/707235):

```text
errors: add AsType

Fixes #51945

GitHub-Pull-Request: golang/go#75621
```

Two things stand out.

**It came in through GitHub.** The `GitHub-Pull-Request: golang/go#75621` trailer means
this arrived as a pull request from a contributor and was mirrored into Gerrit, rather
than being written by the core team. A four year old proposal, implemented by someone
outside the team once it was accepted, is the Go contribution process working exactly as
documented.

**The commit message is two lines.** No design discussion, because that all happened on
the issue. By the time a proposal is accepted, the code is the least contentious part.

The implementation is a thin wrapper over the existing machinery: allocate a `T`, call
the same reflection-based chain walk that `errors.As` uses, and return the value with a
boolean rather than reporting through the pointer. There is no new matching logic, which
is exactly why behaviour is identical, including custom `As(any) bool` methods.

## What generics actually bought here

It is worth being precise, because "generics make this nicer" is doing real work in this
story.

`errors.As` has the signature it has because in 2019 there was no way to write "give me
back a value of the type the caller names". The out-parameter is a workaround for a
missing language feature, and it drags three consequences with it: a declaration line, a
widened scope, and a run-time type check where a compile-time one belongs.

```go
func As(err error, target any) bool           // 2019: target is the only option
func AsType[T error](err error) (T, bool)     // 2026: T says it in the signature
```

The lesson generalises. If your API takes an `any` and immediately reflects on it to
decide what the caller meant, that is often a type parameter waiting to be written. The
version with the type parameter moves the mistake from run time to compile time, and it
usually reads better too.

## Read more

- [Proposal #51945: errors: add AsType](https://github.com/golang/go/issues/51945)
- [CL 707235](https://go-review.googlesource.com/c/go/+/707235) — the implementation
- [Go blog: working with errors in Go 1.13](https://go.dev/blog/go1.13-errors) — where
  `Is`, `As` and `%w` came from
