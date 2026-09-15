# Challenge 1: Fan Out Without the Bookkeeping

## The task

Two functions, both a fan-out, both written with `wg.Go` rather than `Add` and `Done`.

### `FetchAll`

```go
func FetchAll(ids []int, fetch func(int) string) []string
```

Call `fetch` for every id **concurrently** and return the results in the **same order
as `ids`**, no matter what order the goroutines finish in.

### `CountMatching`

```go
func CountMatching(inputs []string, match func(string) bool) int
```

Call `match` for every input concurrently and return how many returned `true`.

## The two things worth getting right

**Order comes from the index, not from completion.** Appending to a shared slice as
each goroutine finishes gives you completion order, and needs a mutex on top. Writing
to `out[i]` needs neither: every goroutine owns one element, so there is no shared
write at all.

**A counter is shared.** Unlike the slice, `CountMatching` really does have every
goroutine touching the same value, so `count++` is a data race that quietly loses
increments. Use `sync/atomic`.

## What the tests do

| Test | What it catches |
|---|---|
| `TestFetchAllKeepsInputOrder` | the slowest fetch is first, so appending returns it last |
| `TestFetchAllActuallyRunsConcurrently` | a barrier that only releases when all four are in flight |
| `TestCountMatchingUnderContention` | 500 goroutines; an unsynchronised counter loses some |
| `TestCountMatchingWaitsForEveryGoroutine` | returning before `Wait` shows up as a short count |

The concurrency test is the interesting one. Every `fetch` blocks until all of them
have arrived, so a loop that calls `fetch` one at a time can never reach the fourth
call and the test fails on its timeout rather than passing by accident.

## Rules

- Use `wg.Go`. That is the point of the exercise, and it is shorter.
- Do not add a mutex around the result slice. If you need one, you are appending.
- Both functions must have finished all their work before they return.

## Run it

```bash
go test -v ./...
```

Worth also running with the race detector, which is how you would catch the counter
bug in real code:

```bash
go test -race ./...
```
