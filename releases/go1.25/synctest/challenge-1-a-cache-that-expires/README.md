# Challenge 1: A Cache That Expires

## The situation

`Cache` stores strings and every entry is supposed to die after a fixed TTL. `Set` and
`Len` are written. Two things are missing, and both of them are only interesting after
some time has passed, which is exactly the kind of code that used to be miserable to
test.

## The task

### `Get`

```go
func (c *Cache) Get(key string) (string, bool)
```

Return the value for `key`. Treat an entry that has passed its expiry time as a **miss**,
even though it is still sitting in the map. An unknown key is also a miss.

Take the lock. `Set` already stamps each entry with `expires`.

### `Cleanup`

```go
func (c *Cache) Cleanup(ctx context.Context, interval time.Duration)
```

A janitor. Every `interval`, delete the entries that have expired. Return when `ctx` is
done. It runs in its own goroutine, so it must not busy-loop and it must actually stop.

## What the tests do

This is the part worth reading, because the tests are the reason this challenge exists.

Every test body runs inside a bubble:

```go
synctest.Test(t, func(t *testing.T) {
    c := New(time.Hour)
    c.Set("k", "v")

    time.Sleep(time.Hour + time.Second)   // costs nothing

    if _, ok := c.Get("k"); ok {
        t.Fatal("entry is past its TTL")
    }
})
```

That `Sleep` is an hour of *virtual* time. The suite verifies TTL behaviour across an
hour, ten minutes and twenty-five hours, and finishes in about a quarter of a second.

`TestTheWholeSuiteCostsNoRealTime` checks this on purpose: it measures real elapsed
time around a 25 hour simulation and fails if more than two seconds actually passed.

The janitor tests use the other half of the package:

```go
go c.Cleanup(ctx, 10*time.Second)
synctest.Wait()   // returns once the janitor is blocked again
```

`Wait` does not move the clock. It returns once every *other* goroutine in the bubble
is durably blocked, which is how the test knows the janitor has finished a pass and
that asserting now is safe.

## Rules

- Compare against the stored expiry time. Do not round to whole seconds, one test
  checks a moment one second before expiry.
- `Cleanup` must return on cancellation. A bubble whose goroutines are still running
  when the body ends is a failure, so a janitor that ignores `ctx` shows up as a leak.
- Keep the mutex discipline. `Len` and `Set` already lock; your code should too.

## Run it

```bash
go test -v ./...
```
