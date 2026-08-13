# Learning Materials for A Cache That Expires

## Two Ways to Expire an Entry

There are only two designs, and this challenge uses both.

**Lazy expiry** happens in `Get`. The entry stays in the map after it dies, and the
reader notices. Cheap, and it means a stale entry is never served even if no janitor is
running.

**Active expiry** is the janitor. It reclaims memory, which lazy expiry alone never
does: a key that is written once and never read again sits there forever.

Real caches do both, which is why `Get` checks the timestamp *and* `Cleanup` sweeps.

## Writing the Expiry Check

```go
e, ok := c.items[key]
if !ok || !time.Now().Before(e.expires) {
    return "", false
}
return e.value, true
```

Store an absolute deadline at write time rather than a duration plus a start time.
Comparing two `time.Time` values is one operation and there is no arithmetic to get
wrong at read time.

Prefer `Before` and `After` to subtracting and comparing durations. `time.Since(x) >
ttl` works but reads worse and invites off-by-one thinking around the boundary.

## The Ticker Loop

```go
tk := time.NewTicker(interval)
defer tk.Stop()
for {
    select {
    case <-ctx.Done():
        return
    case <-tk.C:
        c.sweep()
    }
}
```

Three details worth keeping:

- **`defer tk.Stop()`.** An unstopped ticker keeps a runtime timer alive. In a long
  running process that is a leak.
- **`ctx.Done()` first in the select.** Not required, `select` picks randomly among
  ready cases, but it documents the intent.
- **A ticker, not `time.Sleep` in a loop.** A ticker's period includes the time your
  sweep takes; a sleep loop drifts by the duration of the work.

## Deleting From a Map While Ranging It

The spec allows it:

> The iteration order over maps is not specified... If a map entry that has not yet
> been reached is removed during iteration, the corresponding iteration value will not
> be produced.

So the obvious sweep is correct:

```go
now := time.Now()
for k, e := range c.items {
    if !now.Before(e.expires) {
        delete(c.items, k)
    }
}
```

Read `now` once. Calling `time.Now()` inside the loop judges each entry against a
slightly different instant, which is harmless here and sloppy in general.

## What synctest Changes About Testing This

The cache code is ordinary. The tests are where the new package shows up.

### Bubbles

```go
synctest.Test(t, func(t *testing.T) { ... })
```

Inside, the `time` package is fake. `time.Sleep(time.Hour)` returns immediately, but
`time.Now()` really has moved an hour forward, so any code reading the clock behaves
exactly as it would after a real hour. `context.WithTimeout` uses the same fake clock,
because the interception is below `time` rather than in your own code.

### The advancement rule

The clock moves only when **every** goroutine in the bubble is durably blocked. That
is what makes the test deterministic: either something is still running, in which case
time is frozen, or nothing can run, in which case time jumps to the next timer. Your
assertion never races the scheduler.

### `Wait` is not "advance"

```go
go c.Cleanup(ctx, 10*time.Second)
synctest.Wait()
```

`Wait` returns once every *other* goroutine in the bubble is durably blocked. It says
"the background work has settled", not "time has passed". You reach for it when you
want to assert on state without knowing how long the work takes.

### Leaks become failures

If a goroutine started inside the bubble is still running when the body returns, the
test fails. That is why `Cleanup` has to respect `ctx.Done()`, and it is a genuinely
useful side effect: goroutine leaks normally show up as slow memory growth in
production rather than a red test on your laptop.

## Further Reading

- [Go 1.25 release notes: testing/synctest](https://go.dev/doc/go1.25#testingsynctest)
- [Proposal #67434](https://github.com/golang/go/issues/67434)
- [Go spec: for statements with range clause](https://go.dev/ref/spec#For_range)
