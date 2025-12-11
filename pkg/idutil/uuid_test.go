package idutil

import (
	"testing"

	"github.com/apotourlyan/ludus-studii/pkg/testutil"
)

func TestUUID_Length(t *testing.T) {
	t.Parallel()

	uuid := UUID()

	// UUID format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx (36 characters)
	testutil.GotWant(t, len(uuid), 36)
}

func TestUUID_Version(t *testing.T) {
	t.Parallel()

	uuid := UUID()

	// Version 4 is indicated by '4' at position 14
	testutil.GotWant(t, uuid[14], '4')
}

func TestUUID_Variant(t *testing.T) {
	t.Parallel()

	uuid := UUID()

	// Variant bits should be 10xx in binary, which means first hex digit
	// at position 19 should be 8, 9, a, or b
	valid := []byte{'8', '9', 'a', 'b'}
	testutil.GotWantOneOf(t, uuid[19], valid)
}

func TestUUID_ValidCharacters(t *testing.T) {
	t.Parallel()

	uuid := UUID()

	valid := "0123456789abcdef-"
	for _, char := range uuid {
		testutil.GotWantOneOf(t, char, []rune(valid))
	}
}

func TestUUID_HyphenPositions(t *testing.T) {
	t.Parallel()

	uuid := UUID()

	// Hyphens should be at positions 8, 13, 18, and 23
	// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	testutil.GotWant(t, uuid[8], '-')
	testutil.GotWant(t, uuid[13], '-')
	testutil.GotWant(t, uuid[18], '-')
	testutil.GotWant(t, uuid[23], '-')
}

func TestUUID_Uniqueness(t *testing.T) {
	t.Parallel()

	iterations := 1000
	uuids := make(map[string]bool, iterations)

	for range iterations {
		uuid := UUID()
		uuids[uuid] = true
	}

	testutil.GotWant(t, len(uuids), iterations)
}
