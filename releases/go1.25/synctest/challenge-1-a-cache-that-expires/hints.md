# Hints for A Cache That Expires

## Hint 1: Get Is a Lookup Plus a Time Comparison

Take the lock, read the entry, and decide whether it is still alive:

```go
e, ok := c.items[key]
if !ok {
    return "", false
}
```

## Hint 2: Compare Against the Stored Expiry

`Set` already recorded `expires`. An entry is alive while now is *before* it:

```go
if !time.Now().Before(e.expires) {
    return "", false
}
return e.value, true
```

Using `time.Now().After(e.expires)` is almost the same but treats the exact expiry
instant as still alive. Either passes the tests; `Before` is the tidier reading.

## Hint 3: Cleanup Is a Ticker and a Select

```go
tk := time.NewTicker(interval)
defer tk.Stop()
for {
    select {
    case <-ctx.Done():
        return
    case <-tk.C:
        // sweep
    }
}
```

`defer tk.Stop()` matters: a ticker that is never stopped keeps a timer alive.

## Hint 4: Deleting While Ranging Is Fine

Deleting from a map during `range` over that same map is explicitly allowed in Go:

```go
now := time.Now()
c.mu.Lock()
for k, e := range c.items {
    if !now.Before(e.expires) {
        delete(c.items, k)
    }
}
c.mu.Unlock()
```

Take `now` once outside the loop so every entry is judged against the same instant.

## Hint 5: If the Cleanup Test Hangs

The bubble only advances the clock when every goroutine is durably blocked. A janitor
written as `for { ... }` with no blocking operation never lets that happen, so the test
hangs instead of failing. Make sure the loop blocks on the ticker channel.

## Hint 6: If You See a Goroutine Leak Failure

`synctest.Test` fails if goroutines started inside the bubble are still running when
the body returns. That means `Cleanup` ignored `ctx.Done()`. The `select` above fixes
it.
