## The feature that was always supposed to exist

Proposal [#46477](https://github.com/golang/go/issues/46477) is not really a feature
request. It is a gap report.

When generics landed in Go 1.18, defined types got type parameters and aliases did not.
That was not a design decision anyone argued for; it was scope control. Generics were
already the largest change in the language's history, and aliases were the piece that
could be left until later without blocking anything.

"Later" turned out to be three releases, and the reason is instructive.

## Why an alias is harder than it looks

A defined type is a fresh entity. `type Set[T any] map[T]struct{}` creates a new named
type, and the type checker can treat `Set[int]` as its own thing with its own identity.

An alias has to be **transparent**. `Set[int]` must be indistinguishable from
`map[int]struct{}` everywhere: in assignability, in type identity, in method sets, in
error messages. Adding type parameters means the substitution has to happen *before*
any of those questions are asked, and it has to happen consistently across two separate
type checkers, `go/types` and `cmd/compile/internal/types2`, plus the compiler's own
export data format.

The hard cases are the ones where an alias meets another generic construct:

- a generic alias used as a type argument to another generic type
- a generic alias whose right-hand side is itself an instantiated generic type
- an alias instantiated with the wrong number of type arguments
- a defined type whose underlying type is a generic alias

Each of those is a place where "just substitute the parameters" can go wrong, and the
commit history is largely a record of finding them.

## The change itself, in three acts

This one is a good example of how Go ships a risky language change: behind a flag, then
on by default, then the flag is deleted.

**Act one, the flag.**
[CL 586955](https://go-review.googlesource.com/c/go/+/586955),
`internal/goexperiment: add aliastypeparams GOEXPERIMENT flag`. The feature goes in
switched off, so it can be developed on master without affecting anyone.

**Act two, the type checkers.**
[CL 567617](https://go-review.googlesource.com/c/go/+/567617),
`go/types, types2: instantiate generic alias types`, is the core of it. Both checkers
learn to substitute type arguments into an alias's right-hand side.

Then the edge cases arrive as their own fixes, and they are worth listing because they
are exactly the cases above:

- [CL 601115](https://go-review.googlesource.com/c/go/+/601115) —
  `types2, go/types: fix instantiation of named type with generic alias`
- [CL 641857](https://go-review.googlesource.com/c/go/+/641857) —
  `go/types, types2: don't panic when instantiating generic alias with wrong number of type arguments`
- [CL 629080](https://go-review.googlesource.com/c/go/+/629080) —
  `go/types, types2: disallow new methods on generic alias and instantiated types`

That third one is the rule you actually feel as a user: no methods on an alias. It is
not an arbitrary restriction; it falls out of an alias not defining a type. The CL is
the type checker being taught to say so clearly instead of accepting something
meaningless.

There is also [CL 608595](https://go-review.googlesource.com/c/go/+/608595),
`cmd/compile/internal/noder: write V2 bitstream aliastypeparams=1`, which is the export
data format learning to represent a generic alias so it can cross a package boundary at
all.

**Act three, the flag goes away.**
[CL 691956](https://go-review.googlesource.com/c/go/+/691956),
`all: delete aliastypeparams GOEXPERIMENT`. Once the feature shipped on by default in
1.24 and survived a release, the escape hatch was removed. A GOEXPERIMENT that never
gets deleted is a maintenance burden and a signal of doubt; deleting it is the project
saying the question is settled.

## What to take from it

**Transparency is expensive to implement and cheap to use.** The reason this took three
releases is the same reason it is pleasant: an alias has to be invisible, and invisible
things have to be correct in every context rather than just the common one.

**Flags are for features, not just for experiments.** The `aliastypeparams` arc is a
template worth copying in your own projects: build behind a switch, default it on when
you believe it, delete the switch once you are sure. All three steps, including the
last one.

**The restriction is a consequence, not a rule.** "No methods on an alias" was not
chosen. It follows from what an alias is, and the CL that enforces it exists to make the
compiler explain that rather than to impose it. When you meet a limitation in a language
you use, it is worth checking which of the two kinds it is.

## Read more

- [Proposal #46477: generic type aliases](https://github.com/golang/go/issues/46477)
- [Go 1.24 release notes](https://go.dev/doc/go1.24#language)
- [Go spec: alias declarations](https://go.dev/ref/spec#Alias_declarations)
