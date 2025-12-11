package syncutil

import (
	"fmt"
	"sync"
	"testing"

	"github.com/apotourlyan/ludus-studii/pkg/testutil"
)

func TestConcurrentMap_SetGet(t *testing.T) {
	t.Parallel()

	m := NewConcurrentMap[string, int]()
	m.Set("key1", 100)

	got, ok := m.Get("key1")
	testutil.GotWant(t, ok, true)
	testutil.GotWant(t, got, 100)
}

func TestConcurrentMap_OverwriteValue(t *testing.T) {
	t.Parallel()

	m := NewConcurrentMap[string, int]()
	m.Set("counter", 1)
	m.Set("counter", 2)

	got, ok := m.Get("counter")
	testutil.GotWant(t, ok, true)
	testutil.GotWant(t, got, 2)
}

func TestConcurrentMap_MultipleKeys(t *testing.T) {
	t.Parallel()

	m := NewConcurrentMap[string, int]()
	m.Set("key1", 1)
	m.Set("key2", 2)
	m.Set("key3", 3)

	got, ok := m.Get("key1")
	testutil.GotWant(t, ok, true)
	testutil.GotWant(t, got, 1)

	got, ok = m.Get("key2")
	testutil.GotWant(t, ok, true)
	testutil.GotWant(t, got, 2)

	got, ok = m.Get("key3")
	testutil.GotWant(t, ok, true)
	testutil.GotWant(t, got, 3)
}

func TestConcurrentMap_ZeroValue(t *testing.T) {
	t.Parallel()

	m := NewConcurrentMap[string, int]()
	m.Set("zero", 0)

	got, ok := m.Get("zero")
	testutil.GotWant(t, ok, true)
	testutil.GotWant(t, got, 0)

	got, ok = m.Get("nonexistent")
	testutil.GotWant(t, ok, false)
	testutil.GotWant(t, got, 0)
}

func TestConcurrentMap_GetFromEmptyMap(t *testing.T) {
	t.Parallel()

	m := NewConcurrentMap[string, int]()
	_, ok := m.Get("nonexistent")
	testutil.GotWant(t, ok, false)
}

func TestConcurrentMap_GetNonexistentKey(t *testing.T) {
	t.Parallel()

	m := NewConcurrentMap[string, int]()
	m.Set("exists", 123)

	_, ok := m.Get("nonexistent")
	testutil.GotWant(t, ok, false)
}

func TestConcurrentMap_RemoveExistingKey(t *testing.T) {
	t.Parallel()

	m := NewConcurrentMap[string, int]()
	m.Set("key1", 100)

	_, ok := m.Get("key1")
	testutil.GotWant(t, ok, true)

	m.Remove("key1")

	_, ok = m.Get("key1")
	testutil.GotWant(t, ok, false)
}

func TestConcurrentMap_RemoveNonexistentKey(t *testing.T) {
	t.Parallel()

	m := NewConcurrentMap[string, int]()
	m.Remove("nonexistent") // Should not panic
}

func TestConcurrentMap_RemoveFromEmptyMap(t *testing.T) {
	t.Parallel()

	m := NewConcurrentMap[string, int]()
	m.Remove("key") // Should not panic
}

func TestConcurrentMap_RemoveOneOfMultipleKeys(t *testing.T) {
	t.Parallel()

	m := NewConcurrentMap[string, int]()
	m.Set("key1", 1)
	m.Set("key2", 2)
	m.Set("key3", 3)

	m.Remove("key2")

	got, ok := m.Get("key1")
	testutil.GotWant(t, ok, true)
	testutil.GotWant(t, got, 1)

	_, ok = m.Get("key2")
	testutil.GotWant(t, ok, false)

	got, ok = m.Get("key3")
	testutil.GotWant(t, ok, true)
	testutil.GotWant(t, got, 3)
}

func TestConcurrentMap_RemoveAndReAddKey(t *testing.T) {
	t.Parallel()

	m := NewConcurrentMap[string, int]()
	m.Set("key", 1)
	m.Remove("key")
	m.Set("key", 2)

	got, ok := m.Get("key")
	testutil.GotWant(t, ok, true)
	testutil.GotWant(t, got, 2)
}

func TestConcurrentMap_ConcurrentSetThenGet(t *testing.T) {
	t.Parallel()

	m := NewConcurrentMap[string, int]()
	var wg sync.WaitGroup
	goroutines := 100
	ops := 100

	wg.Add(goroutines)
	for i := range goroutines {
		go func(id int) {
			defer wg.Done()
			for j := range ops {
				key := fmt.Sprintf("key-%d-%d", id, j)
				m.Set(key, id*ops+j)
			}
		}(i)
	}
	wg.Wait()

	for i := range goroutines {
		for j := range ops {
			key := fmt.Sprintf("key-%d-%d", i, j)
			got, ok := m.Get(key)
			testutil.GotWant(t, ok, true)
			testutil.GotWant(t, got, i*ops+j)
		}
	}
}

func TestConcurrentMap_ConcurrentSetThenRemove(t *testing.T) {
	t.Parallel()

	m := NewConcurrentMap[string, int]()
	var wg sync.WaitGroup
	goroutines := 50
	ops := 100

	// Concurrent setters
	wg.Add(goroutines)
	for i := range goroutines {
		go func(id int) {
			defer wg.Done()
			for j := range ops {
				key := fmt.Sprintf("key-%d-%d", id, j)
				m.Set(key, id*ops+j)
			}
		}(i)
	}
	wg.Wait()

	// Concurrent removers
	wg.Add(goroutines)
	for i := range goroutines {
		go func(id int) {
			defer wg.Done()
			for j := range ops {
				key := fmt.Sprintf("key-%d-%d", id, j)
				m.Remove(key)
			}
		}(i)
	}
	wg.Wait()

	for i := range goroutines {
		for j := range ops {
			key := fmt.Sprintf("key-%d-%d", i, j)
			_, ok := m.Get(key)
			testutil.GotWant(t, ok, false)
		}
	}
}

func TestConcurrentMap_ConcurrentReaders(t *testing.T) {
	t.Parallel()

	m := NewConcurrentMap[string, int]()
	m.Set("shared", 42)

	var wg sync.WaitGroup
	goroutines := 100
	ops := 100

	wg.Add(goroutines)
	for range goroutines {
		go func() {
			defer wg.Done()
			for range ops {
				got, ok := m.Get("shared")
				testutil.GotWant(t, ok, true)
				testutil.GotWant(t, got, 42)
			}
		}()
	}
	wg.Wait()
}

func TestConcurrentMap_ConcurrentWriters(t *testing.T) {
	t.Parallel()

	m := NewConcurrentMap[string, int]()
	var wg sync.WaitGroup
	goroutines := 100

	wg.Add(goroutines)
	for i := range goroutines {
		go func(value int) {
			defer wg.Done()
			m.Set("counter", value)
		}(i)
	}
	wg.Wait()

	// Should have one of the values
	got, ok := m.Get("counter")
	testutil.GotWant(t, ok, true)
	testutil.GotWantInRange(t, got, 0, 99)
}

func TestConcurrentMap_ConcurrentReaderWriters(t *testing.T) {
	t.Parallel()

	m := NewConcurrentMap[string, int]()
	var wg sync.WaitGroup
	goroutines := 100

	// Concurrent sets, gets, and removes on same key
	wg.Add(goroutines * 3)
	for range goroutines {
		go func() {
			defer wg.Done()
			m.Set("key", 1)
		}()
		go func() {
			defer wg.Done()
			m.Get("key")
		}()
		go func() {
			defer wg.Done()
			m.Remove("key")
		}()
	}

	wg.Wait()
	// Test completes without panics or deadlocks
}
