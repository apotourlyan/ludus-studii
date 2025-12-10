package syncutil

import (
	"sync"
	"testing"

	"github.com/apotourlyan/ludus-studii/pkg/testutil"
)

func TestCounter_Next_Success(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		shouldReset bool
		callCount   int
		want        int64
	}{
		{
			name:        "one call with reset",
			shouldReset: true,
			callCount:   1,
			want:        0,
		},
		{
			name:        "increment without reset",
			shouldReset: false,
			callCount:   1,
			want:        1,
		},
		{
			name:        "multiple increments",
			shouldReset: false,
			callCount:   5,
			want:        5,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			counter := NewCounter()
			var result int64
			for i := 0; i < c.callCount; i++ {
				result = counter.Next(func() bool {
					return c.shouldReset
				})
			}
			testutil.GotWant(t, result, c.want)
		})
	}
}

func TestCounter_Next_Reset(t *testing.T) {
	t.Parallel()

	t.Run("reset after increments", func(t *testing.T) {
		t.Parallel()

		counter := NewCounter()

		// Increment a few times
		counter.Next(func() bool { return false })
		counter.Next(func() bool { return false })
		result := counter.Next(func() bool { return false })
		testutil.GotWant(t, result, int64(3))

		// Reset
		result = counter.Next(func() bool { return true })
		testutil.GotWant(t, result, int64(0))

		// Continue incrementing
		result = counter.Next(func() bool { return false })
		testutil.GotWant(t, result, int64(1))
	})

	t.Run("multiple resets", func(t *testing.T) {
		t.Parallel()

		counter := NewCounter()
		for range 3 {
			result := counter.Next(func() bool { return true })
			testutil.GotWant(t, result, int64(0))
		}
	})
}

func TestCounter_Next_Concurrent(t *testing.T) {
	t.Parallel()

	counter := NewCounter()
	var wg sync.WaitGroup
	goroutines := 100
	callsPerGoroutine := 100

	wg.Add(goroutines)
	for range goroutines {
		go func() {
			defer wg.Done()
			for range callsPerGoroutine {
				counter.Next(func() bool { return false })
			}
		}()
	}
	wg.Wait()

	// One more call to get the final count
	got := counter.Next(func() bool { return false })
	want := int64(goroutines*callsPerGoroutine + 1)
	testutil.GotWant(t, got, want)
}
