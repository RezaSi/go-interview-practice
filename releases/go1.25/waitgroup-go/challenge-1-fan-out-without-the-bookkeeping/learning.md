# Learning Materials for Fan Out Without the Bookkeeping

## The Method

```go
func (wg *WaitGroup) Go(f func())
```

It does three things in the only order that is always correct: increment the counter on
the calling goroutine, start a goroutine, and arrange for the decrement to happen when
`f` returns even if it panics.

Compared with what it replaces:

```go
// before
wg.Add(1)
go func() {
    defer wg.Done()
    work()
}()

// after
wg.Go(func() {
    work()
})
```

## Two Kinds of Shared State

This challenge deliberately contains one case that needs synchronisation and one that
does not. Telling them apart is the skill worth taking away.

### Writing to distinct indices is safe

```go
out := make([]string, len(ids))
for i, id := range ids {
    wg.Go(func() {
        out[i] = fetch(id)   // no two goroutines touch the same element
    })
}
```

The slice header is read-only here: nobody appends, nobody reslices. Each goroutine
writes one element that no other goroutine touches. The race detector agrees, and no
mutex is required.

This only holds because the slice was allocated at full length first. If you had used
`append`, every goroutine would be writing the same slice header and you would need a
lock, and you would lose input order as a bonus.

### Incrementing a counter is not

```go
count++   // read, add, write: three steps, no atomicity
```

Two goroutines can read the same value and both write back the same increment, so an
increment disappears. `atomic.Int64` makes the read-modify-write a single indivisible
operation:

```go
var n atomic.Int64
n.Add(1)
...
n.Load()
```

A `sync.Mutex` around a plain `int` works equally well. For a single counter, the
atomic is smaller and faster.

## Preserving Order

Concurrency destroys ordering unless you rebuild it. Two ways:

**Index into a preallocated slice** (used here). Order is preserved by construction,
there is no shared write and no sorting afterwards.

**Collect and sort.** Send results to a channel with their index attached, then sort at
the end. More code, more allocation, and only worth it if you do not know the length in
advance.

## The Race Detector

The counter bug is invisible in normal test runs most of the time, which is why:

```bash
go test -race ./...
```

is worth running here. It instruments memory access and reports unsynchronised
read/write pairs even when the run happens to produce the right answer. `-race` costs
roughly 5 to 10 times the runtime, so it is a CI job rather than your inner loop.

## What wg.Go Does Not Do

- **No errors.** It takes `func()`, not `func() error`. Use
  `golang.org/x/sync/errgroup` when you need first-error semantics and cancellation.
- **No limits.** A thousand items start a thousand goroutines. Bound it with a
  semaphore channel or a worker pool if that matters.
- **No result plumbing.** Collecting answers is still up to you, which is what this
  challenge is about.
- **No fix for early reuse.** A `WaitGroup` still must not be reused until every
  previous `Wait` has returned.

## Further Reading

- [Go 1.25 release notes: sync](https://go.dev/doc/go1.25#sync)
- [Proposal #63796](https://github.com/golang/go/issues/63796)
- [Go blog: introducing the race detector](https://go.dev/blog/race-detector)
