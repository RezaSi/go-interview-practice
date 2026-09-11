## The four lines everybody copies

```go
var wg sync.WaitGroup

for _, job := range jobs {
	wg.Add(1)
	go func() {
		defer wg.Done()
		work(job)
	}()
}

wg.Wait()
```

You have written this many times. It is fine. It is also the shape of three separate
bugs, and Go 1.25 makes two of them impossible to express.

## What was added

One method:

```go
func (wg *WaitGroup) Go(f func())
```

The same loop:

```go
var wg sync.WaitGroup

for _, job := range jobs {
	wg.Go(func() {
		work(job)
	})
}

wg.Wait()
```

Two lines shorter per goroutine, and the counter is no longer your responsibility.
`Go` increments before starting the goroutine and decrements with a deferred call
inside it.

## The bugs it removes

### Forgetting `Done`

```go
wg.Add(1)
go func() {
	if err := work(job); err != nil {
		return          // wg.Done() never runs
	}
	wg.Done()
}()
```

An early return that skips the decrement, and `Wait` blocks forever. `defer wg.Done()`
is the fix everyone knows, which is exactly why this bug survives: it only appears when
somebody edits the goroutine later and adds a return path above the `Done`.

With `Go`, the deferred decrement is written by the standard library, so there is no
edit that can drop it.

### `Add` inside the goroutine

```go
for _, job := range jobs {
	go func() {
		wg.Add(1)          // too late
		defer wg.Done()
		work(job)
	}()
}
wg.Wait()                  // may return before anything ran
```

This one is nastier because it usually works. If `Wait` runs before any goroutine has
been scheduled, the counter is still zero, so `Wait` returns immediately and your
program carries on as though the work finished. It fails under load, on a busy CI
machine, or on the day you add a slower first job.

`wg.Go` increments on the calling goroutine, before the new one exists. The ordering
cannot be got wrong because you are not the one writing it.

### Reusing a `WaitGroup` too early

```go
wg.Wait()
wg.Add(1)   // racy if any goroutine from the previous round is still finishing
```

**`Go` does not fix this one.** `WaitGroup` reuse still requires that all previous
`Wait` calls have returned before new `Add` calls happen. Worth knowing, because it is
easy to assume the new method solved every `WaitGroup` problem.

## What it is not

`wg.Go` is not an errgroup. It takes a `func()`, not a `func() error`, so there is no
error collection, no first-error cancellation, and no `context`. If you need those,
`golang.org/x/sync/errgroup` is still the answer and its `Go` method is where this one
got its name and shape.

It also does not limit concurrency. A thousand jobs start a thousand goroutines. If you
need a bound, that is still a semaphore or a worker pool.

What you get is exactly one thing: the counting is handled. That is a smaller feature
than an errgroup and it is available without a dependency, which is the trade being
made.

## The loop variable, while we are here

Note that both examples above capture `job` directly inside the closure. Before Go 1.22
that was the single most common Go bug, and the fix was `job := job` or passing it as
an argument. Since 1.22 each iteration has its own variable, so the capture is correct
as written.

If you are reading older code that does `wg.Add(1)` and `go func(j Job){...}(job)`, the
parameter was working around a language wart that no longer exists. When you convert
that to `wg.Go`, you can usually drop the parameter too.

> **Try it:** the challenge below is a small fan-out that collects results from several
> goroutines. One of the tests uses a barrier that only completes if the goroutines
> really do run at the same time, so a sequential implementation deadlocks rather than
> quietly passing.
