package errtext

import (
	"fmt"

	"github.com/apotourlyan/ludus-studii/pkg/stringutil"
)

const (
	Database = "An unexpected database error has occured."
	System   = "An unexpected system error has occured."
	Request  = "The request is invalid."
)

func Required(name string) string {
	return fmt.Sprintf("%v is required", stringutil.Capitalize(name))
}

func Format(name string) string {
	return fmt.Sprintf("%v format is invalid", stringutil.Capitalize(name))
}

func StringLength(name string, limit int) string {
	return fmt.Sprintf(
		"%v must be at least %d characters",
		stringutil.Capitalize(name),
		limit,
	)
}
