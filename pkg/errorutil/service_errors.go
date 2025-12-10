package errorutil

import (
	"github.com/apotourlyan/ludus-studii/pkg/errorutil/errcode"
	"github.com/apotourlyan/ludus-studii/pkg/errorutil/errtext"
)

func DatabaseError(cause error) error {
	if cause == nil {
		return nil
	}

	return Wrap(errcode.Database, errtext.Database, cause)
}

func SystemError(cause error) error {
	if cause == nil {
		return nil
	}

	return Wrap(errcode.System, errtext.System, cause)
}

func RequestError(cause error) error {
	if cause == nil {
		return nil
	}

	return Wrap(errcode.Request, errtext.Request, cause)
}
