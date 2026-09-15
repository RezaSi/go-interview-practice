## A proposal that shipped twice

`testing/synctest` is unusual: it landed as `GOEXPERIMENT=synctest` in Go 1.24, spent
a full release cycle in the wild, changed its main entry point based on what people
hit, and only then graduated in 1.25.

The API changed in that year, and the change tells you what the hard part was.

The experimental version had:

```go
func Run(f func())   // Go 1.24, experimental
```

The released version has:

```go
func Test(t *testing.T, f func(*testing.T))   // Go 1.25
```

Two differences, both learned from use. The bubble now receives a `*testing.T`, so
failures inside it are attributed to the right test and `t.Fatal` behaves the way you
expect. And the name moved from `Run` to `Test`, making it obvious at the call site
that this is a testing construct rather than a general concurrency primitive, which is
a thing people were starting to reach for.

## The idea is older than the package

Fake clocks are not new. Every large Go codebase has one:

```go
type Clock interface {
	Now() time.Time
	After(d time.Duration) <-chan time.Time
}
```

You thread it through every constructor, use the real one in production, a fake one in
tests. It works, and it has two costs that never go away.

**It is viral.** Every type that touches time needs the clock parameter, and so does
everything that constructs it, all the way up. A `Cache` that wants a timeout drags a
`Clock` through six layers of your API.

**It does not cover the standard library.** `context.WithTimeout` calls the real
`time` package. So does `http.Client.Timeout`, and any dependency you did not write.
Your fake clock stops at the edge of your own code.

The proposal's insight is that the interception belongs one level down, in the runtime,
where it can cover `time` itself and therefore everything built on it. You get a fake
clock for `context`, for third-party libraries, for code you have never read, without
changing a single signature.

## Why "durably blocked" is the whole design

The naive fake clock has an obvious failure mode: when do you advance it?

If the test advances it manually (`clock.Advance(time.Second)`), you are back to
guessing, just with better syntax. Advance too early and the goroutine has not reached
its `Sleep` yet, so it sleeps a full second past the point you wanted. That is the
original flake wearing a new hat.

The alternative the proposal settled on is to let the *scheduler* decide. Advance the
clock only when no goroutine can make progress. That single rule gives you:

- **Determinism.** The set of runnable goroutines fully determines what happens next.
  There is no timing to race against.
- **No API for advancing.** You cannot get it wrong because you cannot do it at all.
- **A leak detector for free.** If a goroutine is blocked on something outside the
  bubble, the bubble can never become idle, so the test hangs rather than silently
  passing. Unpleasant the first time it happens, correct every time.

Getting there needs the runtime's cooperation. The scheduler has to distinguish
"blocked on something inside this bubble" from "blocked on a syscall", and it has to
know which goroutines belong to which bubble. That is not something a library can do
from the outside, and it is the reason this took a runtime change rather than a clever
package.

## Reading the constraints as design advice

The bubble's restrictions are usually presented as caveats. They read better as a
description of testable code.

**"Real I/O is not faked."** A unit that does its own networking cannot be tested
deterministically by anyone, with or without this package. Take the dependency as a
parameter.

**"Channels shared with the outside break it."** If your goroutine's liveness depends
on something the test cannot see, the test cannot make assertions about it. Same
conclusion.

**"The bubble must drain."** If your function leaks goroutines, you want to know. This
turns a class of bug that normally surfaces as slow memory growth in production into a
failing test on your laptop.

Every one of those pushes toward the same shape: dependencies in, no hidden global
state, no orphaned goroutines. The package is opinionated, and the opinions are ones
most Go style guides already hold.

## What to read next

- [Proposal #67434: testing/synctest](https://github.com/golang/go/issues/67434) —
  the full discussion, including the `Run` to `Test` change
- [Go 1.25 release notes](https://go.dev/doc/go1.25#testingsynctest)
- [`src/testing/synctest`](https://github.com/golang/go/tree/master/src/testing/synctest)
  — the package is small; most of the machinery lives in the runtime

## The change itself

The first commit is [`d90ce58`](https://github.com/golang/go/commit/d90ce588eac7),
reviewed as [CL 591997](https://go-review.googlesource.com/c/go/+/591997), and its
message is a tidy summary of the whole design:

```text
internal/synctest: new package for testing concurrent code

Add an internal (for now) implementation of testing/synctest.

The synctest.Run function executes a tree of goroutines in an
isolated environment using a fake clock. The synctest.Wait function
allows a test to wait for all other goroutines within the test
to reach a blocking point.

For #67434
For #69687
```

Two details in that message are the whole story.

**"internal (for now)"** is the release strategy stated up front. The package landed as
`internal/synctest`, was exposed as `testing/synctest` behind `GOEXPERIMENT=synctest`
in Go 1.24, and only became a normal import in 1.25 once the API had survived contact
with real test suites. `Run` became `Test(t, f)` in that window.

**"a tree of goroutines"** is why this needed the runtime. Bubble membership is
inherited: a goroutine started inside the bubble joins it, and so does anything it
starts. Only the scheduler is in a position to track that.

The follow-up commits show where the difficulty actually was, and almost none of it is
in the package itself:

- `runtime: record synctest bubble ownership in hchan` ([CL 674515](https://go-review.googlesource.com/c/go/+/674515)) —
  channels have to know which bubble they belong to, otherwise "durably blocked" cannot
  be decided
- `runtime: avoid panic in expired synctest timer chan read` ([CL 634955](https://go-review.googlesource.com/c/go/+/634955)) —
  fake timers and channel reads interacting badly
- `runtime: print blocking status of bubbled goroutines in stacks` ([CL 670976](https://go-review.googlesource.com/c/go/+/670976))
  and `runtime: clarify stack traces for bubbled goroutines` ([CL 679376](https://go-review.googlesource.com/c/go/+/679376)) —
  when a bubble deadlocks you need a stack trace that says *why*, so the panic output
  had to learn about bubbles

That last pair is the honest bit. A fake clock that hangs is useless if the failure is
unreadable, so a meaningful slice of the work went into the error message you see when
your bubble never goes idle.

The adoption commits are worth a look too: `database/sql: use synctest in tests` and
`net/http/internal/http2: call synctest.Test directly`. The standard library moved its
own timing tests over, which is the strongest signal available that the API is the one
they intend you to use.
