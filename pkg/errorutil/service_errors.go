package errorutil

func DatabaseError(cause error) error {
	if cause == nil {
		return nil
	}

	message := "An unexpected database error has occured."
	return Wrap(CodeDatabase, message, cause)
}

func SystemError(cause error) error {
	if cause == nil {
		return nil
	}

	message := "An unexpected system error has occured."
	return Wrap(CodeSystem, message, cause)
}

func RequestError(cause error) error {
	if cause == nil {
		return nil
	}

	message := "The request is invalid."
	return Wrap(CodeRequest, message, cause)
}
