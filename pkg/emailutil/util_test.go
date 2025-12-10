package emailutil

import (
	"strings"
	"testing"

	"github.com/apotourlyan/ludus-studii/pkg/testutil"
)

func TestParse_ValidEmails(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		email      string
		wantLocal  string
		wantDomain string
	}{
		{
			name:       "simple email",
			email:      "user@example.com",
			wantLocal:  "user",
			wantDomain: "example.com",
		},
		{
			name:       "email with subdomain",
			email:      "user@mail.example.com",
			wantLocal:  "user",
			wantDomain: "mail.example.com",
		},
		{
			name:       "email with plus sign",
			email:      "user+tag@example.com",
			wantLocal:  "user+tag",
			wantDomain: "example.com",
		},
		{
			name:       "email with dash in local",
			email:      "user-name@example.com",
			wantLocal:  "user-name",
			wantDomain: "example.com",
		},
		{
			name:       "email with underscore in local",
			email:      "user_name@example.com",
			wantLocal:  "user_name",
			wantDomain: "example.com",
		},
		{
			name:       "email with dots in local",
			email:      "user.name@example.com",
			wantLocal:  "user.name",
			wantDomain: "example.com",
		},
		{
			name:       "email with numbers in local",
			email:      "user123@example.com",
			wantLocal:  "user123",
			wantDomain: "example.com",
		},
		{
			name:       "email with percent in local",
			email:      "user%test@example.com",
			wantLocal:  "user%test",
			wantDomain: "example.com",
		},
		{
			name:       "complex local part",
			email:      "first.last+tag-123@example.com",
			wantLocal:  "first.last+tag-123",
			wantDomain: "example.com",
		},
		{
			name:       "domain with dash",
			email:      "user@ex-ample.com",
			wantLocal:  "user",
			wantDomain: "ex-ample.com",
		},
		{
			name:       "domain with multiple subdomains",
			email:      "user@sub1.sub2.example.com",
			wantLocal:  "user",
			wantDomain: "sub1.sub2.example.com",
		},
		{
			name:       "long TLD",
			email:      "user@example.museum",
			wantLocal:  "user",
			wantDomain: "example.museum",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			parts, err := Parse(c.email)
			testutil.DontWantError(t, err)
			testutil.DontWantNil(t, parts)
			testutil.GotWant(t, parts.Local, c.wantLocal)
			testutil.GotWant(t, parts.Domain, c.wantDomain)
		})
	}
}

func TestParse_InvalidEmails(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		email string
	}{
		{
			name:  "missing @",
			email: "userexample.com",
		},
		{
			name:  "missing local part",
			email: "@example.com",
		},
		{
			name:  "missing domain",
			email: "user@",
		},
		{
			name:  "empty string",
			email: "",
		},
		{
			name:  "consecutive dots in local",
			email: "user..name@example.com",
		},
		{
			name:  "starts with dot",
			email: ".user@example.com",
		},
		{
			name:  "ends with dot before @",
			email: "user.@example.com",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			parts, err := Parse(c.email)
			testutil.WantError(t, err)
			testutil.GotWant(t, err, ErrInvalidFormat)
			if parts != nil {
				t.Errorf("expected parts to be nil for invalid email %q, got Local=%q, Domain=%q",
					c.email, parts.Local, parts.Domain)
			}
		})
	}
}

func TestIsValid_ValidEmails(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		email string
	}{
		{
			name:  "simple email",
			email: "user@example.com",
		},
		{
			name:  "email with subdomain",
			email: "user@mail.example.com",
		},
		{
			name:  "email with plus sign",
			email: "user+tag@example.com",
		},
		{
			name:  "email with dash",
			email: "user-name@example.com",
		},
		{
			name:  "email with underscore",
			email: "user_name@example.com",
		},
		{
			name:  "email with dots",
			email: "user.name@example.com",
		},
		{
			name:  "email with numbers",
			email: "user123@example.com",
		},
		{
			name:  "email with percent",
			email: "user%test@example.com",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			got := IsValid(c.email)
			testutil.GotWant(t, got, true)
		})
	}
}

func TestIsValid_InvalidEmails(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		email string
	}{
		// Local part (before @) issues
		{
			name:  "consecutive dots in local",
			email: "user..name@example.com",
		},
		{
			name:  "starts with dot",
			email: ".user@example.com",
		},
		{
			name:  "ends with dot before @",
			email: "user.@example.com",
		},
		{
			name:  "local part too long (>64 chars)",
			email: strings.Repeat("a", 65) + "@example.com",
		},

		// Domain structure issues
		{
			name:  "domain starts with dash",
			email: "user@-example.com",
		},
		{
			name:  "domain ends with dash",
			email: "user@example-.com",
		},
		{
			name:  "consecutive dots in domain",
			email: "user@example..com",
		},
		{
			name:  "domain starts with dot",
			email: "user@.example.com",
		},
		{
			name:  "domain ends with dot",
			email: "user@example.com.",
		},
		{
			name:  "domain label too long (>63 chars)",
			email: "user@" + strings.Repeat("a", 64) + ".com",
		},
		{
			name: "total domain too long (>253 chars)",
			email: "user@" + strings.Repeat("a", 62) + "." +
				strings.Repeat("b", 62) + "." +
				strings.Repeat("c", 62) + "." +
				strings.Repeat("d", 62) + ".com",
		},

		// TLD issues
		{
			name:  "TLD too short (1 char)",
			email: "user@example.c",
		},
		{
			name:  "TLD with numbers only",
			email: "user@example.123",
		},
		{
			name:  "TLD starts with digit",
			email: "user@example.1com",
		},

		// Multiple @ symbols
		{
			name:  "double @ in local part",
			email: "user@@example.com",
		},
		{
			name:  "@ symbol in domain",
			email: "user@exam@ple.com",
		},
		{
			name:  "@ at start",
			email: "@user@example.com",
		},
		{
			name:  "@ at end",
			email: "user@example.com@",
		},

		// Missing parts
		{
			name:  "missing @",
			email: "userexample.com",
		},
		{
			name:  "missing local part",
			email: "@example.com",
		},
		{
			name:  "missing domain",
			email: "user@",
		},
		{
			name:  "missing TLD",
			email: "user@example",
		},

		// Invalid characters
		{
			name:  "contains space in local",
			email: "user name@example.com",
		},
		{
			name:  "contains space in domain",
			email: "user@exam ple.com",
		},
		{
			name:  "contains space in TLD",
			email: "user@example.co m",
		},
		{
			name:  "contains comma",
			email: "user,name@example.com",
		},
		{
			name:  "domain with invalid char",
			email: "user@exam$ple.com",
		},
		{
			name:  "TLD with invalid char",
			email: "user@example.co*m",
		},

		// Edge cases
		{
			name:  "empty string",
			email: "",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			got := IsValid(c.email)
			testutil.GotWant(t, got, false)
		})
	}
}
