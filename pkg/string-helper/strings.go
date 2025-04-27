package string_helper

import (
	"bytes"
	"strings"
	"unicode"
	"unsafe"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var printer = message.NewPrinter(language.English)

func SnakeToCamel(snakeCase string) string {
	words := strings.Split(snakeCase, "_")
	tc := cases.Title(language.English)
	if len(words) > 1 {
		for i := 1; i < len(words); i++ {
			words[i] = tc.String(words[i])
		}
	}
	return strings.Join(words, "")
}

func Slugify(s string) string {
	var buf bytes.Buffer

	for _, r := range s {
		switch {
		case r > unicode.MaxASCII:
			continue
		case unicode.IsLetter(r):
			buf.WriteRune(unicode.ToLower(r))
		case unicode.IsDigit(r), r == '_', r == '-':
			buf.WriteRune(r)
		case unicode.IsSpace(r):
			buf.WriteRune('-')
		}
	}

	return buf.String()
}

// StringToBytes converts a string to a byte slice without a memory allocation, the returned slice must not be modified.
func StringToBytes(s string) []byte {
	if s == "" {
		return nil
	}
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// BytesToString converts a byte slice to a string without a memory allocation, the returned string must not be modified.
func BytesToString(s []byte) string {
	if len(s) == 0 {
		return ""
	}
	return unsafe.String(unsafe.SliceData(s), len(s))
}
