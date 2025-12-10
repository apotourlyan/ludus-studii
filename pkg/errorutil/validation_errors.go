package errorutil

import (
	"github.com/apotourlyan/ludus-studii/pkg/errorutil/errcode"
	"github.com/apotourlyan/ludus-studii/pkg/errorutil/errtext"
)

func FieldErrorRequired(name string) *FieldError {
	return NewFieldError(name, errcode.Required, errtext.Required(name))
}

func FieldErrorFormat(name string) *FieldError {
	return NewFieldError(name, errcode.Format, errtext.Format(name))
}

func FieldErrorStringLength(name string, limit int) *FieldError {
	return NewFieldError(
		name, errcode.StringLength, errtext.StringLength(name, limit))
}
