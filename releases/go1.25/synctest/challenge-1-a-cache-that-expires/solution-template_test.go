package main

import (
	"context"
	"testing"
	"testing/synctest"
	"time"
)

func TestGetReturnsAFreshEntry(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c := New(time.Minute)
		c.Set("k", "v")

		got, ok := c.Get("k")
		if !ok {
			t.Fatal("Get missed an entry that has not expired yet")
		}
		if got != "v" {
			t.Fatalf("Get = %q, want %q", got, "v")
		}
	})
}

func TestGetMissesAnUnknownKey(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c := New(time.Minute)

		if _, ok := c.Get("nope"); ok {
			t.Fatal("Get reported a hit for a key that was never set")
		}
	})
}

// TestGetMissesAfterTheTTL is the test that used to need a real sleep.
// Inside the bubble the hour passes instantly.
func TestGetMissesAfterTheTTL(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c := New(time.Hour)
		c.Set("k", "v")

		time.Sleep(time.Hour + time.Second)

		if _, ok := c.Get("k"); ok {
			t.Fatal("Get returned an entry that is past its TTL")
		}
	})
}

func TestGetStillHitsJustBeforeTheTTL(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c := New(time.Hour)
		c.Set("k", "v")

		time.Sleep(time.Hour - time.Second)

		if _, ok := c.Get("k"); !ok {
			t.Fatal("Get expired an entry one second early; compare against the expiry time, not a rounded value")
		}
	})
}

// TestTheWholeSuiteCostsNoRealTime proves the point of the package.
func TestTheWholeSuiteCostsNoRealTime(t *testing.T) {
	realStart := time.Now()

	synctest.Test(t, func(t *testing.T) {
		c := New(24 * time.Hour)
		c.Set("k", "v")

		time.Sleep(25 * time.Hour)

		if _, ok := c.Get("k"); ok {
			t.Fatal("entry survived 25 virtual hours with a 24 hour TTL")
		}
	})

	if elapsed := time.Since(realStart); elapsed > 2*time.Second {
		t.Fatalf("25 virtual hours took %v of real time; the body is not running in a bubble", elapsed)
	}
}

func TestCleanupRemovesExpiredEntries(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c := New(time.Minute)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		c.Set("a", "1")
		c.Set("b", "2")

		go c.Cleanup(ctx, 10*time.Second)

		// Let the janitor reach its first tick and block again.
		synctest.Wait()
		if c.Len() != 2 {
			t.Fatalf("entries were removed before they expired: Len = %d, want 2", c.Len())
		}

		// Past the TTL, the janitor should clear them out.
		time.Sleep(2 * time.Minute)
		synctest.Wait()

		if n := c.Len(); n != 0 {
			t.Fatalf("Cleanup left %d expired entries behind, want 0", n)
		}
	})
}

func TestCleanupKeepsEntriesThatAreStillFresh(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c := New(time.Hour)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		c.Set("keep", "me")
		go c.Cleanup(ctx, time.Minute)

		time.Sleep(10 * time.Minute)
		synctest.Wait()

		if n := c.Len(); n != 1 {
			t.Fatalf("Cleanup removed a live entry: Len = %d, want 1", n)
		}
	})
}

// TestCleanupStopsWhenTheContextIsDone also catches a goroutine leak: a bubble
// whose goroutines are still running when the body returns is an error.
func TestCleanupStopsWhenTheContextIsDone(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c := New(time.Minute)
		ctx, cancel := context.WithCancel(context.Background())

		done := make(chan struct{})
		go func() {
			c.Cleanup(ctx, time.Second)
			close(done)
		}()

		synctest.Wait()
		cancel()

		select {
		case <-done:
		case <-time.After(time.Minute):
			t.Fatal("Cleanup did not return after the context was cancelled")
		}
	})
}
