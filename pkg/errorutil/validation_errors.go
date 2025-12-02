package errorutil

import "fmt"

func FieldErrorRequired(name string) *FieldError {
	message := fmt.Sprintf("%v is required", name)
	return NewFieldError(name, CodeRequired, message)
}

func FieldErrorFormat(name string) *FieldError {
	message := fmt.Sprintf("%v format is invalid", name)
	return NewFieldError(name, CodeFormat, message)
}

func FieldErrorStringLength(name string, limit int) *FieldError {
	message := fmt.Sprintf("%v must be at least %d characters", name, limit)
	return NewFieldError(name, CodeStringLength, message)
}
