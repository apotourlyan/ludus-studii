package register

import (
	"testing"

	"github.com/apotourlyan/ludus-studii/pkg/errorutil"
	"github.com/apotourlyan/ludus-studii/pkg/errorutil/errcode"
	"github.com/apotourlyan/ludus-studii/pkg/errorutil/errtext"
	"github.com/apotourlyan/ludus-studii/pkg/testutil"
	rerrcode "github.com/apotourlyan/ludus-studii/services/user/internal/service/user/register/errcode"
	rerrtext "github.com/apotourlyan/ludus-studii/services/user/internal/service/user/register/errtext"
	"github.com/apotourlyan/ludus-studii/services/user/internal/service/user/register/field"
)

func TestValidate_Success(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		email    string
		password string
	}{
		{
			name:     "valid credentials",
			email:    "test@example.com",
			password: "Password123!",
		},
		{
			name:     "valid with spaces trimmed",
			email:    "  test@example.com  ",
			password: "  Password123!  ",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			req := &Request{
				Email:    c.email,
				Password: c.password,
			}

			err := Validate(req)
			testutil.DontWantError(t, err)
		})
	}
}

func TestValidate_TrimmsInput(t *testing.T) {
	t.Parallel()

	req := &Request{
		Email:    "  test@example.com  ",
		Password: "  Password123!  ",
	}

	err := Validate(req)
	testutil.DontWantError(t, err)

	// Verify trimming happened
	testutil.GotWant(t, req.Email, "test@example.com")
	testutil.GotWant(t, req.Password, "Password123!")
}

func TestValidate_EmailErrors(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		email    string
		wantCode string
		wantText string
	}{
		{
			name:     "required",
			email:    "",
			wantCode: errcode.Required,
			wantText: errtext.Required(field.Email),
		},
		{
			name:     "format no @ symbol",
			email:    "invalid.email.com",
			wantCode: errcode.Format,
			wantText: errtext.Format(field.Email),
		},
		{
			name:     "format no domain",
			email:    "invalid@",
			wantCode: errcode.Format,
			wantText: errtext.Format(field.Email),
		},
		{
			name:     "format no local part",
			email:    "@example.com",
			wantCode: errcode.Format,
			wantText: errtext.Format(field.Email),
		},
		{
			name:     "format no extension",
			email:    "test@example",
			wantCode: errcode.Format,
			wantText: errtext.Format(field.Email),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			req := &Request{
				Email:    c.email,
				Password: "Password123!",
			}

			err := Validate(req)
			testutil.WantError(t, err)

			verr, ok := err.(*errorutil.ValidationError)
			testutil.GotWant(t, ok, true)

			got := verr.Has(field.Email, c.wantCode, c.wantText)
			testutil.GotWant(t, got, true)
		})
	}
}

func TestValidate_PasswordErrors(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		password string
		wantCode string
		wantText string
	}{
		{
			name:     "required",
			password: "",
			wantCode: errcode.Required,
			wantText: errtext.Required(field.Password),
		},
		{
			name:     "length",
			password: "Pass1!",
			wantCode: errcode.StringLength,
			wantText: errtext.StringLength(field.Password, 8),
		},
		{
			name:     "missing uppercase",
			password: "password123!",
			wantCode: rerrcode.PasswordUpper,
			wantText: rerrtext.PasswordUpper,
		},
		{
			name:     "missing lowercase",
			password: "PASSWORD123!",
			wantCode: rerrcode.PasswordLower,
			wantText: rerrtext.PasswordLower,
		},
		{
			name:     "missing digit",
			password: "Password!",
			wantCode: rerrcode.PasswordDigit,
			wantText: rerrtext.PasswordDigit,
		},
		{
			name:     "missing special",
			password: "Password123",
			wantCode: rerrcode.PasswordSpecial,
			wantText: rerrtext.PasswordSpecial,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			req := &Request{
				Email:    "test@example.com",
				Password: c.password,
			}

			err := Validate(req)
			testutil.WantError(t, err)

			verr, ok := err.(*errorutil.ValidationError)
			testutil.GotWant(t, ok, true)

			got := verr.Has(field.Password, c.wantCode, c.wantText)
			testutil.GotWant(t, got, true)
		})
	}
}

func TestValidate_EmailAndPasswordRequired(t *testing.T) {
	t.Parallel()

	req := &Request{
		Email:    "",
		Password: "",
	}

	err := Validate(req)
	testutil.WantError(t, err)

	verr, ok := err.(*errorutil.ValidationError)
	testutil.GotWant(t, ok, true)

	got := verr.Has(field.Email, errcode.Required, errtext.Required(field.Email))
	testutil.GotWant(t, got, true)

	got = verr.Has(field.Password, errcode.Required, errtext.Required(field.Password))
	testutil.GotWant(t, got, true)
}

func TestValidate_EmailFormatAndPasswordLength(t *testing.T) {
	t.Parallel()

	req := &Request{
		Email:    "invalid-email",
		Password: "Pass1!",
	}

	err := Validate(req)
	testutil.WantError(t, err)

	verr, ok := err.(*errorutil.ValidationError)
	testutil.GotWant(t, ok, true)

	got := verr.Has(field.Email, errcode.Format, errtext.Format(field.Email))
	testutil.GotWant(t, got, true)

	got = verr.Has(
		field.Password, errcode.StringLength, errtext.StringLength(field.Password, 8))
	testutil.GotWant(t, got, true)
}

func TestValidate_EmailFormatAndPasswordRequirements(t *testing.T) {
	t.Parallel()

	req := &Request{
		Email:    "invalid-email",
		Password: "кирилица",
	}

	err := Validate(req)
	testutil.WantError(t, err)

	verr, ok := err.(*errorutil.ValidationError)
	testutil.GotWant(t, ok, true)

	// Email errors
	got := verr.Has(field.Email, errcode.Format, errtext.Format(field.Email))
	testutil.GotWant(t, got, true)

	// Password errors - all character type requirements
	got = verr.Has(
		field.Password, rerrcode.PasswordUpper, rerrtext.PasswordUpper)
	testutil.GotWant(t, got, true)

	got = verr.Has(
		field.Password, rerrcode.PasswordLower, rerrtext.PasswordLower)
	testutil.GotWant(t, got, true)

	got = verr.Has(
		field.Password, rerrcode.PasswordDigit, rerrtext.PasswordDigit)
	testutil.GotWant(t, got, true)

	got = verr.Has(
		field.Password, rerrcode.PasswordSpecial, rerrtext.PasswordSpecial)
	testutil.GotWant(t, got, true)
}
