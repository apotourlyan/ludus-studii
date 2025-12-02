package httputil

import (
	"errors"
	"net/http"

	"github.com/apotourlyan/ludus-studii/pkg/errorutil"
)

type RequestResult struct {
	Code int
	Data any
}

func CreatedResult(data any) *RequestResult {
	return &RequestResult{
		Code: http.StatusCreated,
		Data: data,
	}
}

func ErrorResult(err error, codeMap map[string]int) *RequestResult {
	// Check for ValidationError
	var validationErr *errorutil.ValidationError
	if errors.As(err, &validationErr) {
		return ValidationErrorResult(validationErr)
	}

	// Check for ServiceError
	var svcErr *errorutil.ServiceError
	if errors.As(err, &svcErr) {
		return ServiceErrorResult(svcErr, codeMap)
	}

	// Unknown error
	return InternalErrorResult()
}

func InternalErrorResult() *RequestResult {
	return &RequestResult{
		Code: http.StatusInternalServerError,
		Data: map[string]string{
			"error": "Internal server error",
		},
	}
}

func ValidationErrorResult(err *errorutil.ValidationError) *RequestResult {
	return &RequestResult{
		Code: http.StatusBadRequest,
		Data: err,
	}
}

func ServiceErrorResult(err *errorutil.ServiceError, errorMap map[string]int) *RequestResult {
	// Get status code from map, default to 500 if not found
	code, ok := errorMap[err.Code()]
	if !ok {
		code = http.StatusInternalServerError
	}

	return &RequestResult{
		Code: code,
		Data: err,
	}
}
