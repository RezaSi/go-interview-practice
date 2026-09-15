## Why a language change for two lines

Changing the language is the most expensive thing the Go team can do. Every addition
lands in every tutorial, every code review and every linter forever. So the bar is not
"is this nicer", it is "does this remove a recurring problem that cannot be solved in a
library".

The argument that carried this one is that the workaround **is already in every
codebase**, just written slightly differently in each:

```go
func ptr[T any](v T) *T { return &v }
func Ptr[T any](v T) *T { return &v }
func addr[T any](v T) *T { return &v }
func toPtr[T any](v T) *T { return &v }
```

When the same six-character function is independently reinvented in thousands of
projects, that is evidence of a missing primitive rather than a style preference. It
also only became expressible as a helper at all in Go 1.18; before generics you needed
one per type.

The counter-argument, made repeatedly, was that a helper is *fine*. The response is
that a helper is fine for your package and useless across a boundary: you cannot use
another module's `ptr` without importing it, so everyone keeps writing their own.

## Extending new rather than adding syntax

Several spellings were considered before this one:

- **`&expr` for any expression.** The obvious minimal change, and the most dangerous.
  `&x` where `x` is a variable already means "alias this variable". Making `&f()` legal
  puts two very different meanings on the same operator, and readers would have to know
  whether the operand was addressable to know which one they were looking at.
- **A new built-in, say `ref(x)`.** Costs a new predeclared identifier, which shadows
  user code named `ref` everywhere. `new` is already reserved.
- **A library function in the standard library.** Would need to be generic, and
  `stdlib.Ptr(3)` reads worse than `new(3)` while solving nothing that a local helper
  does not.

Overloading `new` wins because `new` already means "allocate and give me a pointer".
The change is to what you may pass it, not to what it does. Existing code keeps
compiling because a type argument still behaves exactly as before.

## The change itself

This one landed as a coordinated set rather than a single commit, which is normal for a
language change. From the compiler commit,
[`7bc1935`](https://github.com/golang/go/commit/7bc1935db55c), reviewed as
[CL 705157](https://go-review.googlesource.com/c/go/+/705157):

```text
cmd/compile/internal: support new(expr)

This CL adds compiler support for new(expr),
a feature of go1.26 that allows the user to specify
the initial value of the variable instead of its
type.

Also, a basic test of dynamic behavior.

See CL 704737 for spec change and CL 704935 for
type-checker changes.
```

That message is a map of how a Go language change is built, and the order matters:

| CL | Layer | What it does |
|---|---|---|
| [704737](https://go-review.googlesource.com/c/go/+/704737) | spec | defines the meaning first |
| [704935](https://go-review.googlesource.com/c/go/+/704935) | type checker | `go/types` and `types2` learn the new form |
| [705157](https://go-review.googlesource.com/c/go/+/705157) | compiler | code generation, plus a behavioural test |

The spec goes first on purpose. If the wording cannot be written clearly, the feature is
not ready, and no amount of working code fixes that.

Two follow-ups are worth reading because they show where the edges were:

- [CL 713241](https://go-review.googlesource.com/c/go/+/713241),
  `go/types, types2: only report version errors if new(expr) is ok otherwise` — if you
  compile `new(f())` with an older language version, you should be told "this needs Go
  1.26", not a confusing type error. But only if the code would otherwise be valid,
  since a genuinely broken call deserves the real error.
- [CL 737680](https://go-review.googlesource.com/c/go/+/737680),
  `cmd/compile/internal/staticinit: fix bug in global new(expr)` — `new(expr)` at
  package level interacts with static initialisation, where the compiler tries to
  compute values at build time rather than at start-up. A real bug, found and fixed
  before release.

And [CL 722482](https://go-review.googlesource.com/c/go/+/722482), `spec: more precise
prose for built-in function new`, is the spec being tightened *after* implementation,
which is the usual sign that writing the code exposed a case the words did not quite
cover.

## What to take from it

The transferable idea is **extend an existing concept instead of adding a new one**.

`new` already had a job. Widening what it accepts cost no new keyword, no new
identifier, and no new mental model: if you knew what `new` did, you already know what
`new(expr)` does. Compare that with a hypothetical `ref` builtin, which would have
needed a name, a place in the docs, and a paragraph explaining how it differs from
`new` and `&`.

When you are designing an API and reaching for a new function, it is worth asking
whether an existing one could take a wider argument instead.

## Read more

- [Go 1.26 release notes](https://go.dev/doc/go1.26)
- [Go spec: allocation](https://go.dev/ref/spec#Allocation)
