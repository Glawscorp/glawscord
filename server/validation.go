package server

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

var validChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-"

func validUsername(username string) bool {
	for _, c := range username {
		if !strings.ContainsRune(validChars, c) {
			return false
		}
	}
	if (utf8.RuneCountInString(username)) < 4 || (utf8.RuneCountInString(username)) > 25 {
		return false
	}

	return true

}

func validPassword(password string) bool {
	containsUpper := false
	containsLower := false

	for _, c := range password {
		if !strings.ContainsRune(validChars, c) {

			return false
		}

	}
	if (utf8.RuneCountInString(password)) < 8 || (utf8.RuneCountInString(password)) > 20 {
		return false
	}

	for _, c := range password {
		if unicode.IsUpper(c) {
			containsUpper = true

		}
	}

	for _, c := range password {
		if unicode.IsLower(c) {
			containsLower = true

		}
	}

	if !containsUpper || !containsLower {
		return false
	}
	return true

}
