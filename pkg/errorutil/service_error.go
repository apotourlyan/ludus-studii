package errorutil

import (
	"encoding/json"
	"fmt"
)

type serviceErrorDto struct {
	Type    ErrorType `json:"type"`
	Code    string    `json:"code"`
	Message string    `json:"message"`
}

type ServiceError struct {
	code    string
	message string
	cause   error
}

func NewServiceError(code string, message string) error {
	return &ServiceError{code: code, message: message}
}

func Wrap(code string, message string, cause error) error {
	if cause == nil {
		return nil
	}

	return &ServiceError{code, message, cause}
}

func (e *ServiceError) Unwrap() error {
	return e.cause
}

func (e *ServiceError) Code() string {
	return e.code
}

func (e *ServiceError) Error() string {
	if e.cause == nil {
		return fmt.Sprintf("%s: %s", e.code, e.message)
	}

	return fmt.Sprintf("%s: %s: %v", e.code, e.message, e.cause)
}

func (e *ServiceError) MarshalJSON() ([]byte, error) {
	return json.Marshal(serviceErrorDto{
		Type:    TypeService,
		Code:    e.code,
		Message: e.message,
	})
}
