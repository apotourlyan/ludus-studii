package register

import (
	"strings"

	"github.com/apotourlyan/ludus-studii/pkg/emailutil"
	"github.com/apotourlyan/ludus-studii/pkg/errorutil"
	"github.com/apotourlyan/ludus-studii/pkg/stringutil"
	"github.com/apotourlyan/ludus-studii/services/user/internal/service/user/register/errcode"
	"github.com/apotourlyan/ludus-studii/services/user/internal/service/user/register/errtext"
	"github.com/apotourlyan/ludus-studii/services/user/internal/service/user/register/field"
)

type Request struct {
	Email    string
	Password string
}

func Validate(m *Request) error {
	m.Email = strings.TrimSpace(m.Email)
	m.Password = strings.TrimSpace(m.Password)

	errors := []errorutil.FieldError{}
	errors = validateEmail(errors, m.Email)
	errors = validatePassword(errors, m.Password)
	return errorutil.NewValidationError(errors)
}

func validateEmail(
	errors []errorutil.FieldError, email string,
) []errorutil.FieldError {
	if email == "" {
		return append(errors, *errorutil.FieldErrorRequired(field.Email))
	}

	if !emailutil.IsValid(email) {
		return append(errors, *errorutil.FieldErrorFormat(field.Email))
	}

	return errors
}

func validatePassword(
	errors []errorutil.FieldError, password string,
) []errorutil.FieldError {
	if password == "" {
		return append(errors, *errorutil.FieldErrorRequired(field.Password))
	}

	if len(password) < 8 {
		return append(
			errors, *errorutil.FieldErrorStringLength(field.Password, 8))
	}

	if !stringutil.ContainsUppercase(password) {
		e := *errorutil.NewFieldError(
			field.Password, errcode.PasswordUpper, errtext.PasswordUpper)
		errors = append(errors, e)
	}

	if !stringutil.ContainsLowercase(password) {
		e := *errorutil.NewFieldError(
			field.Password, errcode.PasswordLower, errtext.PasswordLower)
		errors = append(errors, e)
	}

	if !stringutil.ContainsDigit(password) {
		e := *errorutil.NewFieldError(
			field.Password, errcode.PasswordDigit, errtext.PasswordDigit)
		errors = append(errors, e)
	}

	if !stringutil.ContainsSpecial(password) {
		e := *errorutil.NewFieldError(
			field.Password, errcode.PasswordSpecial, errtext.PasswordSpecial)
		errors = append(errors, e)
	}

	return errors
}
