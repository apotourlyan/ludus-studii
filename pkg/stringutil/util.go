package stringutil

import (
	"regexp"
	"unicode"
)

var (
	uppercaseRegex   = regexp.MustCompile(`[A-Z]`)
	lowercaseRegex   = regexp.MustCompile(`[a-z]`)
	digitRegex       = regexp.MustCompile(`[0-9]`)
	specialCharRegex = regexp.MustCompile(`[!@#$%\^&*()\[\]\-_+={}|;:,.<>?]`)
)

func ContainsUppercase(text string) bool {
	return uppercaseRegex.MatchString(text)
}

func ContainsLowercase(text string) bool {
	return lowercaseRegex.MatchString(text)
}

func ContainsDigit(text string) bool {
	return digitRegex.MatchString(text)
}

func ContainsSpecial(text string) bool {
	return specialCharRegex.MatchString(text)
}

func IsWhitespace(text string) bool {
	if text == "" {
		return false
	}

	for _, r := range text {
		if !unicode.IsSpace(r) {
			return false
		}
	}

	return true
}

func Capitalize(text string) string {
	if text == "" {
		return text
	}
	runes := []rune(text)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
