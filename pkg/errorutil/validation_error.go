package errorutil

import (
	"encoding/json"
)

type validationErrorDto struct {
	Type ErrorType        `json:"type"`
	Data []fieldErrorsDto `json:"data"`
}

type fieldErrorsDto struct {
	Name   string          `json:"name"`
	Errors []fieldErrorDto `json:"errors"`
}

type fieldErrorDto struct {
	Code    string `json:"code"`
	Message string `json:"message"`
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
	fmap := make(map[string][]fieldErrorDto)
	for _, fe := range e.data {
		fmap[fe.field] = append(fmap[fe.field], fieldErrorDto{
			Code:    fe.code,
			Message: fe.message,
		})
	}

	var data []fieldErrorsDto
	for name, errors := range fmap {
		data = append(data, fieldErrorsDto{
			Name:   name,
			Errors: errors,
		})
	}

	return json.Marshal(validationErrorDto{
		Type: TypeValidation,
		Data: data,
	})
}
