package errorutil

import (
	"encoding/json"

	"github.com/apotourlyan/ludus-studii/pkg/errorutil/errtype"
)

type ValidationErrorDto struct {
	Type errtype.Type     `json:"type"`
	Data []FieldErrorsDto `json:"data"`
}

type FieldErrorsDto struct {
	Name   string          `json:"name"`
	Errors []FieldErrorDto `json:"errors"`
}

type FieldErrorDto struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (dto *ValidationErrorDto) Has(field, code, text string) bool {
	for _, x := range dto.Data {
		if x.Name == field {
			for _, e := range x.Errors {
				if e.Code == code && e.Message == text {
					return true
				}
			}
		}
	}

	return false
}

type ValidationError struct {
	data []FieldError
}

type FieldError struct {
	field   string
	code    string
	message string
}

func NewValidationError(errors []FieldError) error {
	if len(errors) == 0 {
		return nil
	}

	return &ValidationError{data: errors}
}

func NewFieldError(name string, code string, message string) *FieldError {
	return &FieldError{name, code, message}
}

func (e *ValidationError) Error() string {
	return e.data[0].message
}

func (e *ValidationError) MarshalJSON() ([]byte, error) {
	// Group errors by field
	fmap := make(map[string][]FieldErrorDto)
	for _, fe := range e.data {
		fmap[fe.field] = append(fmap[fe.field], FieldErrorDto{
			Code:    fe.code,
			Message: fe.message,
		})
	}

	var data []FieldErrorsDto
	for name, errors := range fmap {
		data = append(data, FieldErrorsDto{
			Name:   name,
			Errors: errors,
		})
	}

	return json.Marshal(ValidationErrorDto{
		Type: errtype.Validation,
		Data: data,
	})
}
