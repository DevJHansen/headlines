package utils

import (
	"strings"
	"unicode"
)

func TrimWhiteSpace(s string) string {

	start := 0

	for start < len(s) && unicode.IsSpace(rune(s[start])) {
		start++
	}

	end := len(s)

	for end > 0 && unicode.IsSpace(rune(s[end-1])) {
		end--
	}

	return s[start:end]
}

func GetImgUrlFromStyleAtr(s string) string {
	startIndex := strings.Index(s, "('")
	endIndex := strings.Index(s, "')")

	return s[startIndex+2 : endIndex]
}
