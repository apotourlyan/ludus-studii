package testutil

import (
	"fmt"
	"strings"
	"testing"
)

func GotWant[T comparable](t *testing.T, got T, want T) {
	t.Helper()
	if got != want {
		text := getGotWantText(got, want)
		t.Error(text)
	}
}

func getGotWantText[T any](got T, want T) string {
	g := any(got)
	w := any(want)
	switch g.(type) {
	case int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8:
		return fmt.Sprintf("got %d, want %d\n", g, w)
	case float64, float32:
		return fmt.Sprintf("got %f, want %f\n", g, w)
	case bool:
		return fmt.Sprintf("got %t, want %t\n", g, w)
	case string:
		return fmt.Sprintf("got %q, want %q\n", g, w)
	default:
		return fmt.Sprintf("got %#v, want %#v\n", g, w)
	}
}

func DontWant[T comparable](t *testing.T, got T, value T) {
	t.Helper()
	if got == value {
		text := getDontWantText(value)
		t.Error(text)
	}
}

func getDontWantText[T any](value T) string {
	v := any(value)
	switch v.(type) {
	case int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8:
		return fmt.Sprintf("don't want equal to %d\n", v)
	case float64, float32:
		return fmt.Sprintf("don't want equal to %f\n", v)
	case bool:
		return fmt.Sprintf("don't want equal to %t\n", v)
	case string:
		return fmt.Sprintf("don't want equal to %q\n", v)
	default:
		return fmt.Sprintf("don't want equal to%#v\n", v)
	}
}

func GotWantPanic(t *testing.T, f func(), want string) {
	t.Helper()
	panicked, got := catchPanic(f)
	if !panicked {
		t.Errorf("got panic 'nil', want panic %q", want)
	} else if got != want {
		t.Errorf("got panic %q, want panic %q", got, want)
	}
}

func catchPanic(f func()) (panicked bool, message string) {
	defer func() {
		if r := recover(); r != nil {
			panicked = true
			message = fmt.Sprint(r)
		}
	}()
	f()
	return false, ""
}

func WantError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Errorf("got error 'nil', want error")
	}
}

func DontWantError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func DontWantNil(t *testing.T, value any) {
	t.Helper()
	if value == nil {
		t.Fatal("expected non-nil value")
	}
}

func WantPrefix(t *testing.T, got string, prefixes ...string) {
	t.Helper()
	startsWithOneOf := false

	for _, prefix := range prefixes {
		if strings.HasPrefix(got, prefix) {
			startsWithOneOf = true
			break
		}
	}

	if startsWithOneOf {
		return
	}

	if len(prefixes) == 1 {
		t.Errorf("got %q, want prefix %v", got, prefixes[0])
	} else {
		t.Errorf("got %q, want prefix one of %v", got, prefixes)
	}
}
