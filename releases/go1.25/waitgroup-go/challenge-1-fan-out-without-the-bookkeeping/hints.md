# Hints for Fan Out Without the Bookkeeping

## Hint 1: The Shape of wg.Go

No `Add`, no `Done`, no `defer`:

```go
var wg sync.WaitGroup
for _, x := range items {
    wg.Go(func() {
        doSomething(x)
    })
}
wg.Wait()
```

`Go` increments the counter before starting the goroutine and decrements it when your
function returns.

## Hint 2: Allocate the Result Slice Up Front

```go
out := make([]string, len(ids))
```

Full length, not `make([]string, 0, len(ids))`. Every goroutine writes to its own index,
so the slice never grows and never needs a lock.

## Hint 3: Range Gives You Both

```go
for i, id := range ids {
    wg.Go(func() {
        out[i] = fetch(id)
    })
}
```

Since Go 1.22 each iteration has its own `i` and `id`, so capturing them in the closure
is correct as written. No `i := i` needed.

## Hint 4: Counting Needs Synchronisation

The slice case is safe because each goroutine owns a different element. A counter is
the opposite: everyone writes the same thing.

```go
var n atomic.Int64
...
n.Add(1)
...
return int(n.Load())
```

## Hint 5: If the Concurrency Test Times Out

`TestFetchAllActuallyRunsConcurrently` blocks every `fetch` until all four have arrived.
If you call `fetch` directly in the loop instead of inside `wg.Go`, the first call
blocks forever and the test hits its five second timeout. The `fetch` call belongs
inside the goroutine.

## Hint 6: If the Count Is Short

Either the counter is a plain `int` (increments are being lost) or the function returns
before `wg.Wait()`. Both show up as a number lower than expected.
