package errorutil

import (
	"net/http"
)

const (
	CodeSystem       = "SYSTEM_ERROR"
	CodeDatabase     = "DATABASE_ERROR"
	CodeRequest      = "REQUEST_ERROR"
	CodeRequired     = "REQUIRED"
	CodeFormat       = "FORMAT"
	CodeStringLength = "STRING_LENGTH"
)

func GetBaseCodeMap() map[string]int {
	m := make(map[string]int)
	m[CodeSystem] = http.StatusInternalServerError
	m[CodeDatabase] = http.StatusInternalServerError
	m[CodeRequest] = http.StatusBadRequest
	return m
}
