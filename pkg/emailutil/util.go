package emailutil

import (
	"errors"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^([a-zA-Z0-9_%+\-]+(?:\.[a-zA-Z0-9_%+\-]+)*)@([a-zA-Z0-9]+(?:[.\-][a-zA-Z0-9]+)*\.[a-zA-Z]{2,})$`)

var ErrInvalidFormat = errors.New("invalid email format")

// Parts represents the parsed components of an email address
type Parts struct {
	Local  string
	Domain string
}

// Parse extracts the local and domain parts from an email address.
// Returns ErrInvalidFormat if the email format is invalid.
func Parse(text string) (*Parts, error) {
	matches := emailRegex.FindStringSubmatch(text)
	if len(matches) < 3 {
		return nil, ErrInvalidFormat
	}

	return &Parts{
		Local:  matches[1],
		Domain: matches[2],
	}, nil
}

func IsValid(text string) bool {
	parts, err := Parse(text)
	if err != nil {
		return false
	}

	// Local part must not exceed 64 characters
	if len(parts.Local) > 64 {
		return false
	}

	// Total domain must not exceed 253 characters
	if len(parts.Domain) > 253 {
		return false
	}

	// Each domain label must not exceed 63 characters
	for label := range strings.SplitSeq(parts.Domain, ".") {
		if len(label) > 63 {
			return false
		}
	}

	return true
}
