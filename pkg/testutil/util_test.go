package testutil

import "testing"

func TestCatchPanic(t *testing.T) {
	cases := []struct {
		name         string
		f            func()
		wantPanicked bool
		wantMessage  string
	}{
		{
			name: "no panic",
			f: func() {
				// no panic
			},
			wantPanicked: false,
			wantMessage:  "",
		},
		{
			name: "string panic",
			f: func() {
				panic("test error")
			},
			wantPanicked: true,
			wantMessage:  "test error",
		},
		{
			name: "integer panic",
			f: func() {
				panic(42)
			},
			wantPanicked: true,
			wantMessage:  "42",
		},
		{
			name: "struct panic",
			f: func() {
				panic(struct{ msg string }{"custom error"})
			},
			wantPanicked: true,
			wantMessage:  "{custom error}",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			panicked, message := catchPanic(c.f)
			GotWant(t, panicked, c.wantPanicked)
			GotWant(t, message, c.wantMessage)
		})
	}
}
