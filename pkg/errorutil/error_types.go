package errorutil

type ErrorType string

const (
	TypeService    ErrorType = "service"
	TypeValidation ErrorType = "validation"
)
