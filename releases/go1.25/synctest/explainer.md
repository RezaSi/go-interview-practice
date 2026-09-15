## The sleep everyone has written

Here is a test for something that happens after a delay. You have written this test.
Possibly this week.

```go
func TestCacheExpires(t *testing.T) {
	c := New(50 * time.Millisecond)
	c.Set("k", "v")

	time.Sleep(60 * time.Millisecond)

	if _, ok := c.Get("k"); ok {
		t.Fatal("entry should have expired")
	}
}
```

It works on your laptop. It costs 60 milliseconds. Multiply that by a few hundred
tests and your suite is slow for no good reason.

Worse, it is a coin flip. On a loaded CI runner the goroutine that does the expiry may
not get scheduled in those extra 10 milliseconds, and the test fails for a reason that
has nothing to do with your code. So somebody bumps the sleep to 100ms, then 200ms,
and now the suite is slower *and* still occasionally red.

The real problem is that the test has to guess how long to wait, because it has no way
to ask "is everything finished yet?"

## What Go 1.25 added

`testing/synctest` graduated from experiment to the standard library. You wrap the
test body in a **bubble**:

```go
func TestCacheExpires(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c := New(50 * time.Millisecond)
		c.Set("k", "v")

		time.Sleep(60 * time.Millisecond)

		if _, ok := c.Get("k"); ok {
			t.Fatal("entry should have expired")
		}
	})
}
```

Same test. Two differences: it takes **zero** real time, and it gives the same answer
every single run.

Inside a bubble, `time` is fake. `time.Now`, `time.Sleep`, `time.After`, `time.Ticker`,
`time.Timer` and everything built on them, including `context.WithTimeout`, read from a
clock the bubble owns. That clock does not tick along with the wall. It advances only
when there is nothing else to do.

Here is a five second context deadline, verified on a real toolchain:

```go
synctest.Test(t, func(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	start := time.Now()
	<-ctx.Done()
	t.Logf("elapsed=%v err=%v", time.Since(start), ctx.Err())
})
```

```text
elapsed=5s err=context deadline exceeded
--- PASS (0.00s)
```

Five virtual seconds. Zero real seconds.

## The rule that makes it deterministic

This is the part to internalise, because everything else follows from it.

**The clock advances only when every goroutine in the bubble is durably blocked.**

"Durably blocked" means blocked in a way that only another goroutine *in the bubble*
could release: a channel operation, a mutex, a `WaitGroup`, a `time.Sleep`. If any
goroutine is still runnable, the clock is frozen and that goroutine runs first.

So a delay never races against your test. Either work is still happening, in which case
time stands still, or everything is waiting, in which case time jumps straight to
whatever the next timer needs. There is no window in which the scheduler decides your
result.

That is the flakiness gone, and it is gone by construction rather than by picking a
bigger number.

## `Wait` is a different tool

`synctest.Wait()` is the second half of the package and it is regularly
misunderstood, so it is worth being precise:

**`Wait` does not advance the clock. It blocks until every *other* goroutine in the
bubble is durably blocked.**

Use it when you want to assert on state after the background work has settled, without
knowing how long that takes:

```go
go worker()      // does some work, then blocks on a channel
synctest.Wait()  // returns once worker can make no further progress
// now assert
```

A worked example that shows the distinction, and its real output:

```go
synctest.Test(t, func(t *testing.T) {
	ch := make(chan int)
	done := false

	go func() {
		time.Sleep(time.Second)
		done = true
		ch <- 1
	}()

	synctest.Wait()
	t.Logf("after Wait, done=%v", done) // done=false
	<-ch
})
```

`Wait` returned while the goroutine was still parked in `time.Sleep`, because sleeping
counts as durably blocked. `done` is still false. Only when the main goroutine also
blocks, on `<-ch`, is the whole bubble idle, so the clock jumps a second, the sleep
finishes and the send completes.

Two questions, two tools. *Has everything settled?* is `Wait`. *Has enough time
passed?* happens on its own.

## Where the bubble ends

A bubble is only deterministic over things it controls. Its walls are worth knowing.

- **Goroutines started inside are in it.** Goroutines started before it are not.
- **Real I/O is still real.** A network call or a disk read blocks on the outside
  world, which the bubble cannot fake and does not count as durably blocked.
- **Channels shared with the outside break the illusion.** If a goroutine inside is
  waiting on a channel that only an outside goroutine will send to, the bubble cannot
  tell whether it is stuck.
- **The bubble must drain.** If goroutines are still running when the function
  returns, that is an error, not a warning. This is a feature: it catches the goroutine
  leak you did not know you had.

The practical consequence is that `synctest` pushes you toward code that takes its
dependencies as parameters. Which is the design you wanted anyway, now with a test
that rewards you for it.

## What it is really for

Anything where the assertion is "after some time, X should have happened":

- caches and TTLs
- retry and backoff loops, where testing "does it back off correctly" used to mean
  actually waiting
- context deadlines and cancellation propagation
- rate limiters
- heartbeats, reconnection, anything with a ticker
- graceful shutdown paths with a timeout

If your test file imports `time` and calls `Sleep`, it is a candidate.

> **Try it:** the challenge below is the cache above, for real. You will write the
> expiry test twice, once with a sleep and once in a bubble, and the tests check both
> that the behaviour is correct and that the bubble version does not actually wait.
