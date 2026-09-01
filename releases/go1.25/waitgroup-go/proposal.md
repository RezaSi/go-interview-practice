## A two year argument about six lines

Proposal [#63796](https://github.com/golang/go/issues/63796) was opened in 2023 and the
API it landed with is almost exactly the one it opened with. What took the time was not
the design. It was deciding whether a change this small was worth making at all.

The case against is real, and worth stating properly:

- The pattern is already idiomatic, documented and understood.
- `errgroup.Group` exists and does more.
- Every method added to `sync` is a method every Go programmer eventually has to know.

The case for won on one observation: the manual version is not merely verbose, it has
**failure modes that the compiler cannot see**. `wg.Add(1)` in the wrong place produces
a program that passes tests and returns early under load. A method that makes the
ordering impossible to express incorrectly is a different kind of change from one that
merely saves keystrokes.

## The shape it did not take

Several alternatives were discussed and rejected, and each rejection is informative.

**`wg.Go(f func() error)`.** The obvious wish: match `errgroup`. Rejected because
`sync.WaitGroup` has nowhere to put an error. Storing one means deciding whether it is
the first or the last, whether other goroutines get cancelled, and what `Wait` returns.
Those are the decisions that make `errgroup` a different type, and folding them into
`WaitGroup` would have made the simple case worse to serve the complex one.

**A variadic `wg.Go(fs ...func())`.** Saves a loop occasionally, and makes the common
single-function call marginally more confusing to read. Dropped.

**`wg.Go` returning something.** A handle, a channel, anything. Every version made it
possible to leak or ignore the returned value, which is the sort of thing `WaitGroup`
exists to prevent.

The final signature is deliberately the least powerful one that removes the bugs:

```go
func (wg *WaitGroup) Go(f func())
```

## The change itself

Commit [`822031d`](https://github.com/golang/go/commit/822031dffc43), reviewed as
[CL 662635](https://go-review.googlesource.com/c/go/+/662635):

```text
sync: add WaitGroup.Go

Fixes #63796
```

That is the entire commit message, which is a fair reflection of the implementation.
There is no new state on `WaitGroup`, no new synchronisation, and no runtime change.
`Go` is a composition of the two methods that already existed, in the one order that is
always correct:

1. `Add(1)` on the calling goroutine
2. start the goroutine
3. `defer Done()` inside it
4. call `f`

Step 1 happening before step 2 is the whole point. The increment is ordered before the
goroutine exists, so there is no window in which `Wait` can observe a counter that has
not been raised yet.

A follow-up, [CL 662975](https://go-review.googlesource.com/c/go/+/662975)
(`sync: tidy WaitGroup documentation, add WaitGroup.Go example`), did the other half of
the work: making sure the package docs point newcomers at `Go` rather than at the
pattern it replaces. For a change whose entire value is "people will use the safe form
by default", the documentation edit is arguably the more important commit.

## The lesson worth taking

The interesting thing here is not the method. It is the type of change: **making an
ordering constraint unrepresentable rather than documenting it**.

The old pattern's correctness depended on the programmer putting one line above
another. No compiler check, no vet analyser, no type error. Just a convention, and a
failure that only shows up under scheduling pressure.

The new method does not check the ordering. It removes the ability to express it. That
is a design move you can apply in your own APIs, and it is usually available in exactly
these situations: when your documentation contains the phrase "you must call X before
Y".

There is a vet analyser for the old pattern too, added in the same release: the
`waitgroup` check reports `wg.Add` calls misplaced inside the goroutine. That is the
belt to `Go`'s braces, for the code that has not been converted yet.

## Read more

- [Proposal #63796: sync: add WaitGroup.Go](https://github.com/golang/go/issues/63796)
- [Go 1.25 release notes](https://go.dev/doc/go1.25#sync)
- [`golang.org/x/sync/errgroup`](https://pkg.go.dev/golang.org/x/sync/errgroup) — where
  to go when you need errors and cancellation
