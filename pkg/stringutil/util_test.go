package stringutil

import (
	"testing"

	"github.com/apotourlyan/ludus-studii/pkg/testutil"
)

func TestContainsUppercase_WithUppercase(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		text string
	}{
		{
			name: "uppercase at start",
			text: "Hello",
		},
		{
			name: "uppercase in middle",
			text: "heLlo",
		},
		{
			name: "uppercase at end",
			text: "helloW",
		},
		{
			name: "all uppercase",
			text: "HELLO",
		},
		{
			name: "single uppercase letter",
			text: "A",
		},
		{
			name: "mixed with digits",
			text: "Test123",
		},
		{
			name: "mixed with special chars",
			text: "Hello!",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			got := ContainsUppercase(c.text)
			testutil.GotWant(t, got, true)
		})
	}
}

func TestContainsUppercase_WithoutUppercase(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		text string
	}{
		{
			name: "all lowercase",
			text: "hello",
		},
		{
			name: "digits only",
			text: "12345",
		},
		{
			name: "special chars only",
			text: "!@#$%",
		},
		{
			name: "empty string",
			text: "",
		},
		{
			name: "lowercase with digits",
			text: "test123",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			got := ContainsUppercase(c.text)
			testutil.GotWant(t, got, false)
		})
	}
}

func TestContainsLowercase_WithLowercase(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		text string
	}{
		{
			name: "lowercase at start",
			text: "hello",
		},
		{
			name: "lowercase in middle",
			text: "HeLLo",
		},
		{
			name: "lowercase at end",
			text: "HELLOw",
		},
		{
			name: "all lowercase",
			text: "hello",
		},
		{
			name: "single lowercase letter",
			text: "a",
		},
		{
			name: "mixed with digits",
			text: "test123",
		},
		{
			name: "mixed with special chars",
			text: "hello!",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			got := ContainsLowercase(c.text)
			testutil.GotWant(t, got, true)
		})
	}
}

func TestContainsLowercase_WithoutLowercase(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		text string
	}{
		{
			name: "all uppercase",
			text: "HELLO",
		},
		{
			name: "digits only",
			text: "12345",
		},
		{
			name: "special chars only",
			text: "!@#$%",
		},
		{
			name: "empty string",
			text: "",
		},
		{
			name: "uppercase with digits",
			text: "TEST123",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			got := ContainsLowercase(c.text)
			testutil.GotWant(t, got, false)
		})
	}
}

func TestContainsDigit_WithDigit(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		text string
	}{
		{
			name: "digit at start",
			text: "1hello",
		},
		{
			name: "digit in middle",
			text: "hel2lo",
		},
		{
			name: "digit at end",
			text: "hello3",
		},
		{
			name: "all digits",
			text: "12345",
		},
		{
			name: "single digit",
			text: "0",
		},
		{
			name: "multiple digits",
			text: "abc123def456",
		},
		{
			name: "password with digit",
			text: "Password1",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			got := ContainsDigit(c.text)
			testutil.GotWant(t, got, true)
		})
	}
}

func TestContainsDigit_WithoutDigit(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		text string
	}{
		{
			name: "letters only",
			text: "hello",
		},
		{
			name: "special chars only",
			text: "!@#$%",
		},
		{
			name: "empty string",
			text: "",
		},
		{
			name: "mixed letters and special chars",
			text: "Hello!World",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			got := ContainsDigit(c.text)
			testutil.GotWant(t, got, false)
		})
	}
}

func TestContainsSpecial_WithSpecialChar(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		text string
	}{
		{
			name: "exclamation mark",
			text: "hello!",
		},
		{
			name: "at symbol",
			text: "user@domain",
		},
		{
			name: "hash",
			text: "hello#world",
		},
		{
			name: "dollar sign",
			text: "$100",
		},
		{
			name: "percent",
			text: "50%",
		},
		{
			name: "caret",
			text: "x^2",
		},
		{
			name: "ampersand",
			text: "you&me",
		},
		{
			name: "asterisk",
			text: "test*",
		},
		{
			name: "parentheses",
			text: "test(1)",
		},
		{
			name: "square brackets",
			text: "array[0]",
		},
		{
			name: "curly braces",
			text: "obj{key}",
		},
		{
			name: "dash",
			text: "test-case",
		},
		{
			name: "underscore",
			text: "test_case",
		},
		{
			name: "plus",
			text: "1+1",
		},
		{
			name: "equals",
			text: "a=b",
		},
		{
			name: "pipe",
			text: "a|b",
		},
		{
			name: "semicolon",
			text: "a;b",
		},
		{
			name: "colon",
			text: "key:value",
		},
		{
			name: "comma",
			text: "a,b,c",
		},
		{
			name: "period",
			text: "end.",
		},
		{
			name: "angle brackets",
			text: "<tag>",
		},
		{
			name: "question mark",
			text: "really?",
		},
		{
			name: "password with special char",
			text: "Password123!",
		},
		{
			name: "multiple special chars",
			text: "!@#$%^&*()",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			got := ContainsSpecial(c.text)
			testutil.GotWant(t, got, true)
		})
	}
}

func TestContainsSpecial_WithoutSpecialChar(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		text string
	}{
		{
			name: "letters only",
			text: "hello",
		},
		{
			name: "digits only",
			text: "12345",
		},
		{
			name: "empty string",
			text: "",
		},
		{
			name: "letters and digits",
			text: "Hello123",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			got := ContainsSpecial(c.text)
			testutil.GotWant(t, got, false)
		})
	}
}

func TestIsWhitespace_WithWhitespace(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		text string
	}{
		{
			name: "single space",
			text: " ",
		},
		{
			name: "multiple spaces",
			text: "   ",
		},
		{
			name: "tab",
			text: "\t",
		},
		{
			name: "newline",
			text: "\n",
		},
		{
			name: "carriage return",
			text: "\r",
		},
		{
			name: "mixed whitespace",
			text: " \t\n\r",
		},
		{
			name: "multiple tabs",
			text: "\t\t\t",
		},
		{
			name: "multiple newlines",
			text: "\n\n\n",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			got := IsWhitespace(c.text)
			testutil.GotWant(t, got, true)
		})
	}
}

func TestIsWhitespace_WithoutWhitespace(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		text string
	}{
		{
			name: "empty string",
			text: "",
		},
		{
			name: "single letter",
			text: "a",
		},
		{
			name: "word",
			text: "hello",
		},
		{
			name: "letter with spaces",
			text: " a ",
		},
		{
			name: "word with leading space",
			text: " hello",
		},
		{
			name: "word with trailing space",
			text: "hello ",
		},
		{
			name: "digit",
			text: "1",
		},
		{
			name: "special char",
			text: "!",
		},
		{
			name: "spaces with letter",
			text: "  a  ",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			got := IsWhitespace(c.text)
			testutil.GotWant(t, got, false)
		})
	}
}

func TestCapitalize(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		text string
		want string
	}{
		{
			name: "lowercase word",
			text: "hello",
			want: "Hello",
		},
		{
			name: "uppercase word",
			text: "HELLO",
			want: "HELLO",
		},
		{
			name: "already capitalized",
			text: "Hello",
			want: "Hello",
		},
		{
			name: "single lowercase letter",
			text: "a",
			want: "A",
		},
		{
			name: "single uppercase letter",
			text: "A",
			want: "A",
		},
		{
			name: "empty string",
			text: "",
			want: "",
		},
		{
			name: "starts with digit",
			text: "1hello",
			want: "1hello",
		},
		{
			name: "starts with special char",
			text: "!hello",
			want: "!hello",
		},
		{
			name: "starts with space",
			text: " hello",
			want: " hello",
		},
		{
			name: "mixed case",
			text: "hELLO",
			want: "HELLO",
		},
		{
			name: "with punctuation",
			text: "hello!",
			want: "Hello!",
		},
		{
			name: "multiple words",
			text: "hello world",
			want: "Hello world",
		},
		{
			name: "unicode lowercase",
			text: "über",
			want: "Über",
		},
		{
			name: "unicode uppercase",
			text: "Über",
			want: "Über",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			got := Capitalize(c.text)
			testutil.GotWant(t, got, c.want)
		})
	}
}
