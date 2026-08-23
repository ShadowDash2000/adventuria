package activities

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	regParens = regexp.MustCompile(`\s*(\(.*?\)|\[.*?]|{.*?})\s*`)
	regSpaces = regexp.MustCompile(`\s+`)
)

func NormalizeTitle(title string) string {
	s := strings.TrimSpace(title)
	if s == "" {
		return s
	}

	s = regParens.ReplaceAllString(s, " ")

	s = strings.ToLower(s)

	s = strings.Map(func(r rune) rune {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			return r
		case r == '\'' || r == '.' || r == '/' || r == '\\':
			return r
		default:
			return ' '
		}
	}, s)

	s = regSpaces.ReplaceAllString(s, " ")

	s = strings.TrimSpace(s)

	if s == "" {
		return strings.ToLower(title)
	}

	return s
}
