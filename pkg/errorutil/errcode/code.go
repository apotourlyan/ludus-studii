package errcode

import (
	"net/http"
)

const (
	System       = "SYSTEM_ERROR"
	Database     = "DATABASE_ERROR"
	Request      = "REQUEST_ERROR"
	Required     = "REQUIRED"
	Format       = "FORMAT"
	StringLength = "STRING_LENGTH"
)

func GetBaseCodeMap() map[string]int {
	m := make(map[string]int)
	m[System] = http.StatusInternalServerError
	m[Database] = http.StatusInternalServerError
	m[Request] = http.StatusBadRequest
	return m
}
