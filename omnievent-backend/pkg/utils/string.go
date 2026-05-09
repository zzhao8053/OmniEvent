package utils

import (
	"regexp"
	"strings"
)

func IsNullOrEmpty(s string) bool {
	return s == ""
}

func IsNullOrWhiteSpace(s string) bool {
	return strings.TrimSpace(s) == ""
}

func Trim(s string) string {
	return strings.TrimSpace(s)
}

func TrimToNull(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	return s
}

func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func ContainsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func HasPrefix(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

func HasSuffix(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

func Replace(s, old, new string) string {
	return strings.Replace(s, old, new, -1)
}

func ReplaceN(s, old, new string, n int) string {
	return strings.Replace(s, old, new, n)
}

func Split(s, sep string) []string {
	return strings.Split(s, sep)
}

func Join(elems []string, sep string) string {
	return strings.Join(elems, sep)
}

func ToLower(s string) string {
	return strings.ToLower(s)
}

func ToUpper(s string) string {
	return strings.ToUpper(s)
}

func Capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func IsValidUsername(username string) bool {
	pattern := regexp.MustCompile(`^(?i)[a-z0-9_-]+$`)
	return len(username) <= 32 && pattern.MatchString(username)
}

func IsValidEmail(email string) bool {
	emailPattern := regexp.MustCompile(`^(?i)(?:[a-z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+)*|"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)\])$`)
	return len(email) <= 100 && emailPattern.MatchString(email)
}

func IsValidNickName(nickname string) bool {
	return len(nickname) <= 64
}

var hexRGBColorPattern = regexp.MustCompile(`^(?i)([0-9a-f]{6}|[0-9a-f]{3})$`)

func IsValidHexRGBColor(color string) bool {
	return hexRGBColorPattern.MatchString(color)
}
