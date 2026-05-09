package utils

import (
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

func SubString(str string, start int, length int) string {
	chars := []rune(str)
	realLength := len(chars)
	end := 0

	if start < 0 {
		start = realLength + start
	}

	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}

	if start > realLength {
		start = realLength
	}

	if end < 0 {
		end = 0
	}

	if end > realLength {
		end = realLength
	}

	return string(chars[start:end])
}

func ContainsAnyString(s string, substrs []string) bool {
	for i := 0; i < len(substrs); i++ {
		if strings.Index(s, substrs[i]) >= 0 {
			return true
		}
	}

	return false
}

func GetFirstLowerCharString(s string) string {
	if s == "" {
		return s
	}

	chars := []rune(s)

	if chars[0] >= 'a' && chars[0] <= 'z' {
		return s
	}

	chars[0] = chars[0] + 32
	return string(chars)
}

func ContainsOnlyOneRune(s string, r rune) bool {
	if len(s) < 1 {
		return false
	}

	for i := 0; i < len(s); i++ {
		if rune(s[i]) != r {
			return false
		}
	}

	return true
}

